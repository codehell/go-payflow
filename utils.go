package main

import (
	"encoding/json"
	"net/http"
)

// APIResponse Application error responses constructor
func APIResponse(w http.ResponseWriter, description string, code string, httpCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	appErr := make(map[string]string)
	appErr["code"] = code
	appErr["description"] = description
	jError, _ := json.Marshal(appErr)
	w.WriteHeader(httpCode)
	_, _ = w.Write(jError)
}
