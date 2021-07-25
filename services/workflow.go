package services

import (
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleUserInput(update tgbotapi.Update) {
	/* Debugging */
	utils.LogMessage(update)
	utils.LogUpdate(update)
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
			return
		case "/addplace",
			"/addplace@toGoListBot":
			utils.SendMessage(update, "Please enter the name of the place to begin")
			if err := utils.SetUserState(update, constants.SetName); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			return
		case "/query",
			"/query@toGoListBot":
			utils.ResetQuery(update)
			sendQuerySelectType(update, "How many places are you asking for?")
			if err := utils.SetUserState(update, constants.QuerySelectType); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			return
		case "/deleteplace",
			"/deleteplace@toGoListBot":
			sendPlacesToDeleteResponse(update, "Just select place do you want to delete?")
			if err := utils.SetUserState(update, constants.DeleteSelect); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/testredirect":

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

	/* Querying places */
	if constants.IsQuery(userState) {
		queryHandler(update, userState)
		return
	}

	/* Delete place */
	if constants.IsDeletePlace(userState) {
		deletePlaceHandler(update, userState)
		return
	}
}
