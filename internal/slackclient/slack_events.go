package slackclient

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type ReactionAttachment struct {
	Permalink string
	Name      string
}

type ReactionEvent struct {
	ReactionEmoji     string
	ReactionTimestamp string
	ReactionCount     int
	Message           string
	MessageAttachment *ReactionAttachment
	UserIDReactedTo   string
	UserIDReactedBy   string
}

type socketEvent = socketmode.Event
type innerEvent = slackevents.EventsAPIInnerEvent
type slackReactionAddedEvent = slackevents.ReactionAddedEvent

func getInnerEvent(event socketEvent) innerEvent {
	evt, _ := event.Data.(slackevents.EventsAPIEvent)
	return evt.InnerEvent
}

func getFormattedEvent(client *SlackClient, innerEvent innerEvent) ReactionEvent {
	evt := innerEvent.Data.(*slackReactionAddedEvent)
	evtReactionEmoji := evt.Reaction
	evtItem := evt.Item
	reactedMessage := getReactedMessage(client, evtReactionEmoji, evtItem)
	reactionCount := getNumberOfMessageReactions(reactedMessage, evtReactionEmoji)
	formattedEvent := ReactionEvent{
		ReactionEmoji:     evtReactionEmoji,
		ReactionTimestamp: evt.Item.Timestamp,
		ReactionCount:     reactionCount,
		UserIDReactedBy:   evt.User,
		UserIDReactedTo:   evt.ItemUser,
		Message:           reactedMessage.Text,
	}

	if len(reactedMessage.Files) > 0 {
		f := reactedMessage.Files[0]
		formattedEvent.MessageAttachment = &ReactionAttachment{
			Name:      f.Name,
			Permalink: f.Permalink,
		}
	}

	return formattedEvent
}

func handleSlackEventReactionAdded(client *SlackClient, evt socketEvent, callback func(ReactionEvent)) {
	innerEvent := getInnerEvent(evt)
	if innerEvent.Type == slackevents.ReactionAdded {
		event := getFormattedEvent(client, innerEvent)
		callback(event)
	}
}

// HandleSlackEvents executes callbacks for specific reactions
func HandleSlackEvents(client *SlackClient, callback func(ReactionEvent)) {
	socket := newSocket(client)

	go func() {
		for evt := range socket.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				color.White("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				color.Red("Connection failed. Retrying later...")
			case socketmode.EventTypeDisconnect:
				color.Red("Disconnecting!")
			case socketmode.EventTypeConnected:
				color.Green("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeHello:
				numConnections := evt.Request.NumConnections
				color.Green("Well hello there! Reaction Bot has finish starting up (current connections: %d).", numConnections)
			case socketmode.EventTypeEventsAPI:
				socket.Ack(*evt.Request)
				go handleSlackEventReactionAdded(client, evt, callback)
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	socket.Run()
}
