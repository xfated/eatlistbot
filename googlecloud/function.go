package function

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xfated/golistbot/services"
	"github.com/xfated/golistbot/services/utils"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func init() {
	utils.InitTelegram()
	utils.InitFirebase()
}

func TelegramHandler(w http.ResponseWriter, r *http.Request) {

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return
	}

	go services.HandleUserInput(&update)
}
