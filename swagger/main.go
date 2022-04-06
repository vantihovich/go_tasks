package main

import (
	_ "embed"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/vantihovich/go_tasks/tree/master/swagger/handlers"
)

//go:embed  api/apiauth.yaml
var spec []byte

func main() {
	r := chi.NewRouter()

	r.HandleFunc("/auth/register", handlers.RegisterNewUser)
	r.HandleFunc("/auth/login", handlers.UserLogin)
	r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(spec)))

	log.Println("serving on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error starting server")
	}

}
