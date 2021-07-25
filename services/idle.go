package services

import (
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func idleHandler(update tgbotapi.Update) {
	switch update.Message.Text {
	case "/addplace",
		"/addplace@toGoListBot":
		if err := utils.SetUserState(update, constants.SetName); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
		utils.SendMessage(update, "Please enter the name of the place to begin")
	}
}
