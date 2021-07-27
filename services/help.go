package services

import (
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func helpHandler(update *tgbotapi.Update) {
	helpText := "/start or /reset: To reset the bot's status. (in case there are errors somehow) \n" +
		"\n" +
		"/additem: To add a new item to this chat's list (where this command was sent). Can be any item basically. You will be redirected to the bot's chat to add the item. \n" +
		"    /setXX: Adds (or overwrites) the field \n" +
		"    /addXX: Tag or Image. You can add multiple \n" +
		"\n" +
		"/deleteitem: To delete a item. Forever. \n" +
		"\n" +
		"/edititem: To edit a item. Similar process to /additem \n" +
		"\n" +
		"/query: To fetch an item from this chat's list.\n" +
		"    /getOne: Returns one at random \n" +
		"        /withTag: Select multiple tags (or none). Filters for items with at least one matching tag \n" +
		"        /withName: Returns your selection \n" +
		"    /getFew: Returns a few (your choice) at random \n" +
		"        /withTag: Same as above \n" +
		"    /getAll: Returns all"

	utils.SendMessage(update, helpText, false)
}
