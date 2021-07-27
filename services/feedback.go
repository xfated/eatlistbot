package services

import (
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
	}
}
