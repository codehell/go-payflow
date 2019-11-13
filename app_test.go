package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	response, err := http.Get("http://localhost:8080/api/ping")
	if err != nil {
		t.Error(err)
		return
	}
	if http.StatusOK != response.StatusCode {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.StatusCode)
		return
	}
	data := make(map[string]string)
	_ = json.NewDecoder(response.Body).Decode(&data)
	description, ok := data["description"]
	if !ok {
		t.Error("Expected a description")
		return
	}
	if description != "pong" {
		t.Errorf("Expected description %s. Got %s", "pong", description)
	}
}

func TestCrsf(t *testing.T) {
	response, err := http.Get("http://localhost:8080/api/csrf")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Error("Expected status code ok")
	}
	log.Println(response.Header.Get("Set-Cookie"))
}

func ExampleHello() {
	fmt.Println("hello")
	// Output: hello
}
