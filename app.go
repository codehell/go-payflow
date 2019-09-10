package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/codehell/go_firestore/point"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"gopkg.in/yaml.v2"
)

var (
	appConfig      AppConfig
	projectID      string
	sessionManager *scs.SessionManager
)

func init() {
	body, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal("unable to read file")
	}

	err = yaml.Unmarshal(body, &appConfig)
	if err != nil {
		log.Fatal("unable to unmarshal file")
	}

	projectID = appConfig.ProjectID
	log.Printf("Config: %v", appConfig)
}

func main() {
	r := chi.NewRouter()

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	r.Use(middleware.Logger)
	r.Use(SetJSONContentType)

	// Api Payflow
	r.Get("/api/payflow", getPayflowNotifications)
	r.Post("/api/payflow", setPayflowNotification)
	r.Post("/api/test/error/response", testErrorResponse)

	r.Get("/api/set-session", func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Put(r.Context(), "message", "Hello from a session!")
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

	r.Get("/api/get-session", func(w http.ResponseWriter, r *http.Request) {
		msg := sessionManager.GetString(r.Context(), "message")
		APIResponse(w, "sessionData", msg, http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		isProduction := os.Getenv("GCP_ENVIRONMENT") == "production"
		csrfOption := csrf.Secure(isProduction)
		csrfMiddleware := csrf.Protect([]byte(appConfig.Key), csrfOption)

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
			id, err := user.SetUser(projectID)
			if err != nil {
				APIResponse(w, err.Error(), "errSetUser", 500)
				return
			}

			description := "User added: " + id
			APIResponse(w, description, "userAdded", http.StatusCreated)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", sessionManager.LoadAndSave(r)))
}
