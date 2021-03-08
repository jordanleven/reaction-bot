package reactionbot

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func getFormattedMessage(reactedTo string, reactedMessage string) string {
	return fmt.Sprintf("\"%s\" \n- @%s", reactedMessage, reactedTo)
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

func (bot ReactionBot) getReactedMessage(reactionEmoji string, reactionItem slackevents.Item) slack.Message {
	var reactedMessage slack.Message
	timestamp := reactionItem.Timestamp
	channelID := reactionItem.Channel
	payload := slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: timestamp,
		// Required to show messages that are at the limit of the timestamp
		Inclusive: true,
	}
	conversationHistory, _, _, _ := bot.Slack.GetConversationReplies(&payload)

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
func (bot ReactionBot) PostReactedMessageToChannel(reactionEvent *slackevents.ReactionAddedEvent) {
	allUsers := bot.Users
	reactedByUser := GetUserByUserID(allUsers, reactionEvent.User)
	reactedByName := reactedByUser.DisplayName
	reactedToUser := GetUserByUserID(allUsers, reactionEvent.ItemUser)
	reactedToName := reactedToUser.Username

	reactionEmoji := reactionEvent.Reaction
	registeredReaction := bot.GetRegisteredReaction(reactionEmoji)
	channelToPostReaction := registeredReaction.Channel
	reactionType := registeredReaction.Name
	reactionItem := reactionEvent.Item
	reactedMessage := bot.getReactedMessage(reactionEmoji, reactionItem)
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

	reactedMessageTextFormatted := getFormattedMessage(reactedToName, reactedMessageText)

	reactedMessageBlock := slack.NewTextBlockObject("mrkdwn", reactedMessageTextFormatted, false, false)
	blocks := []slack.Block{
		slack.NewSectionBlock(reactedMessageBlock, nil, nil),
	}

	_, _, err := bot.Slack.PostMessage(
		channelToPostReaction,
		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionAttachments(reactionAttachments),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionIconURL(reactedByUser.ProfileImage),
		slack.MsgOptionParse(true),
		slack.MsgOptionUsername(reactedByName),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Successfully sent a \"%s\" reaction to the %s channel.\n", reactionType, channelToPostReaction)
}
