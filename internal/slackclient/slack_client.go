package slackclient

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient = slack.Client
type SlackSocketClient = socketmode.Client

func newSocket(client *SlackClient) *SlackSocketClient {
	return socketmode.New(
		client,
		socketmode.OptionDebug(false),
	)
}

// NewClient returns our Slack client
func NewClient(botToken string, appToken string) *SlackClient {
	return slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(false),
	)
}
