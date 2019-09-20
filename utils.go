package main

import (
	"encoding/json"
	"net/http"
)

// APIResponse Application error responses constructor
func APIResponse(w http.ResponseWriter, description string, code string, httpCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	response := make(map[string]string)
	response["code"] = code
	response["description"] = description
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(httpCode)
	_, _ = w.Write(jsonResponse)
}
