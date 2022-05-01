package main

import (
	"os"

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
	b := telegram.NewBot(os.Getenv("TOKEN"))
	b.Start()

	// server.Start
}
