package reactionbot

import "github.com/jordanleven/reaction-bot/internal/slackclient"

func (r reactionBot) messageShouldBePosted(event slackclient.ReactionEvent) bool {
	return r.reactionIsRegistered(event.ReactionEmoji) && event.ReactionCount == 1
}

func (r reactionBot) handleReaction(event slackclient.ReactionEvent) {
	if r.messageShouldBePosted(event) {
		r.maybePostReactedMessageToChannel(event)
	}
}

func (r reactionBot) handleEvents() {
	slackclient.HandleSlackEvents(r.SlackClient, r.handleReaction)
}
