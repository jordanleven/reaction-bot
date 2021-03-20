package reactionbot

import (
	"time"

	"github.com/fatih/color"
	"github.com/jordanleven/reaction-bot/internal/slackclient"
)

const refreshIntervalInHours = 4

type reactionBot struct {
	SlackClient     *slackclient.SlackClient
	IsDevelopment   bool
	RegisteredEmoji RegisteredReactions
	Users           *slackclient.Users
}

// RegistrationOptions the list of options to init the package
type RegistrationOptions struct {
	SlackTokenApp   string
	SlackTokenBot   string
	RegisteredEmoji RegisteredReactions
}

func (b *reactionBot) handleUpdateUsers() {
	ticker := time.NewTicker(time.Hour * refreshIntervalInHours)
	go func() {
		for range ticker.C {
			color.White("Updating users...")
			b.updateUsers()
		}
	}()
}

func getReactionBot(options RegistrationOptions) reactionBot {
	slackClient := slackclient.NewClient(
		options.SlackTokenBot,
		options.SlackTokenApp,
	)

	b := reactionBot{
		SlackClient:     slackClient,
		RegisteredEmoji: options.RegisteredEmoji,
		Users:           &slackclient.Users{},
	}

	b.updateUsers()

	return b
}

// New function to init the package
func New(options RegistrationOptions) {
	bot := getReactionBot(options)
	go bot.handleUpdateUsers()
	defer bot.handleEvents()
}
