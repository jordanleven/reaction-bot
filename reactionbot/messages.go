package reactionbot

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jordanleven/reaction-bot/internal/slackclient"
)

//SlackErrorMessageContent the error message from slack when the content is unavailable
const SlackErrorMessageContent string = "This content can't be displayed."

func getAuthorAttribution(reactedTo string) string {
	return fmt.Sprintf("\n- @%s", reactedTo)
}

func getFormattedMessage(reactedTo string, reactedMessage string) string {
	attribution := getAuthorAttribution(reactedTo)
	return fmt.Sprintf("\"%s\" %s", reactedMessage, attribution)
}

func (r *reactionBot) postReactedMessageToChannel(channel string, event slackclient.ReactionEvent) (timetamp string, error error) {
	allUsers := r.Users
	reactedByUser := getUserByUserID(*allUsers, event.UserIDReactedBy)
	reactedToUser := getUserByUserID(*allUsers, event.UserIDReactedTo)
	reactedToName := reactedToUser.Username
	reactedMessage := event.Message
	reactedMessageFormatted := getFormattedMessage(reactedToName, reactedMessage)

	if event.MessageAttachment != nil {
		// If a user just posted an image, update the message text to be an empty string (the API
		// requires us to post a non-null string)
		if reactedMessage == "" {
			authorAttribution := getAuthorAttribution(reactedToName)
			reactedMessageFormatted = fmt.Sprintf(":camera: %s", authorAttribution)
		}
	}

	opts := slackclient.SlackPostMessageOptions{
		ReactedByUser: reactedByUser,
		Attachment:    event.MessageAttachment,
		Message:       reactedMessageFormatted,
	}
	return slackclient.PostSlackMessage(r.SlackClient, channel, opts)
}

func (r *reactionBot) maybePostReactedMessageToChannel(event slackclient.ReactionEvent) {
	registeredReaction := r.getRegisteredReaction(event.ReactionEmoji)
	reactionName := registeredReaction.Name
	reactionChannel := registeredReaction.Channel
	reactedMessageText := event.Message

	if reactedMessageText == SlackErrorMessageContent {
		color.Red("Unable to retrieve message for %s reaction to message (dated %s).\n", reactionName, event.ReactionTimestamp)
		color.Red("Message data is posted below\n")
		return
	}

	_, err := r.postReactedMessageToChannel(reactionChannel, event)

	if err != nil {
		color.Red("Error posting message to Slack: %s\n", err)
		return
	}

	color.Green("Successfully sent a \"%s\" reaction to the %s channel.\n", reactionName, reactionChannel)
}
