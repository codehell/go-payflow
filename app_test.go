package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	response, err := http.Get("http://localhost:8080/api/ping")
	if http.StatusOK != response.StatusCode {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.StatusCode)
	}
	data := make(map[string]string)
	json.NewDecoder(response.Body).Decode(&data)
	description, ok := data["description"]
	if !ok {
		t.Error("Expected a description")
	}
	if description != "pong" {
		t.Errorf("Expected description %s. Got %s", "pong", description)
	}
	if err != nil {
		t.Errorf("Encountered an error: %v", err)
	}
}
