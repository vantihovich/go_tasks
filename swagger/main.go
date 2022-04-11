package main

import (
	"context"
	_ "embed"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
	log "github.com/sirupsen/logrus"
	"github.com/vantihovich/go_tasks/tree/master/swagger/handlers"
)

//go:embed  api/apiauth.yaml
var spec []byte

func main() {
	valv := valve.New()
	baseCtx := valv.Context()

	r := chi.NewRouter()
	r.Route("/auth", func(r chi.Router) {

		r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
			// consider moving all the logic from here to function e.g. handlers.GracefulRegisterHandler
			log.Println("register request came to router") //TODO implement debug level logging
			valve.Lever(r.Context()).Open()
			defer valve.Lever(r.Context()).Close()

			log.Println("Checking if shutting the server down is initiated:") //TODO implement debug level logging

			select {
			case <-valve.Lever(r.Context()).Stop(): //Notifies about initiated server shutdown, provides timeout
				log.Println("valve is closed, finish handling") //TODO implement debug level logging
				time.Sleep(8 * time.Second)                     //TODO Check if some actions are required here
			case <-time.After(20 * time.Second): //time.After here for testing purposes, the code will be: " default: func()" like commented below part
				log.Println("valve is open, proceed handling")
				<-time.After(20 * time.Second)
				handlers.RegisterNewUser(w, r)
				// default:
				// 	log.Println("valve is open, proceed handling") //TODO implement debug level logging
				// 	handlers.RegisterNewUser(w, r)
			}
		})

		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			// consider moving all the logic from here to function e.g. handlers.GracefulLoginHandler
			log.Println("login request came to router") //TODO implement debug level logging
			valve.Lever(r.Context()).Open()
			defer valve.Lever(r.Context()).Close()

			log.Println("Checking if shutting the server down is initiated:") //TODO implement debug level logging
			select {
			case <-valve.Lever(r.Context()).Stop(): //Notifies about initiated server shutdown, provides timeout
				log.Println("valve is closed, finish handling") //TODO implement debug level logging
				time.Sleep(8 * time.Second)                     //TODO Check if some actions are required here
			default:
				log.Println("valve is open, proceed handling") //TODO implement debug level logging
				handlers.UserLogin(w, r)
			}

		})

	})
	r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(spec)))

	srv := &http.Server{ //TODO implement configs for project
		Addr:    ":3000",
		Handler: r,
	}

	srv.BaseContext = func(_ net.Listener) context.Context {
		return baseCtx
	}

	c := make(chan os.Signal, 1) // need to reserve to buffer size 1, so the notifier is not blocked

	go func() {

		signal.Notify(c, os.Interrupt, syscall.SIGTERM) //receiving the command to stop server

	}()

	go func() {

		<-c

		log.Println("Initiated shutting the server down...") //TODO implement debug level logging

		valv.Shutdown(10 * time.Second) // grace period is discussible, 10 sec for testing purposes//TODO move time definition to configs

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //timout is discussible, 10 sec for testing purposes//TODO move time definition to configs
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("Could not gracefully shutdown the server")
		}

		//checking if srv.Shutdown has not managed to finish within context timeout

		<-ctx.Done()
		log.Println("Server not stopped within the timeout") //TODO implement the level of logging
		//TODO logic about cancelling server stopping

	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error starting server")
	}

}
