package main

import (
	"encoding/json"
	"log"
)

// APIError struct for api error response
type APIError struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

func (ae APIError) aError() []byte {
	body, err := json.Marshal(ae)
	if err != nil {
		log.Fatal("error: can not marshal this api error")
	}
	return body
}

func systemError() []byte {
	se := APIError{
		"system_error",
		"System error",
	}
	return se.aError()
}
