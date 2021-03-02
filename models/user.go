package models

import (
	"errors"
	"fmt"
)

// User represents a generic user based object with minimal properties (ID, FirstName, LastName)
type User struct {
	ID        int
	FirstName string
	LastName  string
}


var (
	users []*User
	nextID = 1
)

// Return pointer array of User structs
func GetUsers() []*User {
	return users
}

// Adds new user to users var, returns error if ID included
func AddUser(u User) (User, error)  {
	if u.ID != 0 {
		return User{}, errors.New("new user must not include ID or it must be set to zero")
	}

	u = assignID(u)
	users = append(users, &u)
	return u, nil
}

// Updates supplied user details if ID can be found, if not user is delegated to AddUser()
func UpdateUser(id int, u User) (User, error) {
	for _, user := range users {
		if user.ID == id {
			*user = u
			return *user, nil
		}
	}
	return AddUser(u)
}

// Returns specific user by provided ID. Returns error if user cant be found
func GetUserById(id int) (User, error) {
	for _, user := range users {
		if user.ID == id {
			return *user, nil
		}
	}
	return User{}, fmt.Errorf("user with id %v not found", id)
}

// Deletes user by provided ID. Returns error if user cant be found
func RemoveUserById(id int) error {
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user with id %v not found", id)
}

// Assigns incremental ID to passed user
func assignID(u User) User {
	u.ID = nextID
	nextID++
	return u
}