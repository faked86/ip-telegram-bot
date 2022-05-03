package telegram

import (
	"encoding/json"
	"fmt"
	"net"

	ipapi "github.com/faked86/ip-telegram-bot/pkg/ip-API"
	"github.com/faked86/ip-telegram-bot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handleCommandStart(*message)

		case "unique":
			b.handleCommandUnique(message.From.ID)

		default:
			b.sendMessage(message.Chat.ID, "No such command.")
		}
		return
	}

	ip := message.Text

	if net.ParseIP(ip) != nil {
		msg := fmt.Sprintf("IP Address: %s - Valid\n", ip)
		log.Info(msg)
		b.handleValidIp(*message, ip)
	} else {
		msg := fmt.Sprintf("IP Address: %s - Invalid", ip)
		log.Info(msg)
		b.sendMessage(message.Chat.ID, msg)
	}
}

func (b *Bot) handleCommandStart(message tgbotapi.Message) {
	log.Info("Start command")

	user := models.User{
		ID:       message.From.ID,
		UserName: message.From.UserName,
		Admin:    false,
	}

	if res := b.db.FirstOrCreate(&user, user); res.Error != nil {
		log.Error(res.Error)
		b.sendMessage(message.From.ID, "Failed try to register you in our database. History will be unavailable.")
	}

	log.Println("DB User:", user)

	b.sendMessage(message.From.ID, "Hi I am IP checker bot. Send me IP address to get info about it.")
}

func (b *Bot) handleValidIp(message tgbotapi.Message, ip string) {

	var apiResp *models.IpInfo
	if dbRes := b.db.Where("ip = ?", ip).FirstOrCreate(&apiResp, models.IpInfo{IP: ip}); dbRes.Error != nil {
		log.Error(dbRes.Error)
		b.sendMessage(message.From.ID, fmt.Sprint(dbRes.Error))
		return
	}

	if apiResp.Status == "" {
		log.Print("Ip not from db")
		resp, err := ipapi.IpInfo(ip)
		if err != nil {
			b.sendMessage(message.From.ID, fmt.Sprint(err))
			return
		}
		b.db.Model(apiResp).Updates(resp)
	} else {
		log.Print("Ip from db")
	}

	res, err := json.MarshalIndent(apiResp, "", "    ")
	if err != nil {
		log.Error(err)
		b.sendMessage(message.From.ID, fmt.Sprint(err))
		return
	}

	strRes := string(res)
	b.sendMessage(message.From.ID, strRes)

	dbResReq := b.db.Create(&models.Request{UserID: message.From.ID, IpInfoIP: ip})
	if dbResReq.Error != nil {
		log.Error(dbResReq.Error)
		b.sendMessage(message.From.ID, "Failed to save request to database.")
	}
}

func (b *Bot) handleCommandUnique(userID int64) {
	var reqs []models.Request
	res := b.db.Select("DISTINCT ip_info_ip").Where("user_id = ?", userID).Find(&reqs)
	if res.Error != nil {
		log.Error(res.Error)
		b.sendMessage(userID, "Something wrong in unique function.")
		return
	}

	ips := make([]string, len(reqs))
	for i, req := range reqs {
		ips[i] = req.IpInfoIP
	}

	var ipInfos []models.IpInfo
	b.db.Where("ip IN ?", ips).Find(&ipInfos)

	msg := "Unique results:"
	for _, info := range ipInfos {
		res, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			log.Error(err)
			b.sendMessage(userID, fmt.Sprint(err))
			return
		}
		msg = msg + "\n" + string(res)
	}
	b.sendMessage(userID, msg)
}
