package main

import (
	"fmt"
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
	db := database.Initiate()
	fmt.Println(db)

	b := telegram.NewBot(os.Getenv("TOKEN"), db)
	go b.Start()

	s := server.NewServer(os.Getenv("PORT"), db)
	s.Start()
}
