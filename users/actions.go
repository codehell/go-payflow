package users

import (
	"encoding/json"
	"github.com/codehell/go_firestore/utils"
	"net/http"
	"time"
)

func PostUser(projectID string, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.APIResponse(w, "error: decode user data", "errDecodeUser", 500)
		return
	}
	user.Role = userRoleUser
	user.CreateAt = time.Now().Unix()
	id, err := user.SetUser(projectID)
	if err != nil {
		utils.APIResponse(w, err.Error(), "errSetUser", 500)
		return
	}

	description := "User added: " + id
	utils.APIResponse(w, description, "userAdded", http.StatusCreated)
}
