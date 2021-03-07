package internal

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func getFormattedMessage(reactedBy string, reactedTo string, reactedMessage string) string {
	return fmt.Sprintf("%s: \"%s\" \n- @%s", reactedBy, reactedMessage, reactedTo)
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

func getReactedMessage(slackInstance *slack.Client, reactionEmoji string, reactionItem slackevents.Item) slack.Message {
	var reactedMessage slack.Message
	timestamp := reactionItem.Timestamp
	channelID := reactionItem.Channel
	payload := slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: timestamp,
		// Required to show messages that are at the limit of the timestamp
		Inclusive: true,
	}
	conversationHistory, _, _, _ := slackInstance.GetConversationReplies(&payload)

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

// PostReactedMessageToChannel is the function used to post a reaction
func PostReactedMessageToChannel(slackInstance *slack.Client, allUsers map[string]SlackUser, reactionEvent *slackevents.ReactionAddedEvent) {
	reactedBy := reactionEvent.User
	reactedByName := GetUsernameUserID(allUsers, reactedBy)
	reactedTo := reactionEvent.ItemUser
	reactedToName := GetUsernameUserID(allUsers, reactedTo)

	reactionEmoji := reactionEvent.Reaction
	registeredReaction := GetRegisteredReaction(reactionEmoji)
	channelToPostReaction := registeredReaction.Channel
	reactionType := registeredReaction.Name
	reactionBotName := registeredReaction.BotName
	reactionBotIcon := registeredReaction.BotIconEmoji
	reactionItem := reactionEvent.Item
	reactedMessage := getReactedMessage(slackInstance, reactionEmoji, reactionItem)
	reactedMessageText := reactedMessage.Text
	reactedMessageFiles := reactedMessage.Files
	reactionAttachments := slack.Attachment{}

	if len(reactedMessageFiles) > 0 {
		// In case someone decides to add a bunch of photos, we're going to limit them to one
		firstReactedFile := reactedMessage.Files[0]
		reactionAttachments.ImageURL = firstReactedFile.Permalink
		reactionAttachments.Text = " "

		// If a user just posted an image, update the message text to be an empty string (the API
		// requires us to post a non-null string)
		if reactedMessageText == "" {
			reactedMessageText = ":camera:"
		}

	}

	reactedMessageTextFormatted := getFormattedMessage(reactedByName, reactedToName, reactedMessageText)

	_, _, err := slackInstance.PostMessage(
		channelToPostReaction,
		slack.MsgOptionText(reactedMessageTextFormatted, false),
		slack.MsgOptionAttachments(reactionAttachments),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionIconEmoji(reactionBotIcon),
		slack.MsgOptionDisableMarkdown(),
		slack.MsgOptionUsername(reactionBotName),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Successfully sent a \"%s\" reaction to the %s channel.\n", reactionType, channelToPostReaction)
}
