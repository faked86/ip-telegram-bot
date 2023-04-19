package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/faked86/ip-telegram-bot/pkg/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	q := r.URL.Query()
	idStr := q.Get("id")

	if idStr == "" {
		var users []models.User
		queryRes := s.db.Find(&users)
		if queryRes.Error != nil {
			log.Error(queryRes.Error)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, queryRes.Error)))
			return
		}

		res, err := json.MarshalIndent(users, "", "    ")
		if err != nil {
			log.Error(err)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
			return
		}

		w.Write(res)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	var users []models.User
	queryRes := s.db.Where("id = ?", id).Limit(1).Find(&users)
	if queryRes.Error != nil {
		log.Error(queryRes.Error)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, queryRes.Error)))
		return
	}

	if len(users) == 0 {
		w.Write([]byte(`{"error": "No such user in database"}`))
		return
	}

	user := users[0]

	res, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	w.Write(res)
}

func (s *Server) GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	q := r.URL.Query()
	idStr := q.Get("id")

	if idStr == "" {
		log.Error("/get_history_by_tg no id specified")
		w.Write([]byte(`{"error": "Enter user id like: /get_history_by_tg?id=1}"`))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	var results []joinRes
	s.db.Model(&models.Request{}).Select(
		`requests.id,
		ip_infos.ip,
		ip_infos.country,
		ip_infos.status,
		ip_infos.country_code,
		ip_infos.region,
		ip_infos.region_name,
		ip_infos.city,
		ip_infos.zip,
		ip_infos.lat,
		ip_infos.lon,
		ip_infos.timezone,
		ip_infos.isp,
		ip_infos.org,
		ip_infos.as`).Where("user_id = ?", id).Joins("left join ip_infos on requests.ip_info_ip = ip_infos.ip").Scan(&results)

	res, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	w.Write(res)

}

func (s *Server) DeleteFromHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	res := s.db.Delete(&models.Request{}, id)
	if res.Error != nil {
		log.Error(res.Error)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, res.Error)))
	}
}

type joinRes struct {
	Id          string  `json:"request_id"`
	IP          string  `json:"query"`
	Country     string  `json:"country"`
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
