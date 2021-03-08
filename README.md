# Reaction Bot

A simple bot that can listen for reactions with specific emoji in Slack.

## Installing Reaction Bot

To install this package in your existing Go project, run:

```sh
go get github.com/jordanleven/reaction-bot
```

Import into your file:

```go
import "github.com/jordanleven/reaction-bot/reactionbot"
```

## Initializing Reaction Bot

To initialize Reaction Bot, run `reactionbot.New()` with the your options.

```go
func main() {
  registrationOptions := reactionbot.RegistrationOptions{
    SlackTokenApp:   slackTokenApp,
    SlackTokenBot:   slackTokenBot,
    RegisteredEmoji: registeredReactions,
  }

  reactionbot.New(registrationOptions)
}
```

### Configuration

1. `SlackTokenApp`: The Slack app token obtained from Slack.
1. `SlackTokenBot`: The Slack bot token obtained from Slack.
1. `RegisteredEmoji`: The list of emoji to register reactions to (see Registering Reactions below).

### Retrieving the authentication credentials

Although you may retrieve the tokens from your env file yourself, you may also use the `GetSlackTokenApp` and `GetSlackTokenBot` functions retrieve them (see example above). Each of these functions accepts the name of the environmental variable to retrieve the tokens from, and will subsequently check the formatting of the tokens to confirm they are the correct ones.

```go
  slackTokenApp := reactionbot.GetSlackTokenApp("SLACK_TOKEN_APP")
  slackTokenBot := reactionbot.GetSlackTokenBot("SLACK_TOKEN_BOT")
```

### Registering Emoji

You may register an unlimited number of emoji reactions in `RegisteredEmoji`. When the emoji identified as the key is used in any channel where Reaction Bot is added, the message will be posted to the identified channel.

```go
  // The name of the bot reaction. This is only used when logging reactions on the server.
  Name         string
  // The channel to post the reaction to.
  Channel      string
```

When registered in `RegisteredEmoji`, a bot named "Mr. Randomness" with that will post messages to the "random" channel when reacted to with the "laughing" emoji would look like this.

```go
var registeredReactions = reactionbot.RegisteredReactions{
  "laughing": {
    Name:    "Randomness",
    Channel: "random",
  },
}
```
