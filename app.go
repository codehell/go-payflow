package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"google.golang.org/api/iterator"
)

// ProjectID name of the project
const ProjectID = "go-payflow"

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	ctx := context.Background()
	r := chi.NewRouter()

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(csrfMiddleware)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		// Save it before we write to the response/return from the handler.
		session.Save(r, w)
		_, _ = w.Write([]byte("test"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CRSF-Token", csrf.Token(r))
		user := User{
			Username: "codehell",
			Email:    "admin@codehell.net",
		}
		_ = render.Render(w, r, user)
	})

	r.Get("/restaurants", func(w http.ResponseWriter, r *http.Request) {
		client, err := firestore.NewClient(ctx, ProjectID)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		var arry []map[string]interface{}
		col := client.Collection("restaurants")
		iter := col.Documents(ctx)
		for {
			dr, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if dr != nil {
				arry = append(arry, dr.Data())
			}
		}
		res, _ := json.Marshal(arry)
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
