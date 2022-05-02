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
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handleCommandStart(message.Chat.ID)

		default:
			b.sendMessage(message.Chat.ID, "No such command.")
		}
		return
	}

	ip := message.Text

	if net.ParseIP(ip) != nil {
		b.handleValidIp(message.Chat.ID, ip)
	} else {
		msg := fmt.Sprintf("IP Address: %s - Invalid", ip)
		log.Info(msg)
		b.sendMessage(message.Chat.ID, msg)
	}
}

func (b *Bot) handleCommandStart(chatID int64) {
	log.Info("Start command")
	b.sendMessage(chatID, "Hi I am IP checker bot. Send me IP address to get info about it.")
}

func (b *Bot) handleValidIp(chatID int64, ip string) {
	msg := fmt.Sprintf("IP Address: %s - Valid\n", ip)
	log.Info(msg)

	apiResp, err := ipapi.IpInfo(ip)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprint(err))
		return
	}

	res, err := json.MarshalIndent(apiResp, "", "    ")
	if err != nil {
		log.Error(err)
		b.sendMessage(chatID, fmt.Sprint(err))
		return
	}

	strRes := string(res)
	b.sendMessage(chatID, strRes)
}
