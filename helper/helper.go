package helper

import (
	"crypto/rand"
	"fmt"
	"strings"
)

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

// CreateUserFromUsernameAndResource returns the complete user with username@resource
func CreateUserFromUsernameAndResource(username, resource string) string {
	return fmt.Sprintf("%v@%v", username, resource)
}

// RandomString is used for testing purposes
func RandomString() string {
	length := 8
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}