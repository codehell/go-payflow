package main

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// PayflowNotification struct
type PayflowNotification struct {
	DataReturn          string `json:"dataReturn"`
	SignatureDataReturn string `json:"signatureDataReturn"`
}

func setPayflowNotification(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if true {
		APIResponse(w, "error: add payflow document", "errTimeout", http.StatusGatewayTimeout)
		return
	}
	err := r.ParseForm()
	if err != nil {
		APIResponse(w, "error: add payflow document", "errParseForm", http.StatusBadRequest)
		return
	}
	not := PayflowNotification{
		r.Form.Get("dataReturn"),
		r.Form.Get("signatureDataReturn"),
	}
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, ProjectID)
	defer client.Close()
	if err != nil {
		APIResponse(w, "error: create firestore client", "badFirestoreClient", http.StatusBadRequest)
		return
	}
	_, _, err = client.Collection("payflow").Add(ctx, not)
	if err != nil {
		APIResponse(w, "error: add payflow document", "badFirestoreDoc", http.StatusBadRequest)
		return
	}
	APIResponse(w, "cool", "200", http.StatusOK)
}

func getPayflowNotifications(w http.ResponseWriter, r *http.Request) {
	APIResponse(w, "cool", "200", http.StatusOK)
}

func testErrorResponse(w http.ResponseWriter, r *http.Request) {
	APIResponse(w, "test response for api", "500", http.StatusInternalServerError)
}
