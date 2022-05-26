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
	"github.com/vantihovich/go_tasks/tree/master/swagger/validators"
)

type UsersHandler struct {
	userRepo models.UserRepository
}

func NewUsersHandler(userRepo models.UserRepository) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
	}
}

var ErrNoRows = errors.New("no users with provided credentials were found in database")

type registrationRequest struct {
	Login            string   `json:"login"`
	Password         string   `json:"password"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Email            string   `json:"email"`
	SocialMediaLinks []string `json:"social_media_links"`
}

func (h *UsersHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	parameters := registrationRequest{}

	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		if err == io.EOF {
			log.WithError(err).Info("Empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.WithError(err).Info("Error decoding json request occurred")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check if obligatory fields are fullfilled, although looks like it should be done
	//when the json body is formed
	valid := validators.ValidateRegistrationRequest(parameters.Login, parameters.Password, parameters.FirstName, parameters.LastName, parameters.Email)
	if !valid {
		log.Error("One or several obligatory fields in the request body are empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check if login is unique
	exists, err := h.userRepo.CheckIfLoginExists(ctx, parameters.Login)
	if err != nil {
		log.WithError(err).Info("DB request of login returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exists {
		log.Info("Login exists already")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Please choose another login"))
		return
	}

	log.Debug("Login not found in DB, registration is allowed to proceed with provided login")

	//an attempt to add new user
	err = h.userRepo.AddNewUser(ctx, parameters.Login, parameters.Password, parameters.FirstName, parameters.LastName, parameters.Email, parameters.SocialMediaLinks)
	if err != nil {
		log.WithError(err).Info("Error occurred when adding user to DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug("User added succesfully")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

func (h *UsersHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	user, err := h.userRepo.FindByLoginAndPwd(ctx, parameters.Login, parameters.Password)
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
