package handlers

import (
	"net/http"
)

type registrationRequest struct {
	first_name string `json: "first_name"`
	last_name  string `json: "last_name"`
	email      string `json: "login"`
	//social_media_links `json`
	login    string `json: "login"`
	password string `json:"password`
}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type loginRequest struct {
	login    string `json:"login"`
	password string `json:"passord"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
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
