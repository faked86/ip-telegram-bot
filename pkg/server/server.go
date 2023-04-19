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

	log.Printf("HTTP server starting at port: %s", s.port)

	err := http.ListenAndServe(":"+s.port, s.router)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/users", s.GetUsers).Methods("GET")
	s.router.HandleFunc("/users/{id}", s.GetUser).Methods("GET")
	s.router.HandleFunc("/users/{id}/history", s.GetUserHistory).Methods("GET")
	s.router.HandleFunc("/requests/{id}", s.DeleteRequest).Methods("DELETE")
}
