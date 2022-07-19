package rateshandlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	cnfg "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
	"github.com/vantihovich/go_tasks/tree/master/swagger/validators"
	"github.com/vantihovich/go_tasks/tree/master/swagger/wci"
)

type WCIHandler struct {
	worldCoinIndexCfg cnfg.WorldCoinIndexParameters
}

func NewWCIHandler(wci cnfg.WorldCoinIndexParameters) *WCIHandler {
	return &WCIHandler{
		worldCoinIndexCfg: wci,
	}
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
type WorldCoinIndexTickersRequest struct {
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

func (wh *WCIHandler) WorldCoinIndexTickers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := WorldCoinIndexTickersRequest{}
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

	request, err := wci.NewWCIRequest(parameters.Label, parameters.Fiat, wh.worldCoinIndexCfg.Key, wh.worldCoinIndexCfg.URL)
	if err != nil {
		log.WithError(err).Info("error occurred when creating request to API")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiResp, err := wci.SendWCIRequest(request)
	if err != nil {
		log.WithError(err).Info("error occurred when performing request to API")
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
