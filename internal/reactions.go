package internal

type reactionChannels struct {
	Name    string
	Channel string
}

var registeredReactions = map[string]reactionChannels{
	"laughing": {
		Name:    "Randomness",
		Channel: "random",
	},
	"bulb": {
		Name:    "Today I learned",
		Channel: "til",
	},
}

func getRegisteredReactionByEmoji(emoji string) (reactionChannels, bool) {
	registeredReaction, registeredReactionWasFound := registeredReactions[emoji]
	return registeredReaction, registeredReactionWasFound
}

// GetReactionChannelByReaction returns the name of the channel for a specific reaction
func GetReactionChannelByReaction(emoji string) string {
	registeredReaction, _ := getRegisteredReactionByEmoji(emoji)
	return registeredReaction.Channel
}

// GetReactionTypeByEmoji returns the name of the type of reaction
func GetReactionTypeByEmoji(emoji string) string {
	registeredReaction, _ := getRegisteredReactionByEmoji(emoji)
	return registeredReaction.Name
}

// ReactionIsRegistered returns true if the reaction emoji has been registered
func ReactionIsRegistered(emoji string) bool {
	_, reactionWasFound := getRegisteredReactionByEmoji(emoji)
	return reactionWasFound
}
