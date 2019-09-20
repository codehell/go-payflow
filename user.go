package main

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
)

// UserRole is type for user roles
type UserRole int

const (
	userRoleAdmin UserRole = iota + 1
	userRoleEditor
	userRoleReader
	userRoleUser
)

// User is the struct for the application users
type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
	CreateAt int64    `json:"create_at"`
}

// SetUser save an user at firestore
func (u *User) SetUser(projectID string) (string, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return "", errors.New("firestoreNewClient")
	}
	defer client.Close()
	ref, _, err := client.Collection("users").Add(ctx, u)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

// Render is the method for the Render interface
/* func (a User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
} */
