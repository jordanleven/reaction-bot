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

func (bot ReactionBot) getRegisteredReactionByEmoji(emoji string) (Reaction, bool) {

	registeredReaction, registeredReactionWasFound := bot.RegisteredEmoji[emoji]
	return registeredReaction, registeredReactionWasFound
}

// GetRegisteredReaction returns true if the reaction emoji has been registered
func (bot ReactionBot) GetRegisteredReaction(emoji string) Reaction {
	reaction, _ := bot.getRegisteredReactionByEmoji(emoji)
	return reaction
}

// ReactionIsRegistered returns true if the reaction emoji has been registered
func (bot ReactionBot) ReactionIsRegistered(emoji string) bool {
	_, reactionWasFound := bot.getRegisteredReactionByEmoji(emoji)
	return reactionWasFound
}
