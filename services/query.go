package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendQuerySelectType(update tgbotapi.Update, text string) {
	// Create buttons
	getOneButton := tgbotapi.NewKeyboardButton("/getOne")
	getFewButton := tgbotapi.NewKeyboardButton("/getFew")
	getAllButton := tgbotapi.NewKeyboardButton("/getAll")
	// Create rows
	row := tgbotapi.NewKeyboardButtonRow(getOneButton, getFewButton, getAllButton)

	replyKeyboard := tgbotapi.NewReplyKeyboard(row)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.Selective = true
	utils.SetReplyMarkupKeyboard(update, text, replyKeyboard)
}

func sendQueryOneTagOrNameResponse(update tgbotapi.Update, text string) {
	// Create buttons
	withTagButton := tgbotapi.NewKeyboardButton("/withTag")
	withNameButton := tgbotapi.NewKeyboardButton("/withName")
	// Create rows
	row := tgbotapi.NewKeyboardButtonRow(withTagButton, withNameButton)

	replyKeyboard := tgbotapi.NewReplyKeyboard(row)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.Selective = true
	utils.SetReplyMarkupKeyboard(update, text, replyKeyboard)
}

func sendQueryGetImagesResponse(update tgbotapi.Update, text string) {
	// Create buttons
	yesButton := tgbotapi.NewKeyboardButton("/yes")
	noButton := tgbotapi.NewKeyboardButton("/no")
	// Create rows
	row := tgbotapi.NewKeyboardButtonRow(yesButton, noButton)

	replyKeyboard := tgbotapi.NewReplyKeyboard(row)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.Selective = true
	utils.SetReplyMarkupKeyboard(update, text, replyKeyboard)
}

/* Search from available tags to get */
func sendAvailableTagsResponse(update tgbotapi.Update, text string) {
	tags, err := utils.GetTags(update)
	if err != nil {
		log.Printf("error sending tags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
	}

	doneButton := tgbotapi.NewInlineKeyboardButtonData("/done", "/done")
	doneRow := tgbotapi.NewInlineKeyboardRow(doneButton)

	/* No tags, just send done */
	if len(tags) == 0 {
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(doneRow)
		utils.SendInlineKeyboard(update, "No tags found", inlineKeyboard)
	}

	/* Extract tags */
	queryTags, err := utils.GetQueryTags(update)
	if err != nil {
		log.Printf("error getting query tags: %+v", err)
	}

	/* Send current tags */
	if len(queryTags) > 0 {
		curTags := strings.Join(queryTags, ", ")
		utils.SendMessage(update, fmt.Sprintf("Current tags: %s", curTags))
	}
	/* Set each tag as its own inline row */
	var tagButtons = make([][]tgbotapi.InlineKeyboardButton, len(queryTags)+1)
	for i, tag := range queryTags {
		tagButtons[i] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(tag, tag),
		)
	}
	tagButtons[len(tagButtons)] = doneRow
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tagButtons...)
	utils.SendInlineKeyboard(update, "Add another tag or done", inlineKeyboard)
}

/* Search from name of places */
func sendAvailablePlaceNamesResponse(update tgbotapi.Update, text string) {

}

func queryHandler(update tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.QuerySelectType:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "/getOne":
			// getOne markup (/withTag, /withName), GoTo QueryOneTagOrName
			sendQueryOneTagOrNameResponse(update, "How do you want to search?")
			if err := utils.SetUserState(update, constants.QueryOneTagOrName); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/getFew":
			// getFew GoTo QueryFewSetNum. Message how many they want?
			if err := utils.SetUserState(update, constants.QueryFewSetNum); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/getAll":
			// getAll GoTo QueryAllRetrieve
			if err := utils.SetUserState(update, constants.QueryAllRetrieve); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		default:
			sendQuerySelectType(update, "Please select a response from the provided options")
		}
	/* Ask to get one using tag or name */
	case constants.QueryOneTagOrName:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "/withTag":
			// withTag inline (tags, /done), GoTo QuerySetTags
			// Send message "Don't add any to consider all places"
			sendAvailableTagsResponse(update, "Add the tags you'd like to search with!")
			utils.SendMessage(update, "(Don't add any to consider all places)")
			if err := utils.SetUserState(update, constants.QueryOneSetTags); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/withName":
			// withName inline (names)
			sendAvailablePlaceNamesResponse(update, "Which place do you want?")
			if err := utils.SetUserState(update, constants.QueryOneSetName); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		default:
			sendQueryOneTagOrNameResponse(update, "Please select one of the provided resposnes")
		}

	/* Ask for tags to search with */
	case constants.QueryOneSetTags:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "":
		default:
			sendAvailableTagsResponse(update, "Please select one of the provided responses")
		}
		// tag addTag, preview current, inline (show tags not yet added, /done)
		// done GoTo QueryOneRetrieve. Markup("yes, no"), ask with pic

	/* Ask for name to search with */
	case constants.QueryOneSetName:
		// message, _, err := utils.GetMessage(update)
		// if err != nil {
		// 	log.Printf("error setting message: %+v", err)
		// }

		// set name, GoTo QueryOneRetrieve. Markup("yes, no"), ask with pics

	/* Ask whether want pics, and retrieve */
	case constants.QueryOneRetrieve:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "":
		default:
			sendQueryGetImagesResponse(update, "yes or no?")
		}
		// if name != "", get and show place data.
		// if len(tags) == 0, get all, randomly choose one
		// if len(tags) > 0, get all, extract with matching tags. randomly select one
		// if "yes", send pics, goto Idle

	/* Ask how many they want */
	case constants.QueryFewSetNum:
		// message, _, err := utils.GetMessage(update)
		// if err != nil {
		// 	log.Printf("error setting message: %+v", err)
		// }

		// Set number, GoTo QueryFewSetTags. inline (tags, /done)
		// Send message "Don't add any to consider all places"

	/* Ask for tags to search with */
	case constants.QueryFewSetTags:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "":
		default:
			sendAvailableTagsResponse(update, "Please select one of the provided resposnes")
		}
		// tag addTag, preview current, inline (show tags not yet added, /done)
		// done GoTo QueryFewRetrieve, Markup("yes, no") ask with pics

	/* Ask whether want pics and retrieve */
	case constants.QueryFewRetrieve:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "":
		default:
			sendQueryGetImagesResponse(update, "yes or no?")
		}
		// if len(tags) == 0, randomly select
		// if "yes", send pics. GoTo Idle

	/* Ask whether want pics and retrieve */
	case constants.QueryAllRetrieve:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "":
		default:
			sendQueryGetImagesResponse(update, "yes or no?")
		}
		// fetch all, send one by one

		// if "yes", send pics. GoTo Idle
	}
}
