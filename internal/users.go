package internal

import (
	"fmt"

	"github.com/slack-go/slack"
)

// SlackUser is any current user of this Slack workspace
type SlackUser struct {
	Username string
	FullName string
}

func userIsInactive(user slack.User) bool {
	// Don't include bots or deleted users in our list of users
	return user.IsBot ||
		user.Deleted
}

func getUserByUserID(users map[string]SlackUser, userID string) SlackUser {
	return users[userID]
}

// GetUsernameUserID is a function to return a specific user given a user ID
func GetUsernameUserID(users map[string]SlackUser, userID string) string {
	return getUserByUserID(users, userID).Username
}

// GetSlackWorkspaceUsers is a function to return all
// users of the workspace
func GetSlackWorkspaceUsers(slackInstance *slack.Client) map[string]SlackUser {
	users, err := slackInstance.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
		// return nil, err
	}

	userDictionary := make(map[string]SlackUser)
	for _, user := range users {
		// Don't include bots or deleted users in our list of users
		if userIsInactive(user) {
			break
		}
		userDictionary[user.ID] = SlackUser{
			Username: user.Name,
			FullName: user.RealName,
		}
	}

	return userDictionary
}
