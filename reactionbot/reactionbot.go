package reactionbot

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// ReactionBot is our main structure
type ReactionBot struct {
	Slack           *slack.Client
	IsDevelopment   bool
	RegisteredEmoji RegisteredReactions
	Users           SlackUsers
}

// RegistrationOptions the list of options to init the package
type RegistrationOptions struct {
	SlackTokenApp   string
	SlackTokenBot   string
	RegisteredEmoji RegisteredReactions
}

func getSlackInstance(options RegistrationOptions) *slack.Client {
	return slack.New(
		options.SlackTokenBot,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(options.SlackTokenApp),
	)
}

func getReactionBot(options RegistrationOptions) ReactionBot {
	slack := getSlackInstance(options)
	slackUsers := GetSlackWorkspaceUsers(slack)

	bot := ReactionBot{
		Slack:           slack,
		RegisteredEmoji: options.RegisteredEmoji,
		Users:           slackUsers,
	}
	return bot
}

// New function to init the package
func New(options RegistrationOptions) {
	bot := getReactionBot(options)
	RegisterSlackBot(bot)
}
