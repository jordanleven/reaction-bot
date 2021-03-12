package reactionbot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient = slack.Client
type SlackSocketClient = socketmode.Client
type SlackUser = slack.User
type SlackReactionAddedEvent = slackevents.ReactionAddedEvent
type SlackPostMessageOptions struct {
	ReactedByUser User
	Attachment    *ReactionAttachment
	Message       string
}

func newSlackClient(options RegistrationOptions) *SlackClient {
	return slack.New(
		options.SlackTokenBot,
		slack.OptionDebug(false),
		slack.OptionAppLevelToken(options.SlackTokenApp),
	)
}

func (r reactionBot) getSlackUsers() ([]SlackUser, error) {
	return r.SlackClient.GetUsers()
}

func (r reactionBot) newSlackSocketClient() *SlackSocketClient {
	return socketmode.New(r.SlackClient)
}

func getNumberOfMessageReactions(message slack.Message, reactionEmoji string) int {
	reactions := message.Reactions
	reactionCount := 0
	for _, reaction := range reactions {
		messageEmoji := reaction.Name
		if messageEmoji == reactionEmoji {
			reactionCount = reaction.Count
		}
	}
	return reactionCount
}

func messageIsReactedMessage(emoji string, timestamp string, message slack.Message) bool {
	messageHasCorrectReaction := false
	messageTimestamp := message.Timestamp
	messageReactions := message.Reactions
	for _, reaction := range messageReactions {
		messageReactionEmoji := reaction.Name
		if messageReactionEmoji == emoji {
			messageHasCorrectReaction = true
		}
	}

	return messageTimestamp == timestamp && messageHasCorrectReaction
}

func (r reactionBot) getReactedMessage(reactionEmoji string, reactionItem slackevents.Item) slack.Message {
	var reactedMessage slack.Message
	timestamp := reactionItem.Timestamp
	channelID := reactionItem.Channel
	payload := slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: timestamp,
		// Required to show messages that are at the limit of the timestamp
		Inclusive: true,
	}
	conversationHistory, _, _, _ := r.Slack.GetConversationReplies(&payload)

	messageTextFound := false
	conversationHistoryLength := len(conversationHistory)
	index := 0
	for !messageTextFound && index < conversationHistoryLength {
		for messageIndex, message := range conversationHistory {
			if messageIsReactedMessage(reactionEmoji, timestamp, message) {
				messageTextFound = true
				reactedMessage = message
			}
			index = messageIndex
		}
	}
	return reactedMessage
}

func (r reactionBot) getFormattedEvent(innerEvent slackevents.EventsAPIInnerEvent) ReactionEvent {
	evt := innerEvent.Data.(*SlackReactionAddedEvent)
	evtReactionEmoji := evt.Reaction
	evtItem := evt.Item
	reactedMessage := r.getReactedMessage(evtReactionEmoji, evtItem)
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

func getInnerEvent(event socketmode.Event) slackevents.EventsAPIInnerEvent {
	evt, _ := event.Data.(slackevents.EventsAPIEvent)
	return evt.InnerEvent
}

func (r reactionBot) handleSlackEvents(callback func(ReactionEvent)) {
	client := r.newSlackSocketClient()
	go func() {
		for evt := range client.Events {
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
				color.Green("Well hello there! Reaction Bot has finish starting up.")
			case socketmode.EventTypeEventsAPI:
				client.Ack(*evt.Request)
				innerEvent := getInnerEvent(evt)
				if innerEvent.Type == slackevents.ReactionAdded {
					event := r.getFormattedEvent(innerEvent)
					callback(event)
				}
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}

func (r reactionBot) postSlackMessage(channel string, opts SlackPostMessageOptions) (timestamp string, error error) {
	reactedMessageBlock := slack.NewTextBlockObject(slack.MarkdownType, opts.Message, false, false)
	blocks := []slack.Block{
		slack.NewSectionBlock(reactedMessageBlock, nil, nil),
	}

	reactionAttachments := slack.Attachment{}

	if opts.Attachment != nil {
		reactionAttachments.ImageURL = opts.Attachment.Permalink
		reactionAttachments.Text = opts.Attachment.Name
	}

	_, ts, error := r.Slack.PostMessage(
		channel,
		slack.MsgOptionBlocks(blocks...),
		// Fallback text
		slack.MsgOptionText(opts.Message, true),
		slack.MsgOptionAttachments(reactionAttachments),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionIconURL(opts.ReactedByUser.ProfileImage),
		slack.MsgOptionParse(true),
		slack.MsgOptionUsername(opts.ReactedByUser.DisplayName),
	)

	return ts, error
}
