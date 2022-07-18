package wci

import (
	"net/http"
	"net/url"
	"strings"

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

func NewWCIRequest(requestList []string, fiat, key, urlTemplate string) (req *http.Request, err error) {
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

	req, err = http.NewRequest(http.MethodGet, urlWCI.String(), nil)
	if err != nil {
		log.WithError(err).Info("error occurred when creating request to WorldCoinIndex api")
		return req, err
	}

	return req, nil
}
