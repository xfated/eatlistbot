package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendItemsToEditResponse(update *tgbotapi.Update, text string) {
	// Get target chat
	chatID, err := utils.GetChatTarget(update)
	if err != nil {
		log.Printf("error GetChatTarget: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
		return
	}
	chatIDString := strconv.FormatInt(chatID, 10)

	// Extract item names
	itemNames, err := utils.GetItemNames(update, chatIDString)
	if err != nil {
		log.Printf("error GetItemNames: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
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
	msg := utils.SendInlineKeyboard(update, text, inlineKeyboard, false)
	utils.AddMessageToDelete(update, msg)
}

func editItemHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.GetItemToEdit:
		// Expect user to select from inline keyboard markup. (name of item to edit)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options", false)
			return
		}

		name, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}

		/* Get data from target chat */
		chatID, err := utils.GetChatTarget(update)
		if err != nil {
			log.Printf("error GetChatTarget: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		chatIDString := strconv.FormatInt(chatID, 10)
		if err := utils.CopyItemToTempItem(update, name, chatIDString); err != nil {
			log.Printf("error CopyItemToTempItem: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		// Use additem logic to update
		sendTemplateReplies(update, fmt.Sprintf(`You may start editing *%s*`, name))
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
	}
}
