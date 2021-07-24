package services

import (
	"errors"
	"log"
	"os"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	baseURL            = "https://togolist-bot.herokuapp.com/"
	bot                *tgbotapi.BotAPI
)

func InitTelegram() {
	var err error

	// Init bot
	bot, err = tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Println(err)
	}

	// Set webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(baseURL + bot.Token))
	if err != nil {
		log.Println("Problem setting Webhook", err.Error())
	}

	log.Println("Loaded telegram bot")
}

func SendMessage(msg tgbotapi.MessageConfig) {
	bot.Send(msg)
}

func SendStartInstructions(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Haha you sent start")
	bot.Send(msg)
	SetUserState(update, Idle)

	// Debug
	LogMessage(update)
}

func SendUnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command, please use /start for commands")
	bot.Send(msg)
}

func LogMessage(update tgbotapi.Update) {
	log.Printf("Message: %+v", update.Message)
}

func GetChatUserID(update tgbotapi.Update) (chatID, userID string, err error) {
	if update.Message == nil {
		chatID = ""
		userID = ""
		err = errors.New("invalid message")
		return
	}

	chatID = strconv.FormatInt(update.Message.Chat.ID, 10)
	userID = strconv.Itoa(update.Message.From.ID)
	err = nil
	return
}

func GetMessage(update tgbotapi.Update) (message string, messageID int, err error) {
	if update.Message == nil {
		message = ""
		messageID = 0
		err = errors.New("invalid message")
		return
	}

	message = update.Message.Text
	messageID = update.Message.MessageID
	err = nil
	return
}

/* Confirm details and add to list */
