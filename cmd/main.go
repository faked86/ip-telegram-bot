package main

import (
	"os"

	"github.com/faked86/ip-telegram-bot/pkg/database"
	"github.com/faked86/ip-telegram-bot/pkg/server"
	"github.com/faked86/ip-telegram-bot/pkg/telegram"

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
	dbURL := os.Getenv("PG_ADDRESS")
	db := database.Initiate(dbURL)

	b := telegram.NewBot(os.Getenv("TOKEN"), db)
	go b.Start()

	s := server.NewServer(os.Getenv("PORT"), db)
	s.Start()
}
