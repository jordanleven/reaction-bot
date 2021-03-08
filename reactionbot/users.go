package reactionbot

import (
	"fmt"

	"github.com/slack-go/slack"
)

// SlackUser is any current user of this Slack workspace
type SlackUser struct {
	Username     string
	FullName     string
	DisplayName  string
	ProfileImage string
}

// SlackUsers is a key/value ordered list of users by their UUID
type SlackUsers map[string]SlackUser

func userIsInactive(user slack.User) bool {
	// Don't include bots or deleted users in our list of users
	return user.IsBot ||
		user.Deleted
}

// GetUserByUserID is a function to return a specific user given a user ID
func GetUserByUserID(users SlackUsers, userID string) SlackUser {
	return users[userID]
}

// GetSlackWorkspaceUsers is a function to return all
// users of the workspace
func GetSlackWorkspaceUsers(slackInstance *slack.Client) SlackUsers {
	users, err := slackInstance.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	userDictionary := make(SlackUsers)
	for _, user := range users {
		// Don't include bots or deleted users in our list of users
		if userIsInactive(user) {
			continue
		}
		userDictionary[user.ID] = SlackUser{
			Username:     user.Name,
			FullName:     user.RealName,
			DisplayName:  user.Profile.DisplayNameNormalized,
			ProfileImage: user.Profile.Image512,
		}
	}

	return userDictionary
}