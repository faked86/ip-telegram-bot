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
	var outMsg string

	if message.IsCommand() {

		switch message.Command() {
		case "start":
			outMsg = b.handleCommandStart(*message)

		case "unique":
			outMsg = b.handleCommandUnique(message.From.ID)

		default:
			outMsg = "No such command."
		}
		b.sendMessage(message.Chat.ID, outMsg)
		return
	}

	ip := message.Text

	if net.ParseIP(ip) != nil {
		msg := fmt.Sprintf("IP Address: %s - Valid\n", ip)
		log.Info(msg)
		outMsg = b.handleValidIp(*message, ip)
	} else {
		outMsg = fmt.Sprintf("IP Address: %s - Invalid", ip)
		log.Info(outMsg)
	}
	b.sendMessage(message.Chat.ID, outMsg)
}

func (b *Bot) handleCommandStart(message tgbotapi.Message) string {
	log.Info("Start command")

	user := models.User{
		ID:       message.From.ID,
		UserName: message.From.UserName,
		Admin:    false,
	}

	errMsg := ""
	if res := b.db.FirstOrCreate(&user, user); res.Error != nil {
		log.Error(res.Error)
		errMsg = "Failed try to register you in our database. History will be unavailable."
		// b.sendMessage(message.From.ID, "Failed try to register you in our database. History will be unavailable.")
	}

	log.Println("DB User:", user)

	// b.sendMessage(message.From.ID, "Hi I am IP checker bot. Send me IP address to get info about it.")
	return "Hi I am IP checker bot. Send me IP address to get info about it." + errMsg
}

func (b *Bot) handleValidIp(message tgbotapi.Message, ip string) string {

	var apiResp *models.IpInfo
	if dbRes := b.db.Where("ip = ?", ip).FirstOrCreate(&apiResp, models.IpInfo{IP: ip}); dbRes.Error != nil {
		log.Error(dbRes.Error)
		// b.sendMessage(message.From.ID, fmt.Sprint(dbRes.Error))
		return fmt.Sprint(dbRes.Error)
	}

	if apiResp.Status == "" {
		log.Print("Ip not from db")
		resp, err := ipapi.IpInfo(ip)
		if err != nil {
			// b.sendMessage(message.From.ID, fmt.Sprint(err))
			return fmt.Sprint(err)
		}
		b.db.Model(apiResp).Updates(resp)
	} else {
		log.Print("Ip from db")
	}

	res, err := json.MarshalIndent(apiResp, "", "    ")
	if err != nil {
		log.Error(err)
		// b.sendMessage(message.From.ID, fmt.Sprint(err))
		return fmt.Sprint(err)
	}

	errMsg := ""
	dbResReq := b.db.Create(&models.Request{UserID: message.From.ID, IpInfoIP: ip})
	if dbResReq.Error != nil {
		log.Error(dbResReq.Error)
		// b.sendMessage(message.From.ID, "Failed to save request to database.")
		errMsg = " Failed to save request to database."
	}

	strRes := string(res)
	// b.sendMessage(message.From.ID, strRes)

	return strRes + errMsg
}

func (b *Bot) handleCommandUnique(userID int64) string {
	var reqs []models.Request
	res := b.db.Select("DISTINCT ip_info_ip").Where("user_id = ?", userID).Find(&reqs)
	if res.Error != nil {
		log.Error(res.Error)
		// b.sendMessage(userID, "Something wrong in unique function.")
		return "Something wrong in unique function."
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
			// b.sendMessage(userID, fmt.Sprint(err))
			return fmt.Sprint(err)
		}
		msg = msg + "\n" + string(res)
	}
	// b.sendMessage(userID, msg)
	return msg
}
