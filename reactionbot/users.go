package reactionbot

import (
	"github.com/fatih/color"
)

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

func getUserByUserID(u Users, uid string) User {
	return u[uid]
}

func getFormattedUser(user SlackUser) User {
	return User{
		IsBot:        user.IsBot,
		Username:     user.Name,
		FullName:     user.RealName,
		DisplayName:  user.Profile.DisplayNameNormalized,
		ProfileImage: user.Profile.Image512,
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

func (r reactionBot) getUsers() Users {
	users, err := r.getSlackUsers()
	if err != nil {
		color.Red("Error getting users: %s\n", err)
	}

	formattedUsers := getFormattedUsers(users)
	return formattedUsers
}

func (r *reactionBot) updateUsers() {
	*r.Users = r.getUsers()
}
