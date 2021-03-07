package internal

//RegisteredReaction are all emoji that are registered to a channel
type RegisteredReaction struct {
	Name         string
	BotName      string
	BotIconEmoji string
	Channel      string
}

var registeredReactions = map[string]RegisteredReaction{
	"laughing": {
		Name:         "Randomness",
		BotName:      "Mr. Randomness",
		BotIconEmoji: ":laughing:",
		Channel:      "random",
	},
	"bulb": {
		Name:         "Today I learned",
		BotName:      "TIL Bot",
		BotIconEmoji: ":bulb:",
		Channel:      "til",
	},
}

func getRegisteredReactionByEmoji(emoji string) (RegisteredReaction, bool) {
	registeredReaction, registeredReactionWasFound := registeredReactions[emoji]
	return registeredReaction, registeredReactionWasFound
}

// GetRegisteredReaction returns true if the reaction emoji has been registered
func GetRegisteredReaction(emoji string) RegisteredReaction {
	reaction, _ := getRegisteredReactionByEmoji(emoji)
	return reaction
}

// ReactionIsRegistered returns true if the reaction emoji has been registered
func ReactionIsRegistered(emoji string) bool {
	_, reactionWasFound := getRegisteredReactionByEmoji(emoji)
	return reactionWasFound
}
