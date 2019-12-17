package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	response, err := http.Get("http://localhost:8080/api/ping")
	if err != nil {
		t.Fatal(err)
	}
	if http.StatusOK != response.StatusCode {
		t.Fatalf("Expected response code %d. Got %d\n", http.StatusOK, response.StatusCode)
	}
	data := make(map[string]string)
	_ = json.NewDecoder(response.Body).Decode(&data)
	description, ok := data["description"]
	if !ok {
		t.Fatal("Expected a description")
	}
	if description != "pong" {
		t.Errorf("Expected description %s. Got %s", "pong", description)
	}
}

func TestCrsf(t *testing.T) {
	response, err := http.Get("http://localhost:8080/api/csrf")
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Expected status code ok")
	}
	setCookie := response.Header.Get("Set-Cookie")
	if !strings.Contains(setCookie, "_gorilla_csrf=") {
		t.Error("Expect cookie contains _gorilla_csr=")
	}
}

func TestPostUser(t *testing.T) {
	c := http.DefaultClient
	response, err := http.Get("http://localhost:8080/api/csrf")
	if err != nil {
		t.Fatal(err)
	}

	data := strings.NewReader(`{"user":"codehell", "mail":"admin@codehell.net"}`)
	request, err := http.NewRequest("POST", "http://localhost:8080/api/users", data)
	if err != nil {
		t.Fatal(err)
	}

	token := response.Header.Get("X-Csrf-Token")
	cookie := response.Header.Get("Set-Cookie")
	request.Header.Set("Cookie", cookie)
	request.Header.Set("X-CSRF-Token", token)

	response, err = c.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	got := response.StatusCode
	if got != http.StatusCreated {
		t.Errorf("Expect 201 status code, got %d", got)
	}
}

func ExampleHello() {
	fmt.Println("hello")
	// Output: hello
}
