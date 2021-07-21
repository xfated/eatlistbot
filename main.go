// Package p contains an HTTP Cloud Function.
package p

import (
	"log"
	"net/http"
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
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	init_bot()

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	port := os.Getenv("PORT")

	updates := bot.ListenForWebhook("/" + bot.Token)

	http.ListenAndServe(":"+port, nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
