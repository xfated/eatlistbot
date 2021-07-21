// Package p contains an HTTP Cloud Function.
package p

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
var bot *tgbotapi.BotAPI

func init_bot() {
	var err error
	// Init bot
	bot, err = tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	// Set webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://toeatlist-bot.herokuapp.com/" + bot.Token))
}

func main() {
	init_bot()

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
