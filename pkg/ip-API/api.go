package ipapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const baseApiUrl = "http://ip-api.com/json/"

type ApiResponse struct {
	Country     string  `json:"country"`
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

func IpInfo(ip string) (apiResp *ApiResponse, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseApiUrl, ip))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Error(err)
		return nil, err
	}

	return apiResp, nil
}
