package services

import (
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func idleHandler(update *tgbotapi.Update) {
	message, _, err := utils.GetMessage(update)
	if err != nil {
		return
	}
	switch message {
	default:
	}
}
