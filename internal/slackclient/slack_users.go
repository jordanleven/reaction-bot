package slackclient

import "github.com/slack-go/slack"

type SlackUser = slack.User

type User struct {
	Username     string
	FullName     string
	DisplayName  string
	ProfileImage string
	IsBot        bool
}

type Users map[string]User

func userIsInactive(user SlackUser) bool {
	// Don't include Deleted users in our list (since it's unlikely users are
	// reacting to messages that are from now-deleted users)
	return user.Deleted
}

func getFormattedUser(user SlackUser) User {
	return User{
		IsBot:        user.IsBot,
		Username:     user.Name,
		FullName:     user.RealName,
		DisplayName:  user.Profile.DisplayNameNormalized,
		ProfileImage: user.Profile.Image192,
	}
}

func getFormattedUsers(users []SlackUser) Users {
	userDictionary := make(Users)
	for _, u := range users {
		if userIsInactive(u) {
			continue
		}
		userDictionary[u.ID] = getFormattedUser(u)
	}
	return userDictionary
}

// GetSlackUsers returns all current users in the Workspace
func GetSlackUsers(s *SlackClient) (Users, error) {
	users, err := s.GetUsers()
	return getFormattedUsers(users), err
}
