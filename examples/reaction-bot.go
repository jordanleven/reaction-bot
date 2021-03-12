package main

import (
	"github.com/jordanleven/reaction-bot/reactionbot"
)

var registeredReactions = reactionbot.RegisteredReactions{
	"raised_hands": {
		Name:    "Reaction Bot Testing",
		Channel: "til",
	},
}

func main() {
	slackTokenApp := reactionbot.GetSlackTokenApp("SLACK_TOKEN_APP")
	slackTokenBot := reactionbot.GetSlackTokenBot("SLACK_TOKEN_BOT")
	registrationOptions := reactionbot.RegistrationOptions{
		SlackTokenApp:   slackTokenApp,
		SlackTokenBot:   slackTokenBot,
		RegisteredEmoji: registeredReactions,
	}
	reactionbot.New(registrationOptions)
}
