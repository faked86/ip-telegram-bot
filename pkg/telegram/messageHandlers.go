package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	log.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	for i := 0; i < 3; i++ {
		_, err := b.tgBot.Send(msg)
		if err != nil {
			log.Error(err)
			continue
		}
		break
	}
}
