package helper

import "strings"

// GetUsernameAndResourceFromUser returns given user, username and resource
func GetUsernameAndResourceFromUser(user string) (string, string) {
	if len(user) == 0 {
		return "", ""
	}

	// user is of the form username@resource
	userSlice := strings.Split(user, "@")

	if len(userSlice) != 2 {
		return "", ""
	}

	return userSlice[0], userSlice[1]
}