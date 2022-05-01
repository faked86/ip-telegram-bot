package ipapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/faked86/ip-telegram-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

const baseApiUrl = "http://ip-api.com/json/"

func IpInfo(ip string) (apiResp *models.IpInfo, err error) {

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

	if apiResp.Status == "fail" {
		err := errors.New("failed to check this IP, may be it is on private range")
		log.Error(err)
		return nil, err
	}

	return apiResp, nil
}
