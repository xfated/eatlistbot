package services

import (
	"log"
	"os"

	"github.com/xfated/eatlistbot/firebase"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	baseURL            = "https://toeatlist-bot.herokuapp.com/"
	bot                *tgbotapi.BotAPI
)

// Finite state machine for handling adding items
type State int

const (
	Idle State = iota
	SetName
	SetAddress
	SetURL
	SetImages
	SetTags
	Finished
)

func getNextState(cur State) State {
	switch cur {
	case Idle, Finished:
		return cur
	case SetName, SetAddress, SetURL, SetImages, SetTags:
		return cur + 1
	default:
		return cur
	}
}

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
	firebase.SetUserState(update, Idle)
}

func SendUnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command, please use /start for commands")
	bot.Send(msg)
}
