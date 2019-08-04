package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
)

type UserRole int

const (
	UserRoleAdmin UserRole = iota + 1
	UserRoleEditor
	UserRoleReader
	UserRoleUser
)

// User is the struct for the application users
type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
	CreateAt int64    `json:"create_at"`
}

// SetUser save an user at firestore
func (u *User) SetUser() (string, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, ProjectID)
	if err != nil {
		return "", errors.New("firestoreNewClient")
	}
	ref, _, err := client.Collection("users").Add(ctx, u)
	if err != nil {
		return "", err
	}
	defer client.Close()
	return ref.ID, nil
}

// Render is the method for the Render interface
/* func (a User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
} */
