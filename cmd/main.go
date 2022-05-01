package main

import (
	"fmt"
	// "os"
	// "crypto/sha1"
	// "encoding/hex"

	"github.com/faked86/ip-telegram-bot/pkg/database"
	ipapi "github.com/faked86/ip-telegram-bot/pkg/ip-API"
	// "github.com/faked86/ip-telegram-bot/pkg/telegram"

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

	// b := telegram.NewBot(os.Getenv("TOKEN"))
	// b.Start()

	// server.Start

	resp, _ := ipapi.IpInfo("24.48.0.1")
	fmt.Println(*resp)
}
