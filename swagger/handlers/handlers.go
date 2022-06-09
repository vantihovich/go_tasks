package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
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
			log.WithError(err).Info("empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.WithError(err).Info("error decoding request body occurred")
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
		log.Info("login exists already")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Please choose another login"))
		return
	}

	log.Debug("login not found in DB, registration is allowed to proceed with provided login")

	//an attempt to add new user
	err = h.userRepo.AddNewUser(ctx, parameters.Login, parameters.Password, parameters.FirstName, parameters.LastName, parameters.Email, parameters.SocialMediaLinks)
	if err != nil {
		log.WithError(err).Info("error occurred when adding user to DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug("user added succesfully")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func (h *UsersHandler) UserLogin(w http.ResponseWriter, r *http.Request, jwtParam string) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	parameters := loginRequest{}
	response := loginResponse{}

	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		if err == io.EOF {
			log.WithError(err).Info("empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.WithError(err).Info("error decoding request body occurred")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.FindByLoginAndPwd(ctx, parameters.Login, parameters.Password)
	if err != nil {
		if err == ErrNoRows {
			log.WithError(err).Error("no rows found")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Please check user login or password"))
			return
		}
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//check if user is deactivated
	if !user.Active {
		log.Debug("user is deactivated")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User is deactivated"))
		return
	}

	claims := &claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Issuer: "AuthService",
		},
	}

	jwToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	response.Token, err = jwToken.SignedString([]byte(jwtParam))
	if err != nil {
		log.WithError(err).Info("JWT creation returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.WithError(err).Info("error encoding struct to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type deactivationRequest struct {
	Login string `json:"login"`
}

type contextKey string

var ContextKeyUserID = contextKey("userID")

func (h *UsersHandler) UserDeactivation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parameters := deactivationRequest{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		if err == io.EOF {
			log.WithError(err).Info("empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.WithError(err).Info("error decoding request body occurred")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deactivatorID := (r.Context().Value(ContextKeyUserID)).(int)

	//check if logged in user is allowed to deactivate the user from the request
	deactivator, err := h.userRepo.GetAdminAttrUserLogin(ctx, deactivatorID)
	if err != nil {
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !deactivator.Admin && (deactivator.Login != parameters.Login) {
		// user is not admin or trying to deactivate another`s profile
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User wasn`t deactivated"))
		return
	}

	//deactivation request
	exists, err := h.userRepo.DeactivateUser(ctx, parameters.Login)

	if err != nil {
		log.WithError(err).Info("DB request of login returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Info("login does not exist in the system")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please choose another login"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deactivated successfully"))
}
