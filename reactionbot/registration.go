package reactionbot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// RegisterSlackBot is the function used to start the bot and listen for reactions
func (bot ReactionBot) RegisterSlackBot() {

	client := socketmode.New(bot.Slack)

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
					if bot.ReactionIsRegistered(reactionEmoji) {
						bot.HandleMessageReaction(reactionAddedEvent)
					}
				}

			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}
