package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	cnfg "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
	"github.com/vantihovich/go_tasks/tree/master/swagger/email"
	"github.com/vantihovich/go_tasks/tree/master/swagger/models"
	secretKey "github.com/vantihovich/go_tasks/tree/master/swagger/secretKey"
	"github.com/vantihovich/go_tasks/tree/master/swagger/validators"
)

type UsersHandler struct {
	userRepo          models.UserRepository
	cache             Cache
	jwtParam          string
	cfgLogin          cnfg.LoginLimitParameters
	mailClient        email.Client
	worldCoinIndexKey cnfg.WorldCoinIndexKey
}

type Cache interface {
	Get(string) (int, error)
	Set(string, int, int) error
}

func NewUsersHandler(userRepo models.UserRepository, c Cache, jwtParam string, cfgLogin cnfg.LoginLimitParameters, mailCli email.Client, wciKey cnfg.WorldCoinIndexKey) *UsersHandler {
	return &UsersHandler{
		userRepo:          userRepo,
		cache:             c,
		jwtParam:          jwtParam,
		cfgLogin:          cfgLogin,
		mailClient:        mailCli,
		worldCoinIndexKey: wciKey,
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

	//hashing password
	hash, err := bcrypt.GenerateFromPassword([]byte(parameters.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Info("error occurred when hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//an attempt to add new user
	err = h.userRepo.AddNewUser(ctx, parameters.Login, string(hash), parameters.FirstName, parameters.LastName, parameters.Email, parameters.SocialMediaLinks)
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

func (h *UsersHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	parameters := loginRequest{}
	response := loginResponse{}
	ttl := int(h.cfgLogin.InvalidLoginAttemptTTL.Seconds())
	infoTtl := int(h.cfgLogin.InvalidLoginAttemptTTL.Minutes())

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

	//trying to find user by provided login
	user, err := h.userRepo.FindByLogin(ctx, parameters.Login)
	if err != nil {
		if err == ErrNoRows {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Please check user login"))
			return
		}
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//if user found , checking if user`s profile is banned due to invalid password inputs
	invalidLogins, err := h.cache.Get(parameters.Login)
	if err != nil {
		log.WithError(err).Info("error getting the cache from Redis")
	}

	if invalidLogins >= h.cfgLogin.MaxAllowedInvalidLogins {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "you`ve reached the login attempts limit, please try again after %d minutes", infoTtl)
		return
	}

	//if user login found and profile is not banned, comparing provided password with hash from DB
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(parameters.Password))
	if errPassword != nil {
		//increasing the amount of password input attempts
		invalidLogins++
		err = h.cache.Set(parameters.Login, invalidLogins, ttl)
		if err != nil {
			log.WithError(err).Info("error setting the cache to Redis")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Please check user password"))
		return
	}

	//provided password is correct, so checking if user is deactivated
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
	response.Token, err = jwToken.SignedString([]byte(h.jwtParam))
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

	deactivatorID := (r.Context().Value(ContextKeyUserID))
	if deactivatorID == nil {
		log.Error("empty authorization context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//check if logged in user is allowed to deactivate the user from the request
	deactivator, err := h.userRepo.GetAdminAttrUserLogin(ctx, deactivatorID.(int))
	if err != nil {
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !deactivator.IsAdmin() && (deactivator.Login != parameters.Login) {
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

type resetPasswordRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

func (h *UsersHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parameters := resetPasswordRequest{}
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

	pwdChangerID := r.Context().Value(ContextKeyUserID)
	if pwdChangerID == nil {
		log.Error("empty authorization context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//confirming old password
	user, err := h.userRepo.FindByID(ctx, pwdChangerID.(int))
	if err != nil {
		if err == ErrNoRows {
			log.WithError(err).Error("no rows found")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Please check user login"))
			return
		}
		log.WithError(err).Info("DB request returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(parameters.OldPassword)); err != nil {
		log.WithError(err).Error("hashes don't match")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Please check your password"))
		return
	}

	//validating and hashing new password
	valid := validators.ValidateChangePasswordRequest(parameters.NewPassword, parameters.ConfirmNewPassword)
	if !valid {
		log.Error("new password or new confirmed password are empty or not equal")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password was not changed: new password or new confirmed password are empty or not equal"))
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(parameters.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Info("error occurred when hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//PWD change request
	err = h.userRepo.ChangePassword(ctx, pwdChangerID.(int), string(newHash))
	if err != nil {
		log.WithError(err).Info("error occurred when resetting user's password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug("password changed succesfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password is changed"))
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *UsersHandler) ForgotPasswordSendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parameters := forgotPasswordRequest{}
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
	//generate random secret to store in DB for recovering Password
	secret := secretKey.NewRandomString(10)

	//write the secret to DB
	exists, err := h.userRepo.WriteSecretToDBForPasswordRecovery(ctx, parameters.Email, secret)
	if err != nil {
		log.WithError(err).Info("DB request of writing secret returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Info("email does not exist in the system")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please specify the correct email address"))
		return
	}
	//send email to provided address with the link containing secret
	err = h.mailClient.SendForgottenPasswordEmail(parameters.Email, secret)
	if err != nil {
		log.WithError(err).Info("An attempt to send email to user`s email returned error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The email with secret key has been sent"))
}

type forgotPasswordResetPassword struct {
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

func (h *UsersHandler) ForgotPasswordResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parameters := forgotPasswordResetPassword{}
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

	secretParameter := r.URL.Query().Get("parameter")
	if secretParameter == "" {
		log.WithError(err).Info("parameter was not specified in the request URL")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//validating and hashing new password
	valid := validators.ValidateChangePasswordRequest(parameters.NewPassword, parameters.ConfirmNewPassword)
	if !valid {
		log.Error("new password or new confirmed password are empty or not equal")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password was not changed: new password or new confirmed password are empty or not equal"))
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(parameters.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Info("error occurred when hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//the request to DB to update password if provided in URL secret matches, also this request clears recovery field
	err = h.userRepo.CheckSecretChangePassword(ctx, secretParameter, string(newHash))
	if err != nil {
		if err == ErrNoRows {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Please check provided secret, changing password is forbidden"))
			return
		}
		log.WithError(err).Info("error occurred when resetting user's password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug("password changed succesfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password is changed"))
}

type worldCoinIndexTickers struct {
	Label      string  `json:"label"`
	Name       string  `json:"name"`
	Price      float32 `json:"price"`
	Volume_24h float32 `json:"volume_24h"`
	Timestamp  int64   `json:"timestamp"`
}
type worldCoinIndexTickersAPIResponse struct {
	Markets []worldCoinIndexTickers
}
type worldCoinIndexTickersRequest struct {
	Label []string `json:"label"`
	Fiat  string   `json:"fiat"`
}

type response struct {
	Label      string
	Name       string
	Price      float32
	Volume_24h float32
	Timestamp  time.Time
}

type responseArray struct {
	Response []response
}

func (h *UsersHandler) WorldCoinIndexTickers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := worldCoinIndexTickersRequest{}
	apiResponse := worldCoinIndexTickersAPIResponse{}
	resp := response{}
	respArray := responseArray{}

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

	//validating request parameters
	valid := validators.ValidateWorldCoinIndexRequest(parameters.Label, parameters.Fiat)
	if !valid {
		log.Error("empty list of cryptocurrencies or fiat")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please check if there is at least 1 cryptocurrency ticker and fiat"))
		return
	}

	//concatenating values from labels array in the request to use one string in url
	var list string
	for i, val := range parameters.Label {
		if i < (len(parameters.Label) - 1) {
			list += val + "-"
		} else {
			list += val
		}
	}

	fiat := parameters.Fiat
	url := "https://www.worldcoinindex.com/apiservice/ticker?key=" + h.worldCoinIndexKey.Key + "&label=" + list + "&fiat=" + fiat

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithError(err).Info("error occurred when creating request to WorldCoinIndex api")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	apiResp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Info("error occurred when performing request to WorldCoinIndex api")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(apiResp.Body).Decode(&apiResponse)
	if err != nil {
		if err == io.EOF {
			log.WithError(err).Info("empty response body")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.WithError(err).Info("error decoding response body occurred")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//converting timestamp from unix to date in the response
	for _, val := range apiResponse.Markets {
		resp.Label = val.Label
		resp.Name = val.Name
		resp.Price = val.Price
		resp.Timestamp = time.Unix(val.Timestamp, 0)
		resp.Volume_24h = val.Volume_24h
		respArray.Response = append(respArray.Response, resp)
	}
	//writing the response to user
	err = json.NewEncoder(w).Encode(&respArray)
	if err != nil {
		log.WithError(err).Info("error encoding struct to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
