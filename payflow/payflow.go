package payflow

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/codehell/notifiertester/utils"
)

// PayflowNotification struct
type Notification struct {
	DataReturn          string `json:"dataReturn"`
	SignatureDataReturn string `json:"signatureDataReturn"`
	CreateAt            int64  `json:"date"`
}

func SetPayflowNotification(pid string) http.HandlerFunc {
	return commonPayflow(pid, false, true)
}

func SetPayflowErrorNotification(pid string) http.HandlerFunc {
	return commonPayflow(pid, true, true)
}

func SetPayflowOkUrl(pid string) http.HandlerFunc {
	return commonPayflow(pid, false, false)
}

func SetPayflowErrorUrl(pid string) http.HandlerFunc {
	return commonPayflow(pid, false, false)
}

func commonPayflow(projectID string, withError bool, isNotification bool) http.HandlerFunc {
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
			utils.APIResponse(w, err.Error(), "badFirestoreDoc", http.StatusBadRequest)
			return
		}
		time.Sleep(2 * time.Second)
		if withError {
			utils.APIResponse(w, err.Error(), "errTimeout", http.StatusInternalServerError)
			return
		}
		if isNotification {
			utils.APIResponse(w, "cool notification", "200", http.StatusOK)
		} else {
			utils.APIResponse(w, "cool ok", "200", http.StatusOK)
		}
	}
}

func GetPayflowNotification(pid string) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		utils.APIResponse(w, pid, "200", http.StatusOK)
	}
}
