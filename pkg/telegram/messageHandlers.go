package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	log.Printf("[%s] %s", message.From.UserName, message.Text)
	b.sendMessage(message.Chat.ID, message.Text)
}
