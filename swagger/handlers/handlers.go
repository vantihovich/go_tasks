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

	//check if user is deactivated
	active, err := h.userRepo.CheckIfUserActive(ctx, parameters.Login)
	if err != nil {
		log.WithError(err).Info("DB request of active attribute returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !active {
		log.Debug("User is deactivated")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User is deactivated"))
		return
	}

	response.UserID = user.UserID

	claims := &claims{
		UserID: response.UserID,
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
		log.WithError(err).Info("Error encoding struct to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type deactivationRequest struct {
	Login string `json:"login"`
}

func (h *UsersHandler) UserDeactivation(w http.ResponseWriter, r *http.Request, jwtParam string) {
	ctx := r.Context()
	parameters := deactivationRequest{}
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		log.Error("empty header authorization data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := &claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam), nil
	})
	if err != nil {
		if err.Error() == jwt.ErrSignatureInvalid.Error() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.WithError(err).Info("Error parsing the token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&parameters)
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

	//check if logged in user is allowed to deactivate the user from the request
	deactivator, err := h.userRepo.GetAdminAttrUserLogin(ctx, claims.UserID)
	if err != nil {
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !deactivator.Admin {
		if deactivator.UserLogin != parameters.Login {
			// user is not admin or trying to deactivate another`s profile
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("User wasn`t deactivated"))
			return
		}
	}

	//check if login from request body exists in DB
	exists, err := h.userRepo.CheckIfLoginExists(ctx, parameters.Login)
	if err != nil {
		log.WithError(err).Info("DB request of login returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Debug("Login does not exist in DB")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid login - doesn`t exist in DB"))
		return
	}

	//check if user is deactivated already
	active, err := h.userRepo.CheckIfUserActive(ctx, parameters.Login)
	if err != nil {
		log.WithError(err).Info("DB request of active attribute returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !active {
		log.Debug("User is deactivated already")
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("User is deactivated already"))
		return
	}

	//deactivation request
	err = h.userRepo.DeactivateUser(ctx, parameters.Login)
	if err != nil {
		log.WithError(err).Info("DB user deactivation request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deactivated successfully"))
}
