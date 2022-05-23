package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vantihovich/go_tasks/tree/master/swagger/models"
)

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type UsersHandler struct {
	userRepo models.UserRepository
}

func NewUsersHandler(userRepo models.UserRepository) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
	}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

var ErrNoRows = errors.New("no users with provided credentials were found in database")

func (h *UsersHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := loginRequest{}
	response := loginResponse{}

	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		if err == io.EOF {
			log.WithError(err).Info("Empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.WithError(err).Info("Error occurred")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.FindByLoginAndPwd(parameters.Login, parameters.Password)
	if err != nil {
		if err == ErrNoRows {
			log.WithError(err).Error("No rows found")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Please check user login or password"))
			return
		}
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.UserID = user.UserID

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
		log.WithError(err).Info("Error encoding struct to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
