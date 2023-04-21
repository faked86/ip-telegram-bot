package database

import (
	"github.com/faked86/ip-telegram-bot/pkg/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initiate(dbURL string) *gorm.DB {

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.IpInfo{}, &models.User{}, &models.Request{})

	return db
}
