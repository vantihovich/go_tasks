package wci

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func SendWCIRequest(req *http.Request) (resp *http.Response, err error) {
	client := &http.Client{}
	apiResp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Info("error occurred when performing request to WorldCoinIndex api")
		return apiResp, err
	}
	return apiResp, nil
}

func NewWCIRequest(requestList []string, fiat, key, urlWci string) (req *http.Request, err error) {
	var list string
	for i, val := range requestList {
		if i < (len(requestList) - 1) {
			list += val + "-"
		} else {
			list += val
		}
	}

	url := urlWci + "key=" + key + "&label=" + list + "&fiat=" + fiat

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithError(err).Info("error occurred when creating request to WorldCoinIndex api")
		return req, err
	}

	return req, nil
}
