package slackclient

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SlackPostMessageOptions struct {
	ReactedByUser User
	Attachment    *ReactionAttachment
	Message       string
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

func getReactedMessage(client *SlackClient, reactionEmoji string, reactionItem slackevents.Item) slack.Message {
	var reactedMessage slack.Message
	timestamp := reactionItem.Timestamp
	channelID := reactionItem.Channel
	payload := slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: timestamp,
		// Required to show messages that are at the limit of the timestamp
		Inclusive: true,
	}
	conversationHistory, _, _, _ := client.GetConversationReplies(&payload)

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

// PostSlackMessage allows posting messages to specific channels
func PostSlackMessage(client *SlackClient, channel string, opts SlackPostMessageOptions) (timestamp string, error error) {
	reactedMessageBlock := slack.NewTextBlockObject(slack.MarkdownType, opts.Message, false, false)

	blocks := []slack.Block{
		slack.NewSectionBlock(reactedMessageBlock, nil, nil),
	}

	reactionAttachments := slack.Attachment{}

	if opts.Attachment != nil {
		reactionAttachments.ImageURL = opts.Attachment.Permalink
		reactionAttachments.Text = opts.Attachment.Name
	}

	_, ts, error := client.PostMessage(
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
