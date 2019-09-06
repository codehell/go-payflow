package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go_firestore/point"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

// ProjectID name of the project
const ProjectID = "go-payflow"

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(SetJSONContentType)

	// Api Payflow
	r.Get("/api/payflow", getPayflowNotifications)
	r.Post("/api/payflow", setPayflowNotification)
	r.Post("/api/test/error/response", testErrorResponse)

	r.Get("/api/test", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		// Save it before we write to the response/return from the handler.
		_ = session.Save(r, w)
		fooSession := session.Values["foo"]
		fs := fmt.Sprintf("%v", fooSession)
		fmt.Println(fs)
		pa := point.Point{
			X: 10,
			Y: 7,
		}
		jPoint, err := json.Marshal(pa)
		if err != nil {
			APIResponse(w, "error: i am dumb sorry", "imDumb", 500)
		}
		_, _ = w.Write(jPoint)
	})

	r.Group(func(r chi.Router) {
		isProduction := os.Getenv("GCP_ENVIRONMENT") == "production"
		csrfOption := csrf.Secure(isProduction)
		csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key2224262930x"), csrfOption)

		r.Use(csrfMiddleware)

		r.Get("/api/crsf", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-CRSF-Token", csrf.Token(r))
		})

		r.Post("/api/users", func(w http.ResponseWriter, r *http.Request) {
			var user User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				APIResponse(w, "error: decode user data", "errDecodeUser", 500)
				return
			}
			user.Role = userRoleUser
			user.CreateAt = time.Now().Unix()
			id, err := user.SetUser(ProjectID)
			if err != nil {
				APIResponse(w, err.Error(), "errSetUser", 500)
				return
			}

			description := "User added: " + id
			APIResponse(w, description, "userAdded", http.StatusCreated)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
