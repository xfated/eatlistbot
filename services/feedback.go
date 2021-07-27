package services

import (
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func feedbackHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.Feedback:
		utils.AddFeedback(update)
		utils.SendToFeedbackChat(update)
		utils.SendMessage(update, "Thank you for your feedback!\nIt has been well received", false)
		if err := utils.SetUserState(update, constants.Idle); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			break
		}
	}
}
