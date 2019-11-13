package payflow

import (
	"context"
	"github.com/codehell/go_firestore/utils"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// PayflowNotification struct
type PayflowNotification struct {
	DataReturn          string `json:"dataReturn"`
	SignatureDataReturn string `json:"signatureDataReturn"`
}

func SetPayflowNotification(pid string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		if false {
			utils.APIResponse(w, "error: add payflow document", "errTimeout", http.StatusGatewayTimeout)
			return
		}
		err := r.ParseForm()
		if err != nil {
			utils.APIResponse(w, "error: add payflow document", "errParseForm", http.StatusBadRequest)
			return
		}
		notification := PayflowNotification{
			r.Form.Get("dataReturn"),
			r.Form.Get("signatureDataReturn"),
		}
		ctx := context.Background()

		client, err := firestore.NewClient(ctx, pid)
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
		utils.APIResponse(w, "cool", "200", http.StatusOK)
	}

}

func GetPayflowNotification(pid string) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		utils.APIResponse(w, pid, "200", http.StatusOK)
	}
}
