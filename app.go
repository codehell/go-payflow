package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"google.golang.org/api/iterator"
)

// ProjectID name of the project
const ProjectID = "go-payflow"

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type documentID struct {
	ID string `json:"documentId"`
}

type Restaurant struct {
	Name string `json:"name"`
}

func (res Restaurant) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func SetJsonContentType (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	isProduction := os.Getenv("GCP_ENVIRONMENT") == "production"
	csrfOption := csrf.Secure(isProduction)
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"), csrfOption)

	r.Use(middleware.Logger)
	r.Use(SetJsonContentType)
	r.Use(csrfMiddleware)

	r.Get("/api/test", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		// Save it before we write to the response/return from the handler.
		_ = session.Save(r, w)
		_, _ = w.Write([]byte("test"))
	})

	r.Get("/api/crsf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CRSF-Token", csrf.Token(r))
		_, _ = w.Write([]byte(""))
	})

	r.Post("/api/restaurants", storeRestaurants)

	r.Get("/api/restaurants", getRestaurants)

	r.Post("/api/users", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			Error(w, "error: decode user data", "errDecodeUser", 500)
			return
		}
		user.Role = UserRoleUser
		user.CreateAt = time.Now().Unix()
		id, err := user.SetUser()
		if err != nil {
			Error(w, err.Error(), "errSetUser", 500)
			return
		}
		docID := documentID{
			ID: id,
		}
		jDocID, err := json.Marshal(docID)
		if err != nil {
			Error(w, err.Error(), "errMarshalUserID,", 500)
			return
		}
		_, _ = w.Write(jDocID)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getRestaurants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CRSF-Token", csrf.Token(r))
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, ProjectID)
	if err != nil {
		log.Println("an error has ocurred")
		Error(w, "wtf Restaurants", "where are my restaurants", 500)
		return
	}
	var arry []Restaurant
	col := client.Collection("restaurants")
	iter := col.Documents(ctx)
	for {
		dr, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if dr != nil {
			nameInter := dr.Data()["name"]
			name := fmt.Sprintf("%v", nameInter)
			arry = append(arry, Restaurant{Name:name})
		}
	}
	res, _ := json.Marshal(arry)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func storeRestaurants(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["hello"] = "world"
	jData, _ := json.Marshal(data)
	_, _ = w.Write(jData)
}

// Error Application error responses constructor
func Error(w http.ResponseWriter, description string, code string, httpCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	appErr := make(map[string]string)
	appErr["code"] = code
	appErr["description"] = description
	jError, _ := json.Marshal(appErr)
	w.WriteHeader(httpCode)
	_, _ = w.Write(jError)
}
