package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/vantihovich/go_tasks/tree/master/swagger/models"
	"github.com/vantihovich/go_tasks/tree/master/swagger/validators"
)

type WCIHandler struct {
	wciClient models.WCISource
}

func NewWCIHandler(wciCli models.WCISource) *WCIHandler {
	return &WCIHandler{
		wciClient: wciCli,
	}
}

type WorldCoinIndexTickersRequest struct {
	Label []string `json:"label"`
	Fiat  string   `json:"fiat"`
}

func (wh *WCIHandler) WorldCoinIndexTickers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := WorldCoinIndexTickersRequest{}

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

	tickers, err := wh.wciClient.LoadTickers(parameters.Label, parameters.Fiat)
	if err != nil {
		log.WithError(err).Info("error occurred when sending request to API")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&tickers)
	if err != nil {
		log.WithError(err).Info("error encoding struct to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
