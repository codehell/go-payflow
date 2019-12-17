package payflow

import (
	"context"
	"github.com/codehell/go_firestore/utils"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// PayflowNotification struct
type Notification struct {
	DataReturn          string `json:"dataReturn"`
	SignatureDataReturn string `json:"signatureDataReturn"`
	CreateAt            int64  `json:"date"`
}

func SetPayflowNotification(pid string) http.HandlerFunc {
	return commonPayflow(pid, false)
}

func SetPayflowErrorNotification(pid string) http.HandlerFunc {
	return commonPayflow(pid, true)
}

func SetPayflowOkUrl(pid string) http.HandlerFunc {
	return commonPayflow(pid, false)
}

func SetPayflowErrorUrl(pid string) http.HandlerFunc {
	return commonPayflow(pid, false)
}

func commonPayflow(projectID string, withError bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			utils.APIResponse(w, "error: add payflow document", "errParseForm", http.StatusBadRequest)
			return
		}
		notification := Notification{
			r.Form.Get("dataReturn"),
			r.Form.Get("signatureDataReturn"),
			time.Now().UTC().Unix(),
		}
		ctx := context.Background()

		client, err := firestore.NewClient(ctx, projectID)
		defer client.Close()
		if err != nil {
			utils.APIResponse(w, "error: create firestore client", "badFirestoreClient", http.StatusBadRequest)
			return
		}
		_, _, err = client.Collection("payflow").Add(ctx, notification)
		if err != nil {
			utils.APIResponse(w, "error: add payflow document", "badFirestoreDoc", http.StatusBadRequest)
			return
		}
		time.Sleep(2 * time.Second)
		if withError {
			utils.APIResponse(w, "error: mock error endpoint", "errTimeout", http.StatusInternalServerError)
			return
		}
		utils.APIResponse(w, "cool", "200", http.StatusOK)
	}
}

func GetPayflowNotification(pid string) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		utils.APIResponse(w, pid, "200", http.StatusOK)
	}
}
