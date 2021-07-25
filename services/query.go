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
	yesButton := tgbotapi.NewInlineKeyboardButtonData("/yes", "/yes")
	noButton := tgbotapi.NewInlineKeyboardButtonData("/no", "/no")
	// Create rows
	row := tgbotapi.NewInlineKeyboardRow(yesButton, noButton)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

/* Search from available tags to get */
func addAndSendSelectedTags(update tgbotapi.Update, tag string) {
	utils.AddQueryTag(update, tag)

	/* Extract tags */
	queryTagsMap, err := utils.GetQueryTags(update)
	if err != nil {
		log.Printf("error getting query tags: %+v", err)
	}

	/* Send current tags */
	if len(queryTagsMap) > 0 {
		var queryTags = make([]string, len(queryTagsMap))
		i := 0
		for tag := range queryTagsMap {
			queryTags[i] = tag
			i++
		}
		curTags := strings.Join(queryTags, ", ")
		utils.SendMessage(update, fmt.Sprintf("Selected tags: %s", curTags))
	}
}

func sendAvailableTagsResponse(update tgbotapi.Update, text string) {
	tagsMap, err := utils.GetTags(update)
	if err != nil {
		log.Printf("error sending tags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
	}

	doneButton := tgbotapi.NewInlineKeyboardButtonData("/done", "/done")
	doneRow := tgbotapi.NewInlineKeyboardRow(doneButton)

	/* No tags, just send done */
	if len(tagsMap) == 0 {
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(doneRow)
		utils.SendInlineKeyboard(update, "No tags found", inlineKeyboard)
	}

	/* Set each tag as its own inline row */
	var tagButtons = make([][]tgbotapi.InlineKeyboardButton, len(tagsMap)+1)
	i := 0
	for tag := range tagsMap {
		tagButtons[i] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(tag, tag),
		)
		i++
		// /* Show if tag not chosen yet */
		// if !queryTagsMap[tag] {
		// 	tagButtons = append(tagButtons, tgbotapi.NewInlineKeyboardRow(
		// 		tgbotapi.NewInlineKeyboardButtonData(tag, tag),
		// 	))
		// }
	}
	tagButtons[len(tagsMap)] = doneRow
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tagButtons...)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

/* Search from name of places */
func sendAvailablePlaceNamesResponse(update tgbotapi.Update, text string) {

}

func queryHandler(update tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.QuerySelectType:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error getting message: %+v", err)
		}
		switch message {
		case "/getOne":
			// getOne markup (/withTag, /withName), GoTo QueryOneTagOrName
			sendQueryOneTagOrNameResponse(update, "How do you want to search?")
			utils.SetQueryNum(update, 1)
			if err := utils.SetUserState(update, constants.QueryOneTagOrName); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/getFew":
			// getFew GoTo QueryFewSetNum. Message how many they want?
			utils.RemoveMarkupKeyboard(update, "How many places do you want?")
			if err := utils.SetUserState(update, constants.QueryFewSetNum); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/getAll":
			// getAll GoTo QueryAllRetrieve
			placeNames, err := utils.GetPlaceNames(update)
			if err != nil {
				log.Printf("error getting place names: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.SetQueryNum(update, len(placeNames))
			// Go straight to retrieve
			sendQueryGetImagesResponse(update, "Do you want the images as well?")
			if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		default:
			sendQuerySelectType(update, "Please select a response from the provided options")
		}
	/* Ask to get one using tag or name */
	case constants.QueryOneTagOrName:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error getting message: %+v", err)
		}
		switch message {
		case "/withTag":
			// withTag inline (tags, /done), GoTo QuerySetTags
			// Send message "Don't add any to consider all places"
			utils.RemoveMarkupKeyboard(update, "Starting search with tags")
			sendAvailableTagsResponse(update, "Add the tags you'd like to search with! Press \"done\" once finished")
			utils.SendMessage(update, "(Don't add any to consider all places)")
			if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		case "/withName":
			// withName inline (names)
			utils.RemoveMarkupKeyboard(update, "Starting search with name")
			sendAvailablePlaceNamesResponse(update, "Which place do you want?")
			if err := utils.SetUserState(update, constants.QueryOneSetName); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		default:
			sendQueryOneTagOrNameResponse(update, "Please select one of the provided resposnes")
		}

	/* Ask how many records to get */
	case constants.QueryFewSetNum:

	/* Ask for tags to search with */
	case constants.QuerySetTags:
		// tag addTag, preview current, inline (show tags not yet added, /done)
		tag, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
		}
		// done GoTo QueryRetrieve. Markup("yes, no"), ask with pic
		if tag == "/done" {
			sendQueryGetImagesResponse(update, "Do you want the images too? (if there is)")
			if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		} else {
			addAndSendSelectedTags(update, tag)
		}

		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
		}

	/* Ask for name to search with */
	case constants.QueryOneSetName:
		// message, _, err := utils.GetMessage(update)
		// if err != nil {
		// 	log.Printf("error setting message: %+v", err)
		// }

		// set name, GoTo QueryOneRetrieve. Markup("yes, no"), ask with pics

	/* Ask whether want pics, and retrieve */
	case constants.QueryRetrieve:
		sendImage, err := utils.GetCallbackQueryMessage(update)
		log.Printf("sendImage: %+s", sendImage)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
		}
		if len(sendImage) > 0 {
			queryName, _ := utils.GetQueryName(update)

			// if name != "", get and show place data. (one result)
			if len(queryName) > 0 {
				// DEBUG
				// log.Printf("Sending query with name: %s", queryName)

				placeData, err := utils.GetPlace(update, queryName)
				if err != nil {
					log.Printf("error GetPlace: %+v", err)
					if err := utils.SetUserState(update, constants.Idle); err != nil {
						log.Printf("error SetUserState: %+v", err)
					}
					utils.SendMessage(update, "Sorry, error with getting data on the place.")
					return
				}
				utils.SendPlaceDetails(update, placeData, sendImage == "/yes")
				if err := utils.SetUserState(update, constants.Idle); err != nil {
					log.Printf("error SetUserState: %+v", err)
				}
				return
			}

			// Get number of queries to return
			queryNum, err := utils.GetQueryNum(update)
			if err != nil {
				log.Printf("error GetQueryNum: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			// Get tags for filter
			queryTags, err := utils.GetQueryTags(update)
			if err != nil {
				log.Printf("error GetQueryTags: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}

			// Get matching places
			// if len(tags) == 0, get all, randomly choose QueryNum
			// if len(tags) > 0, get all, extract with matching tags. randomly select queryNum
			places, err := utils.GetPlaces(update, queryTags)
			if err != nil {
				log.Printf("error GetPlaces: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			for _, placeData := range places[:queryNum] {
				log.Printf("senting details for: %+v", placeData)
				utils.SendPlaceDetails(update, placeData, sendImage == "/yes")
			}
		}

		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
		}
	}
}
