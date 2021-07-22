package services

import (
	"log"
	"os"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	baseURL            = "https://toeatlist-bot.herokuapp.com/"
	bot                *tgbotapi.BotAPI
)

func InitTelegram() {
	var err error

	// Init bot
	bot, err = tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Fatalln(err)
	}

	// Set webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(baseURL + bot.Token))
	if err != nil {
		log.Fatalln("Problem setting Webhook", err.Error())
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
}

func SendUnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command, please use /start for commands")
	bot.Send(msg)
}
