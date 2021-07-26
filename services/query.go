package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendQuerySelectType(update *tgbotapi.Update, text string) {
	utils.CreateAndSendInlineKeyboard(update, text, 3, "/getOne", "/getFew", "/getAll")
}

func sendQueryOneTagOrNameResponse(update *tgbotapi.Update, text string) {
	utils.CreateAndSendInlineKeyboard(update, text, 2, "/withTag", "/withName")
}

func sendQueryGetImagesResponse(update *tgbotapi.Update, text string) {
	utils.CreateAndSendInlineKeyboard(update, text, 2, "/yes", "/no")
}

/* Search from available tags to get */
func addAndSendSelectedTags(update *tgbotapi.Update, tag string) {
	utils.AddQueryTag(update, tag)

	/* Extract tags */
	queryTagsMap, err := utils.GetQueryTags(update)
	if err != nil {
		log.Printf("error getting query tags: %+v", err)
		utils.SendMessage(update, "Sorry an error occured!")
		return
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

func sendAvailableTagsResponse(update *tgbotapi.Update, text string) {
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserID: %+v", chatID)
		utils.SendMessage(update, "Sorry, an error occured!")
		return
	}

	tagsMap, err := utils.GetTags(update, chatID)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
		return
	}

	/* No tags, just send done */
	if len(tagsMap) == 0 {
		utils.CreateAndSendInlineKeyboard(update, "No tags found. Just help me click that done button thanks", 1, "/done")
		return
	}

	tags := make([]string, len(tagsMap)+1)
	i := 0
	for tag := range tagsMap {
		tags[i] = tag
		i++
	}
	tags[len(tagsMap)] = "/done"
	utils.CreateAndSendInlineKeyboard(update, text, 1, tags...)
}

/* Search from name of places */
func sendAvailablePlaceNamesResponse(update *tgbotapi.Update, text string) {
	placeNamesMap, err := utils.GetPlaceNames(update)
	if err != nil {
		log.Printf("error GetPlaceNames: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
		return
	}

	/* Set each name as its own inline row */
	placeNames := make([]string, len(placeNamesMap)+1)
	i := 0
	for placeName := range placeNamesMap {
		placeNames[i] = placeName
		i++
	}
	placeNames[len(placeNamesMap)] = "/done"
	utils.CreateAndSendInlineKeyboard(update, text, 1, placeNames...)
}

func queryHandler(update *tgbotapi.Update, userState constants.State) {
	/* Check if there are any places registed */
	placeNames, err := utils.GetPlaceNames(update)
	if err != nil {
		log.Printf("error GetPlaceNames: %+v", err)
		utils.SendMessage(update, "Sorry an error occured")
		return
	}
	if len(placeNames) == 0 {
		utils.SendMessage(update, "No places registered :( go add some")
		return
	}

	switch userState {
	case constants.QuerySelectType:
		// Expect user to select from inline markup keyboard
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}
		message, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		switch message {
		case "/getOne":
			// getOne markup (/withTag, /withName), GoTo QueryOneTagOrName
			sendQueryOneTagOrNameResponse(update, "How do you want to search?")
			utils.SetQueryNum(update, 1)
			if err := utils.SetUserState(update, constants.QueryOneTagOrName); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		case "/getFew":
			// getFew GoTo QueryFewSetNum. Message how many they want?
			utils.RemoveMarkupKeyboard(update, fmt.Sprintf("You have %v recorded", len(placeNames)))
			messageID, err := utils.GetMessageTarget(update)
			if err != nil {
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			utils.SendMessageForceReply(update, "How many places do you want?", messageID)
			if err := utils.SetUserState(update, constants.QueryFewSetNum); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		case "/getAll":
			// getAll GoTo QueryAllRetrieve
			placeNames, err := utils.GetPlaceNames(update)
			if err != nil {
				log.Printf("error GetPlaceNames: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			utils.SetQueryNum(update, len(placeNames))
			// Go straight to retrieve
			utils.RemoveMarkupKeyboard(update, fmt.Sprintf("All in I see. Shall go fetch your %v places", len(placeNames)))
			sendQueryGetImagesResponse(update, "Do you want the images as well?")
			if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		}
	/* Ask to get one using tag or name */
	case constants.QueryOneTagOrName:
		// Expect user to select from inline markup keyboard (use tag or name to search)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}
		message, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		switch message {
		case "/withTag":
			utils.RemoveMarkupKeyboard(update, "Searching for tags")
			sendAvailableTagsResponse(update, "Add the tags you'd like to search with! Press \"done\" once finished")
			utils.SendMessage(update, "(Don't add any to consider all places)")
			if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		case "/withName":
			utils.RemoveMarkupKeyboard(update, "Searching for places")
			sendAvailablePlaceNamesResponse(update, "Which place do you want?")
			if err := utils.SetUserState(update, constants.QueryOneSetName); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		}

	/* Ask for name to search with */
	case constants.QueryOneSetName:
		// Expect user to select from inline markup keyboard (select name of place)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		name, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		utils.SetQueryName(update, name)
		sendQueryGetImagesResponse(update, "Do you want the images too? (if there is)")
		if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}

	/* Ask how many records to get */
	case constants.QueryFewSetNum:
		// Expect user to send a number (number of records to query)
		// Get queryNum
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error GetMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		numQuery, err := strconv.Atoi(message)
		if err != nil || numQuery < 0 {
			utils.SendMessage(update, "comeon, send a proper number")
			return
		}

		// Add queryNum
		if numQuery > len(placeNames) {
			utils.SendMessage(update, fmt.Sprintf("thats too many. I'll just assume you want %v", len(placeNames)))
			utils.SetQueryNum(update, len(placeNames))
		} else {
			utils.SetQueryNum(update, numQuery)
		}

		sendAvailableTagsResponse(update, "Add the tags you'd like to search with! Press \"done\" once finished")
		utils.SendMessage(update, "(Don't add any to consider all places)")
		if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}

	/* Ask for tags to search with */
	case constants.QuerySetTags:
		// Expect user to select from inline keyboard markup (tags to include)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		// preview current, inline (show tags not yet added, /done)
		tag, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		// done GoTo QueryRetrieve. Markup("yes, no"), ask with pic
		if tag == "/done" {
			sendQueryGetImagesResponse(update, "Do you want the images too? (if there is)")
			if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		} else {
			addAndSendSelectedTags(update, tag)
		}

	/* Ask whether want pics, and retrieve */
	case constants.QueryRetrieve:
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		// Expect user to select from inline keyboard markup (yes or no to image)
		sendImage, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		if len(sendImage) > 0 {
			queryName, _ := utils.GetQueryName(update)

			// if name != "", get and show place data. (one result)
			if len(queryName) > 0 {
				placeData, err := utils.GetPlace(update, queryName)
				if err != nil {
					log.Printf("error GetPlace: %+v", err)
					if err := utils.SetUserState(update, constants.Idle); err != nil {
						log.Printf("error SetUserState: %+v", err)
						utils.SendMessage(update, "Sorry an error occured!")
						return
					}
					utils.SendMessage(update, "Sorry, error with getting data on the place.")
					return
				}
				utils.SendPlaceDetails(update, placeData, sendImage == "yes")
				if err := utils.SetUserState(update, constants.Idle); err != nil {
					log.Printf("error SetUserState: %+v", err)
					utils.SendMessage(update, "Sorry an error occured!")
					return
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
				return
			}
			// less than queryNum found
			if len(places) < queryNum {
				utils.SendMessage(update, fmt.Sprintf("Only %v results found with matching tags", len(places)))
				queryNum = len(places)
			}
			for _, placeData := range places[:queryNum] {
				utils.SendPlaceDetails(update, placeData, sendImage == "yes")
			}
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		}
	}
}
