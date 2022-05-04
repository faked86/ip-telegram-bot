package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	port   string
	db     *gorm.DB
	router *mux.Router
}

func NewServer(port string, db *gorm.DB) Server {
	return Server{
		port:   port,
		db:     db,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() {
	s.configureRouter()

	err := http.ListenAndServe(":"+s.port, s.router)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/get_users", s.handleGetUsers()).Methods("GET")
	s.router.HandleFunc("/get_history_by_tg", s.handleGetHistory()).Methods("GET")
	s.router.HandleFunc("/delete_from_history/{id}", s.handleDeleteFromHistory()).Methods("DELETE")
}
