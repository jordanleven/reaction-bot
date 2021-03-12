package reactionbot

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

func (r reactionBot) messageShouldBePosted(event ReactionEvent) bool {
	return r.reactionIsRegistered(event.ReactionEmoji) && event.ReactionCount == 1
}

func (r reactionBot) handleReaction(event ReactionEvent) {
	if r.messageShouldBePosted(event) {
		r.maybePostReactedMessageToChannel(event)
	}
}

func (r reactionBot) handleEvents() {
	r.handleSlackEvents(r.handleReaction)
}
