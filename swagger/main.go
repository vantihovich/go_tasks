package main

import (
	_ "embed"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	log "github.com/sirupsen/logrus"
	"github.com/vantihovich/go_tasks/tree/master/swagger/handlers"
)

//go:embed  api\apiauth.yaml
var spec []byte

func main() {

	http.HandleFunc("/auth/register", handlers.RegisterNewUser)
	http.HandleFunc("/auth/login", handlers.UserLogin)
	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))

	log.Println("serving on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error starting server")
	}

}
