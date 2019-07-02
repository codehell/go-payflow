package main

import "net/http"

// User is the struct for the application users
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Render is the method for the Render interface
func (a User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
