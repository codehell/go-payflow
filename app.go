package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/codehell/notifiertester/config"
	"github.com/codehell/notifiertester/middlewares"
	"github.com/codehell/notifiertester/payflow"
	"github.com/codehell/notifiertester/users"
	"github.com/codehell/notifiertester/utils"

	"github.com/alexedwards/scs/v2"
	"github.com/codehell/notifiertester/point"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"gopkg.in/yaml.v2"
)

const (
	ConfigFile = "conf.yaml"
)

var (
	appConfig      config.AppConfig
	projectID      string
	sessionManager *scs.SessionManager
)

func init() {
	body, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal("unable to read file")
	}

	err = yaml.Unmarshal(body, &appConfig)
	if err != nil {
		log.Fatal("unable to unmarshal file")
	}

	projectID = appConfig.ProjectID
}

func main() {
	r := chi.NewRouter()
	sf := utils.StringHandlerFunc(projectID)
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	r.Use(middleware.Logger)
	r.Use(middlewares.SetJSONContentType)

	// Api Payflow
	r.Get("/api/payflow", payflow.GetPayflowNotification(projectID))
	r.Post("/api/payflow/ok", sf(payflow.SetPayflowOkUrl))
	r.Post("/api/payflow/error", sf(payflow.SetPayflowErrorUrl))
	r.Post("/api/payflow/notification", sf(payflow.SetPayflowNotification))
	r.Post("/api/payflow/notification-error", sf(payflow.SetPayflowErrorNotification))

	r.Get("/api/set-session", func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Put(r.Context(), "message", "Hello from a session!")
		pa := point.Point{
			X: 10,
			Y: 7,
		}
		jPoint, err := json.Marshal(pa)
		if err != nil {
			utils.APIResponse(w, "error: i am dumb sorry", "imDumb", 500)
		}
		_, _ = w.Write(jPoint)
	})

	r.Get("/api/get-session", func(w http.ResponseWriter, r *http.Request) {
		msg := sessionManager.GetString(r.Context(), "message")
		utils.APIResponse(w, "sessionData", msg, http.StatusOK)
	})

	r.Get("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		utils.APIResponse(w, "pong", "testResponse", http.StatusOK)
	})

	r.Get("/api/version", func(w http.ResponseWriter, r *http.Request) {
		utils.APIResponse(w, runtime.Version(), "version", http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		isProduction := os.Getenv("GCP_ENVIRONMENT") == "production"
		csrfOption := csrf.Secure(isProduction)
		csrfMiddleware := csrf.Protect([]byte(appConfig.Key), csrfOption)

		r.Use(csrfMiddleware)

		r.Get("/api/csrf", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-CSRF-Token", csrf.Token(r))
		})

		r.Post("/api/users", func(w http.ResponseWriter, r *http.Request) {
			users.PostUser(projectID, w, r)
		})
	})

	if err := http.ListenAndServe(":5000", sessionManager.LoadAndSave(r)); err != nil {
		log.Fatal(err)
	}
}
