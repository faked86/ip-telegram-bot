package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Bot struct {
	tgBot *tgbotapi.BotAPI
	db    *gorm.DB
}

func NewBot(token string, db *gorm.DB) Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	res := Bot{
		tgBot: bot,
		db:    db,
	}

	return res
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tgBot.GetUpdatesChan(u)
	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {

		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}
