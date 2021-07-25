package services

import (
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendPlacesToDeleteResponse(update tgbotapi.Update, text string) {
	placeNames, err := utils.GetPlaceNames(update)
	if err != nil {
		log.Printf("error GetPlaceNames: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
	}

	/* Set each name as its own inline row */
	var nameButtons = make([][]tgbotapi.InlineKeyboardButton, len(placeNames))
	i := 0
	for name := range placeNames {
		nameButtons[i] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(name, name),
		)
		i++
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(nameButtons...)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func sendConfirmDeleteResponse(update tgbotapi.Update, text string) {
	// Create buttons
	yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	// Create rows
	row := tgbotapi.NewInlineKeyboardRow(yesButton, noButton)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func deletePlaceHandler(update tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.DeleteSelect:
		// Expect user to select from inline keyboard markup. (name of places to delete)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		name, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
		}
		utils.SetPlaceTarget(update, name)

	case constants.DeleteConfirm:
	}
}
