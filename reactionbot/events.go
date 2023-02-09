package reactionbot

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jordanleven/reaction-bot/internal/slackclient"
)

type eventReactionStatus int

// Enums for handling eligible and ineligible posts
const (
	eventReactionEligible eventReactionStatus = iota
	eventReactionIneligibleUserReactedToDoesNotExist
	eventReactionIneligibleMessageContainsDisallowedAttachmentType
	eventReactionIneligibleUnknownReason
)

const (
	eventReactionIneligibleMessageIntro = "Hey {user}, Reaction Bot here :robot_face:! It looks like you recently added a \":{reactionType}:\" reaction to a message (\"{message}\") that I'm unable to work with."
	eventReactionIneligibleMessageOutro = "Sorry about that, but I hope we can still be friends! :heart:"
)

var eventReactionStatusErrorMessages = map[eventReactionStatus]string{
	eventReactionIneligibleUserReactedToDoesNotExist:               "Posting a reaction on messages that originate from Slack Bots and Apps is unsupported.",
	eventReactionIneligibleMessageContainsDisallowedAttachmentType: "Posting a reaction on messages that contain `{attachmentExtension}` attachments is currently not supported by Slack.",
	eventReactionIneligibleUnknownReason:                           "",
}

var eventReactionDisallowedAttachmentTypes = []string{
	"webm",
}

func eventReactionMessageContainsDisallowedAttachmentType(event slackclient.ReactionEvent) bool {
	// Early exist if no attachment is contained
	if event.MessageAttachment == nil {
		return false
	}

	attachmentExtension := filepath.Ext(event.MessageAttachment.Name)

	for _, v := range eventReactionDisallowedAttachmentTypes {
		vF := fmt.Sprintf(".%s", v)
		if attachmentExtension == vF {
			return true
		}
	}

	return false
}

func (r *reactionBot) eventReactionIsEligibleForPost(event slackclient.ReactionEvent) (bool, eventReactionStatus) {
	// User IDs will only exist for messages posted by real users
	userIdDoesNotExist := event.UserIDReactedTo == ""
	messageContainsDisallowedAttachment := eventReactionMessageContainsDisallowedAttachmentType(event)

	messageIsEligibleForPost := !userIdDoesNotExist && !messageContainsDisallowedAttachment

	switch {
	case messageIsEligibleForPost:
		return messageIsEligibleForPost, eventReactionEligible
	case userIdDoesNotExist:
		return messageIsEligibleForPost, eventReactionIneligibleUserReactedToDoesNotExist
	case messageContainsDisallowedAttachment:
		return messageIsEligibleForPost, eventReactionIneligibleMessageContainsDisallowedAttachmentType
	default:
		return messageIsEligibleForPost, eventReactionIneligibleUnknownReason
	}
}

func (r *reactionBot) eventReactionIsRegistered(event slackclient.ReactionEvent) bool {
	return r.reactionIsRegistered(event.ReactionEmoji)
}

func eventReactionIsUnique(event slackclient.ReactionEvent) bool {
	return event.ReactionCount == 1
}

func (r *reactionBot) sentEventIneligibleMessageToUser(status eventReactionStatus, event slackclient.ReactionEvent) {
	allUsers := r.Users
	reactedByUser := getUserByUserID(*allUsers, event.UserIDReactedBy)
	var attachmentExtension string

	if event.MessageAttachment != nil {
		attachmentExtension = filepath.Ext(event.MessageAttachment.Name)
	}

	statusSpecificErrorMessage := eventReactionStatusErrorMessages[status]
	errorMessage := fmt.Sprintf("%s %s \n \n %s", eventReactionIneligibleMessageIntro, statusSpecificErrorMessage, eventReactionIneligibleMessageOutro)

	errorMessageReplacements := strings.NewReplacer(
		"{user}", reactedByUser.DisplayName,
		"{message}", event.Message,
		"{attachmentExtension}", attachmentExtension,
		"{reactionType}", event.ReactionEmoji,
	)
	errorMessageF := errorMessageReplacements.Replace(errorMessage)

	slackclient.SendDirectMessage(r.SlackClient, event.UserIDReactedBy, errorMessageF)
	color.Yellow("Sent error message to %s. Message unable to reacted to (code %d). Original message was \"%s\". (dated %s).\n", reactedByUser.DisplayName, status, event.Message, event.ReactionTimestamp)
}

func (r *reactionBot) handleReaction(event slackclient.ReactionEvent) {
	// Early exist if the reaction isn't something we need to deal with
	if !r.eventReactionIsRegistered(event) || !eventReactionIsUnique(event) {
		return
	}

	isEligibleForPost, status := r.eventReactionIsEligibleForPost(event)
	if isEligibleForPost {
		r.maybePostReactedMessageToChannel(event)
	} else {
		r.sentEventIneligibleMessageToUser(status, event)
	}
}

func (r *reactionBot) handleEvents() {
	slackclient.HandleSlackEvents(r.SlackClient, r.handleReaction)
}
