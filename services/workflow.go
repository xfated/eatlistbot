package services

import (
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleUserInput(update tgbotapi.Update) {
	/* Debugging */
	// utils.LogMessage(update)
	// utils.LogUpdate(update)
	utils.LogCallbackQuery(update)

	/* Check for main commands */
	message, _, err := utils.GetMessage(update)
	if err == nil {
		switch message {
		case "/start",
			"/start@toGoListBot",
			"/reset",
			"/reset@toGoListBot":
			utils.RemoveMarkupKeyboard(update, "I am ready!")
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/query",
			"/query@toGoListBot":
			sendQuerySelectType(update, "How many places are you asking for?")
			if err := utils.SetUserState(update, constants.QuerySelectType); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			return
		}
	}

	/* Get user state for Targeted handling */
	userState, err := utils.GetUserState(update)
	if err != nil {
		log.Printf("error getting user state: %+v", err)
		return
	}

	/* Idle state */
	if userState == constants.Idle {
		idleHandler(update)
		return
	}

	/* Adding new place */
	if constants.IsAddingNewPlace(userState) {
		addPlaceHandler(update, userState)
		return
	}

	if constants.IsQuery(userState) {
		queryHandler(update, userState)
		return
	}
}
