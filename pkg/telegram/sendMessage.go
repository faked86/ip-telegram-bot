package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) sendMessage(chatID int64, message string) {

	msg := tgbotapi.NewMessage(int64(chatID), message)

	for i := 0; i < 3; i++ {
		_, err := b.tgBot.Send(msg)
		if err != nil {
			log.Error(err)
			continue
		}
		break
	}
}
