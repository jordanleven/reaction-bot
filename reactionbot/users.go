package reactionbot

import (
	"github.com/fatih/color"
	"github.com/jordanleven/reaction-bot/internal/slackclient"
)

func getUserByUserID(u slackclient.Users, uid string) slackclient.User {
	return u[uid]
}

func (r *reactionBot) getUsers() slackclient.Users {
	users, err := slackclient.GetSlackUsers(r.SlackClient)
	if err != nil {
		color.Red("Error getting users: %s\n", err)
	}

	return users
}

func (r *reactionBot) updateUsers() {
	*r.Users = r.getUsers()
}
