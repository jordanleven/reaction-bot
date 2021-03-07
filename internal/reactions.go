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

var developmentReaction = map[string]RegisteredReaction{
	// The registered emoji on development reactions should be obscure, less frequently used emoji.
	"white_check_mark": {
		Name:         "Reaction Bot Testing",
		BotName:      "Reaction Bot Development",
		BotIconEmoji: ":test_tube:",
		Channel:      "reaction-bot-testing",
	},
}

func getRegisteredReactionByEmoji(emoji string) (RegisteredReaction, bool) {
	var registeredEmojiSet map[string]RegisteredReaction

	// If we're working locally, we'll load our developmentReaction set to avoid issue posting
	// reactions to production channels
	if isDevelopmentEnvironment() {
		registeredEmojiSet = developmentReaction
	} else {
		registeredEmojiSet = registeredReactions
	}

	registeredReaction, registeredReactionWasFound := registeredEmojiSet[emoji]
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
