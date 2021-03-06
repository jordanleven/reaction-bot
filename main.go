package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	l "github.com/jordanleven/reaction-bot/lib"
)

func getSlackInstance() *slack.Client {
	slackTokenBot := l.GetSlackTokenBot()
	slackTokenApp := l.GetSlackTokenApp()
	return slack.New(
		slackTokenBot,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(slackTokenApp),
	)
}

func registerSlackBot() {
	slackInstance := getSlackInstance()
	slackUsers := l.GetSlackWorkspaceUsers(slackInstance)
	client := socketmode.New(slackInstance)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				color.White("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				color.Red("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				color.Green("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeHello:
				color.Green("Well hello there! Reaction Bot has finish starting up.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, _ := evt.Data.(slackevents.EventsAPIEvent)
				client.Ack(*evt.Request)
				innerEvent := eventsAPIEvent.InnerEvent
				if innerEvent.Type == slackevents.ReactionAdded {
					reactionAddedEvent := innerEvent.Data.(*slackevents.ReactionAddedEvent)
					reactionEmoji := reactionAddedEvent.Reaction
					if l.ReactionIsRegistered(reactionEmoji) {
						l.PostReactedMessageToChannel(slackInstance, slackUsers, reactionAddedEvent)
					}
				}

			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}

func main() {
	registerSlackBot()
}
