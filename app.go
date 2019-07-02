package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"google.golang.org/api/iterator"
)

// ProjectID name of the project
const ProjectID = "bookshelf-2019-3"

func main() {
	ctx := context.Background()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
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
