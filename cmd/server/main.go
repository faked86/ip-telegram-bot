package main

import (
	"fmt"
	"os"

	"github.com/faked86/ip-telegram-bot/pkg/database"
	"github.com/faked86/ip-telegram-bot/pkg/server"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgDBname := os.Getenv("POSTGRES_DB")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPass, pgHost, pgPort, pgDBname)
	db := database.Initiate(dbURL)

	s := server.NewServer(os.Getenv("PORT"), db)
	s.Start()
}
