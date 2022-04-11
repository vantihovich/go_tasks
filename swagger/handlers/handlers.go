package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type registrationRequest struct {
	firstName string `json: "first_name"`
	lastName  string `json: "last_name"`
	email     string `json: "login"`
	//social_media_links `json`
	login    string `json: "login"`
	password string `json:"password`
}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	//Below code is for testing purposes
	log.Println("register request came to func")
	log.Println("register request handeled successfully ")
	w.Write([]byte("Handeled the request "))
	//w.WriteHeader(http.StatusNotImplemented)// will be instead of testing code
}

type loginRequest struct {
	login    string `json:"login"`
	password string `json:"passord"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("login request came ")
	w.WriteHeader(http.StatusNotImplemented)

	// w.Header().Set("Content-Type", "application/json")
	// parameters := loginRequest{}

	// err := json.NewDecoder(r.Body).Decode(&parameters)

	// if err == io.EOF {
	// 	log.WithFields(log.Fields{"Error": err}).Info("Error empty request body")
	// 	fmt.Fprint(w, "Empty request body")
	// 	return
	// } else if err != nil {
	// 	log.WithFields(log.Fields{"Error": err}).Info("Error occurred")
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

}

func GracefulLoginHandler() {
	//possible place to hold graceful shutdown logic here
}

func GracefulRegisterHandler() {
	//possible place to hold graceful shutdown logic here

}
