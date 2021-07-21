// Package p contains an HTTP Cloud Function.
package p

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	baseURL            = "https://toeatlist-bot.herokuapp.com/"
	bot                *tgbotapi.BotAPI
)

func initTelegram() {
	var err error
	// Init bot
	bot, err = tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	// Set webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(baseURL + bot.Token))
	if err != nil {
		log.Fatalln("Problem setting Webhook", err.Error())
	}
}

func fetchUpdates(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	updates := bot.ListenForWebhook("/" + bot.Token)
	return updates
}

func main() {
	initTelegram()

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	port := os.Getenv("PORT")

	updates := fetchUpdates(bot)

	fmt.Println(updates)

	go http.ListenAndServe(":"+port, nil)
}
