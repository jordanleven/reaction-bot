# Reaction Bot

A simple bot that can listen for reactions with specific emoji in Slack.

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
