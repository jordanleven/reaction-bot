package lib

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

func getReactedMessageText(slackInstance *slack.Client, reactionEmoji string, reactionItem slackevents.Item) string {
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
	reactedMessage := ""
	for !messageTextFound && index < conversationHistoryLength {
		for messageIndex, message := range conversationHistory {
			messageText := message.Text
			if messageIsReactedMessage(reactionEmoji, timestamp, message) {
				messageTextFound = true
				reactedMessage = messageText
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
	channelToPostReaction := GetReactionChannelByReaction(reactionEmoji)
	reactionType := GetReactionTypeByEmoji(reactionEmoji)
	reactionItem := reactionEvent.Item
	reactedMessageText := getReactedMessageText(slackInstance, reactionEmoji, reactionItem)
	channelPostReactionMessage := getFormattedMessage(reactedByName, reactedToName, reactedMessageText)

	_, _, err := slackInstance.PostMessage(
		channelToPostReaction,
		slack.MsgOptionText(channelPostReactionMessage, false),
		slack.MsgOptionAttachments(),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Successfully sent a \"%s\" reaction to the %s channel.", reactionType, channelToPostReaction)
}
