package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	cnfg "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
	"github.com/vantihovich/go_tasks/tree/master/swagger/handlers"
	mw "github.com/vantihovich/go_tasks/tree/master/swagger/middleware"
	postgr "github.com/vantihovich/go_tasks/tree/master/swagger/postgres"
	"github.com/vantihovich/go_tasks/tree/master/swagger/redis"
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
			log.WithError(err).Fatal("Could not shutdown the server")
		}

		serverStopCtx()
		shutDownCnclFunc()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("An error starting server")
	}

	<-serverCtx.Done()
}

func service() http.Handler {
	log.Info("Configs loading")

	cfgDB, err := cnfg.LoadDB()
	if err != nil {
		log.WithError(err).Fatal("Failed to load DB config")
	}

	cfgJWT, err := cnfg.LoadJWT()
	if err != nil {
		log.WithError(err).Fatal("Failed to load JWT config")
	}

	cfgLogin, err := cnfg.LoadLogin()
	if err != nil {
		log.WithError(err).Fatal("Failed to load Login config")
	}

	cfgRedis, err := cnfg.LoadRedisConfigs()
	if err != nil {
		log.WithError(err).Fatal("Failed to load Redis config")
	}

	maxAllowedInvalidLogins, err := strconv.Atoi(cfgLogin.MaxAllowedInvalidLogins)
	if err != nil {
		log.WithError(err).Fatal("Failed to convert env parameter to int type")
	}

	log.Info("Connecting to Redis")
	cache, err := redis.New(cfgRedis.RedisServerConnectionType, cfgRedis.RedisServer)
	if err != nil {
		log.WithError(err).Fatal("Failed to establish connection with Redis")
	}

	log.Info("Connecting to DB")
	db := postgr.New(cfgDB)
	if err := db.Open(); err != nil {
		log.WithError(err).Fatal("Failed to establish connection with DB")
	}

	UsersProvider := handlers.NewUsersHandler(&db, cache, cfgJWT.SecretKey, cfgLogin.InvalidLoginAttemptTTL, maxAllowedInvalidLogins)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", UsersProvider.RegisterNewUser)
		r.Post("/login", UsersProvider.UserLogin)
		r.Post("/deactivate", mw.Authorize(cfgJWT.SecretKey, UsersProvider.UserDeactivation))
		r.Post("/password_reset", mw.Authorize(cfgJWT.SecretKey, UsersProvider.PasswordReset))

	})
	r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	return r
}
