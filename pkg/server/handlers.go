package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/faked86/ip-telegram-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

func (s *Server) handleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		idStr := q.Get("id")

		if idStr == "" {
			var users []models.User
			queryRes := s.db.Find(&users)
			if queryRes.Error != nil {
				log.Error(queryRes.Error)
				w.Write([]byte(fmt.Sprint(queryRes.Error)))
				return
			}

			res, err := json.MarshalIndent(users, "", "    ")
			if err != nil {
				log.Error(err)
				w.Write([]byte(fmt.Sprint(err)))
				return
			}

			w.Write(res)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error(err)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		var users []models.User
		queryRes := s.db.Where("id = ?", id).Limit(1).Find(&users)
		if queryRes.Error != nil {
			log.Error(queryRes.Error)
			w.Write([]byte(fmt.Sprint(queryRes.Error)))
			return
		}

		if len(users) == 0 {
			w.Write([]byte("No such user in database."))
			return
		}

		user := users[0]

		res, err := json.MarshalIndent(user, "", "    ")
		if err != nil {
			log.Error(err)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		w.Write(res)
	}
}
