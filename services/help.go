package services

import (
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func helpHandler(update *tgbotapi.Update) {
	helpText := "/start or /reset: To reset the bot's status. (in case there are errors somehow) \n" +
		"/addplace: To add a new item to this chat's list (where this command was sent). Can be any item basically. You will be redirected to the bot's chat to add the item.\n" +
		"/deleteplace: To delete a place duh" +
		"/query: To fetch an item from this chat's list.\n" +
		"  /getOne: Returns one at random" +
		"    /withTag: Select multiple tags (or none). Filters for items with at least one matching tag" +
		"    /withName: Returns your selection" +
		"  /getFew: Returns a few (your choice) at random" +
		"    /withTag: Same as above" +
		"  /getAll: Returns all"

	utils.SendMessage(update, helpText)
}
