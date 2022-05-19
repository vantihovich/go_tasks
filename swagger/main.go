package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgr "github.com/vantihovich/go_tasks/tree/master/swagger/postgres"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	cnfg "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
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

	log.WithFields(log.Fields{}).Info("Configs loading")

	cfg, err := cnfg.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("Failed to load app config")
	}

	log.WithFields(log.Fields{}).Info("Connecting to DB")
	db := postgr.New(cfg)
	if err := db.Open(); err != nil {
		log.WithFields(log.Fields{}).Panic("Failed to establish DB connection")
	}

	UsersProvider := handlers.NewUsersHandler(db)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			handlers.RegisterNewUser(w, r)
		})

		r.Post("/login", UsersProvider.UserLogin)
	})
	r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	return r
}
