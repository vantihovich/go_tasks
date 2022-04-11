package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/vantihovich/go_tasks/tree/master/swagger/handlers"
)

//go:embed  api/apiauth.yaml
var spec []byte

func main() {
	srv := &http.Server{ //TODO implement configs for project
		Addr:    ":3000",
		Handler: service(),
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {

		<-sig

		shutDownCtx, shutDownCnclFunc := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutDownCtx.Done()
			if shutDownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := srv.Shutdown(shutDownCtx); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Could not shutdown the server")
		}

		serverStopCtx()
		shutDownCnclFunc()

	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error starting server")
	}

	<-serverCtx.Done()

}

func service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {

		r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
			handlers.RegisterNewUser(w, r)
		})

		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			handlers.UserLogin(w, r)
		})

	})
	r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	return r
}
