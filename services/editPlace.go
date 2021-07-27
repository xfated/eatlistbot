package services

import (
	"log"
	"strconv"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendPlacesToEditResponse(update *tgbotapi.Update, text string) {
	// Get target chat
	chatID, err := utils.GetChatTarget(update)
	if err != nil {
		log.Printf("error GetChatTarget: %+v", err)
		utils.SendMessage(update, "Sorry an error occured!")
		return
	}
	chatIDString := strconv.FormatInt(chatID, 10)

	// Extract place names
	placeNames, err := utils.GetPlaceNames(update, chatIDString)
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
	msg := utils.SendInlineKeyboard(update, text, inlineKeyboard)
	utils.AddMessageToDelete(update, msg)
}

func editPlaceHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.GetPlaceToEdit:
		// Expect user to select from inline keyboard markup. (name of place to edit)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		name, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}

		/* Get data from target chat */
		chatID, err := utils.GetChatTarget(update)
		if err != nil {
			log.Printf("error GetChatTarget: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		chatIDString := strconv.FormatInt(chatID, 10)
		if err := utils.CopyPlaceToTempPlace(update, name, chatIDString); err != nil {
			log.Printf("error CopyPlaceToTempPlace: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		// Use addplace logic to update
		sendTemplateReplies(update, "You may start editing")
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	}
}
