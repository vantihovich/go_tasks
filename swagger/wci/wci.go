package wci

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	config "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
	"github.com/vantihovich/go_tasks/tree/master/swagger/models"
)

type WCIClient struct {
	client http.Client
	cfg    config.WorldCoinIndexParameters
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

func New(cfg config.WorldCoinIndexParameters) (wciCli WCIClient) {
	wciCli.client = http.Client{}
	wciCli.cfg = cfg

	return wciCli
}

func (w *WCIClient) LoadTickers(requestList []string, fiat string) (response models.WorldCoinIndexHandlerResponseArray, err error) {
	apiResponse := worldCoinIndexTickersAPIResponse{}

	request, err := newWCIRequest(requestList, fiat, w.cfg.Key, w.cfg.URL)
	if err != nil {
		log.WithError(err).Info("error occurred when creating new WCI request")
		return response, err
	}

	apiResp, err := w.client.Do(request)
	if err != nil {
		log.WithError(err).Info("error occurred when performing request to WorldCoinIndex api")
		return response, err
	}

	err = json.NewDecoder(apiResp.Body).Decode(&apiResponse)
	if err != nil {
		if err == io.EOF {
			log.WithError(err).Info("empty response body")
			return response, err
		}
		log.WithError(err).Info("error decoding response body occurred")
		return response, err
	}
	//converting timestamp from unix to date in the response
	response = convertTimeStamp(apiResponse)

	return response, nil
}

func newWCIRequest(requestList []string, fiat, key, urlTemplate string) (request *http.Request, err error) {
	list := strings.Join(requestList, "-")

	urlWCI, err := url.Parse(urlTemplate)
	if err != nil {
		log.WithError(err).Info("error occurred when parsing WCI url template")
		return
	}

	values := urlWCI.Query()
	values.Set("key", key)
	values.Set("label", list)
	values.Set("fiat", fiat)

	urlWCI.RawQuery = values.Encode()

	request, err = http.NewRequest(http.MethodGet, urlWCI.String(), nil)
	if err != nil {
		log.WithError(err).Info("error occurred when creating request to WorldCoinIndex api")
		return request, err
	}

	return request, nil
}

func convertTimeStamp(apiResponse worldCoinIndexTickersAPIResponse) models.WorldCoinIndexHandlerResponseArray {
	resp := models.WorldCoinIndexHandlerResponse{}
	respArray := models.WorldCoinIndexHandlerResponseArray{}
	for _, val := range apiResponse.Markets {
		resp.Label = val.Label
		resp.Name = val.Name
		resp.Price = val.Price
		resp.Timestamp = time.Unix(val.Timestamp, 0)
		resp.Volume_24h = val.Volume_24h
		respArray.Response = append(respArray.Response, resp)
	}
	return respArray
}
