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
			outMsg = b.handleCommandStart(message)

		case "unique":
			outMsg = b.handleCommandUnique(message.From.ID)

		case "spam":
			outMsg = b.handleCommandSpam(message)

		case "admin":
			outMsg = b.handleCommandAdmin(message)

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
		outMsg = b.handleValidIp(message, ip)
	} else {
		outMsg = fmt.Sprintf("IP Address: %s - Invalid", ip)
		log.Info(outMsg)
	}
	b.sendMessage(message.Chat.ID, outMsg)
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) string {
	log.Info("Start command")

	user := models.User{
		ID:       message.From.ID,
		Username: message.From.UserName,
		Admin:    false,
	}

	errMsg := ""
	if res := b.db.FirstOrCreate(&user, user); res.Error != nil {
		log.Error(res.Error)
		errMsg = "Failed try to register you in our database. History will be unavailable."
	}

	log.Println("DB User:", user)

	return "Hi I am IP checker bot. Send me IP address to get info about it." + errMsg
}

func (b *Bot) handleValidIp(message *tgbotapi.Message, ip string) string {

	var apiResp *models.IpInfo
	if dbRes := b.db.Where("ip = ?", ip).FirstOrCreate(&apiResp, models.IpInfo{IP: ip}); dbRes.Error != nil {
		log.Error(dbRes.Error)
		return fmt.Sprint(dbRes.Error)
	}

	if apiResp.Status == "" {
		log.Print("Ip not from db")
		resp, err := ipapi.IpInfo(ip)
		if err != nil {
			return fmt.Sprint(err)
		}
		b.db.Model(apiResp).Updates(resp)
	} else {
		log.Print("Ip from db")
	}

	res, err := json.MarshalIndent(apiResp, "", "    ")
	if err != nil {
		log.Error(err)
		return fmt.Sprint(err)
	}

	errMsg := ""
	dbResReq := b.db.Create(&models.Request{UserID: message.From.ID, IpInfoIP: ip})
	if dbResReq.Error != nil {
		log.Error(dbResReq.Error)
		errMsg = " Failed to save request to database."
	}

	strRes := string(res)

	return strRes + errMsg
}

func (b *Bot) handleCommandUnique(userID int64) string {
	var reqs []models.Request
	res := b.db.Select("DISTINCT ip_info_ip").Where("user_id = ?", userID).Find(&reqs)
	if res.Error != nil {
		log.Error(res.Error)
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
			return fmt.Sprint(err)
		}
		msg = msg + "\n" + string(res)
	}
	return msg
}

func (b *Bot) handleCommandSpam(message *tgbotapi.Message) string {
	var user models.User
	qRes := b.db.Where("id = ?", message.From.ID).First(&user)
	if qRes.Error != nil {
		log.Error(qRes.Error)
		return "Something wrong with query to database."
	}

	if !user.Admin {
		log.Printf("Not admin trying admin command (%s)", user.Username)
		return "You should be admin to use this command."
	}

	msg := message.CommandArguments()
	if msg == "" {
		log.Printf("Wrong /spam usage by %s", user.Username)
		return "Command format: '/spam <message>' e.g. '/spam Hello!'"
	}

	var allUsers []models.User
	res := b.db.Find(&allUsers)
	if res.Error != nil {
		log.Error(qRes.Error)
		return "Something wrong with query to database."
	}

	log.Printf("%s initiated mass spam to all our users [%d]", user.Username, len(allUsers))
	for _, target := range allUsers {
		b.sendMessage(target.ID, msg)
	}

	return "Done."
}

func (b *Bot) handleCommandAdmin(message *tgbotapi.Message) string {
	var user models.User
	qRes := b.db.Where("id = ?", message.From.ID).First(&user)
	if qRes.Error != nil {
		log.Error(qRes.Error)
		return "Something wrong with query to database."
	}

	if !user.Admin {
		log.Printf("Not admin trying admin command (%s)", user.Username)
		return "You should be admin to use this command."
	}

	username := message.CommandArguments()
	if username == "" {
		log.Printf("Wrong /admin usage by %s", user.Username)
		return "Command format: '/admin <username>' e.g. '/admin user1'"
	}

	var targetUsers []models.User
	qRes = b.db.Where("username = ?", username).Find(&targetUsers)
	if qRes.Error != nil {
		log.Error(qRes.Error)
		return "Something wrong with query to database."
	}

	if len(targetUsers) == 0 {
		log.Printf("Wrong /admin usage by %s", user.Username)
		return "No such user in my database."
	}

	targetUser := targetUsers[0]
	if !targetUser.Admin {
		targetUser.Admin = true
		b.db.Save(&targetUser)
		log.Printf("%s made %s admin", user.Username, targetUser.Username)
		b.sendMessage(targetUser.ID, "You are admin now.")
		return fmt.Sprintf("User %s is now admin.", targetUser.Username)
	}

	targetUser.Admin = false
	b.db.Save(&targetUser)
	b.sendMessage(targetUser.ID, "You are no longer admin.")
	log.Printf("%s made %s no admin", user.Username, targetUser.Username)
	return fmt.Sprintf("User %s is no longer admin.", targetUser.Username)
}
