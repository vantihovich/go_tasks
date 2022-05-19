package handlers

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"github.com/vantihovich/go_tasks/tree/master/swagger/postgres"
)

// Not implemented yet
// type registrationRequest struct {
// 	firstName          string   `json: "first_name"`
// 	lastName           string   `json: "last_name"`
// 	email              string   `json: "login"`
// 	social_media_links []string `json: "social_media_links"`
// 	login              string   `json: "login"`
// 	password           string   `json: "password`
// }

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type UsersHandler struct {
	userRepo postgres.DB
}

func NewUsersHandler(userRepo postgres.DB) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
	}
}

func (h *UsersHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := loginRequest{}
	response := loginResponse{}

	err := json.NewDecoder(r.Body).Decode(&parameters)

	if err != nil {
		if err == io.EOF {
			log.WithFields(log.Fields{"Error": err}).Info("Error empty request body")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			log.WithFields(log.Fields{"Error": err}).Info("Error occurred")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = h.userRepo.QueryRow("SELECT user_id FROM users WHERE login=$1 AND password=$2 ", parameters.Login, parameters.Password).Scan(&response.UserId)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.WithFields(log.Fields{"Error": err}).Info("No rows found")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Please check user login or password"))
			return
		} else {
			log.WithFields(log.Fields{"Error": err}).Info("DB request returned error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	log.WithFields(log.Fields{}).Info("User found in DB")

	//at this moment token is just a random integer from 1 up to 100
	//TODO implement JWT
	token := func() string {
		rand.Seed(time.Now().UnixMicro())
		b := rand.Intn(100)

		return strconv.Itoa(b)
	}()

	response.Token = token

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Error encoding struct to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write(response)

}
