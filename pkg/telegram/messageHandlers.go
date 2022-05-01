package telegram

import (
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
	b.sendMessage(message.Chat.ID, message.Text)
}

func (b *Bot) handleCommandStart(chatID int64) {

	b.sendMessage(chatID, "Hi I am IP checker bot. Send me IP address to get info about it.")
}
