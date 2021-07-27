// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xfated/golistbot/services"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func TelegramHandler(w http.ResponseWriter, r *http.Request) {

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return
	}
	return

	go services.HandleUserInput(&update)
}
