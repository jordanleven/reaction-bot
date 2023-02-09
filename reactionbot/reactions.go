package reactionbot

//Reaction is an individual registered reaction
type Reaction struct {
	Name         string
	BotName      string
	BotIconEmoji string
	Channel      string
}

// RegisteredReactions are the registered reactions
type RegisteredReactions map[string]Reaction

func (b *reactionBot) getRegisteredReactionByEmoji(emoji string) (Reaction, bool) {
	registeredReaction, registeredReactionWasFound := b.RegisteredEmoji[emoji]
	return registeredReaction, registeredReactionWasFound
}

func (b *reactionBot) getRegisteredReaction(emoji string) Reaction {
	reaction, _ := b.getRegisteredReactionByEmoji(emoji)
	return reaction
}

func (b *reactionBot) reactionIsRegistered(emoji string) bool {
	_, reactionWasFound := b.getRegisteredReactionByEmoji(emoji)
	return reactionWasFound
}
