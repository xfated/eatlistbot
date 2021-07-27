package services

import (
	"fmt"
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendItemsToDeleteResponse(update *tgbotapi.Update, text string) {
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserIDString: %+v", err)
	}

	itemNames, err := utils.GetItemNames(update, chatID)
	if err != nil {
		log.Printf("error GetItemNames: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
	}

	/* Set each name as its own inline row */
	var nameButtons = make([][]tgbotapi.InlineKeyboardButton, len(itemNames))
	i := 0
	for name := range itemNames {
		nameButtons[i] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(name, name),
		)
		i++
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(nameButtons...)
	msg := utils.SendInlineKeyboard(update, text, inlineKeyboard)
	utils.AddMessageToDelete(update, msg)
}

func sendConfirmDeleteResponse(update *tgbotapi.Update, text string) {
	// Create buttons
	yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	// Create rows
	row := tgbotapi.NewInlineKeyboardRow(yesButton, noButton)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func deleteItemHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.DeleteSelect:
		// Expect user to select from inline keyboard markup. (name of items to delete)
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
		utils.SetItemTarget(update, name)
		sendConfirmDeleteResponse(update, "Are you sure?")
		if err := utils.SetUserState(update, constants.DeleteConfirm); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	case constants.DeleteConfirm:
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		confirm, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		if confirm == "yes" {
			target, err := utils.GetItemTarget(update)
			if err != nil {
				log.Printf("error GetItemTarget: %+v", err)
				utils.SendMessage(update, "Sorry an error occured")
				return
			}
			utils.DeleteItem(update, target)
			utils.SendMessage(update, fmt.Sprintf("%s has been deleted", target))
		} else if confirm == "no" {
			utils.SendMessage(update, "Deletion process cancelled")
		}
		if err := utils.SetUserState(update, constants.Idle); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	}
}
