# Reaction Bot

A simple bot that can listen for reactions with specific emoji in Slack.

## Configuring Reaction Bot

In `internal/reactions.go`, you may register an unlimited number of emoji reactions in `registeredReactions`. When the emoji identified as the key is used in any channel where Reaction Bot is added, the message will be posted to the identified channel.

```go
  // The name of the bot reaction. This is only used when logging reactions on the server.
  Name         string
  // The username of the bot when posting the reaction to te channel.
  BotName      string
  // The emoji of the bot when posting the reaction to the channel.
  BotIconEmoji string
  // The channel to post the reaction to.
  Channel      string
```

When registered in `registeredReactions`, a bot named "Mr. Randomness" with that will post messages to the "random" channel when reacted to with the "laughing" emoji would look like this.

```go
"laughing": {
  Name:         "Randomness",
  BotName:      "Mr. Randomness",
  BotIconEmoji: ":joy:",
  Channel:      "random",
},
```

## Running locally

### Setting up your local environment

To run this app, copy the contents of `.env.sample` to `.env`. Replace the values of the tokens according to the
API bot credentials obtained from Slack.

### Starting the app

To install dependencies, run the following:

```sh
go get
```

After installing, you can run the app by running the following command:

```sh
go run cmd/main.go
```

After successfully starting up, you will see the startup greeting:

```sh
➜ go run cmd/main.go
Connecting to Slack with Socket Mode...
Connected to Slack with Socket Mode.
Well hello there! Reaction Bot has finish starting up.
```

If the bot was unable to start up in Socket Mode, it'll attempt to reconnect. If it is unable to connect, please confirm your authentication credentials and that you're not exceeding the rate limits set by Slack.
