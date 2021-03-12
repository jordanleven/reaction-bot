package reactionbot

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/slack-go/slack"
)

const refreshIntervalInHours = 4

type reactionBot struct {
	Slack           *slack.Client
	SlackClient     *SlackClient
	IsDevelopment   bool
	RegisteredEmoji RegisteredReactions
	Users           *Users
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
	slack := getSlackInstance(options)
	slackClient := newSlackClient(options)

	b := reactionBot{
		Slack:           slack,
		SlackClient:     slackClient,
		RegisteredEmoji: options.RegisteredEmoji,
		Users:           &Users{},
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
