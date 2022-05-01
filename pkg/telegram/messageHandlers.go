package telegram

import (
	"encoding/json"
	"fmt"
	"net"

	ipapi "github.com/faked86/ip-telegram-bot/pkg/ip-API"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handleCommandStart(message.Chat.ID)

		default:
			b.sendMessage(message.Chat.ID, "No such command.")
		}
		return
	}

	log.Printf("[%s] %s", message.From.UserName, message.Text)

	ip := message.Text

	var msg string

	if net.ParseIP(ip) == nil {

		msg = fmt.Sprintf("IP Address: %s - Invalid", ip)
		log.Info(msg)
		b.sendMessage(message.Chat.ID, msg)

	} else {
		msg = fmt.Sprintf("IP Address: %s - Valid\n", ip)
		log.Info(msg)

		apiResp, err := ipapi.IpInfo(ip)

		if err != nil {
			b.sendMessage(message.Chat.ID, "Failed to check this IP, may be it is on private range.")
			return
		}

		res, err := json.Marshal(apiResp)
		if err != nil {
			log.Error(err)
			b.sendMessage(message.Chat.ID, fmt.Sprint(err))
			return
		}

		strRes := string(res)
		b.sendMessage(message.Chat.ID, strRes)
	}
}

func (b *Bot) handleCommandStart(chatID int64) {
	b.sendMessage(chatID, "Hi I am IP checker bot. Send me IP address to get info about it.")
}
