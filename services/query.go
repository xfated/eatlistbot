package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendQuerySelectType(update *tgbotapi.Update, text string) {
	msg := utils.CreateAndSendInlineKeyboard(update, text, 3, "/getOne", "/getFew", "/getAll")
	utils.AddMessageToDelete(update, msg)
}

func sendQueryOneTagOrNameResponse(update *tgbotapi.Update, text string) {
	msg := utils.CreateAndSendInlineKeyboard(update, text, 2, "/withTag", "/withName")
	utils.AddMessageToDelete(update, msg)
}

func sendQueryGetImagesResponse(update *tgbotapi.Update, text string) {
	msg := utils.CreateAndSendInlineKeyboard(update, text, 2, "/yes", "/no")
	utils.AddMessageToDelete(update, msg)
}

func checkAnyPlace(update *tgbotapi.Update) error {
	/* Check if there are any places registed */
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserIDString: %+v", err)
	}

	placeNames, err := utils.GetPlaceNames(update, chatID)
	if err != nil {
		log.Printf("error GetPlaceNames: %+v", err)
		utils.SendMessage(update, "Sorry an error occured")
		return err
	}
	if len(placeNames) == 0 {
		utils.SendMessage(update, "No places registered :( add some")
		return errors.New("no place registered")
	}
	return nil
}

/* Search from available tags to get */
func addAndSendSelectedTags(update *tgbotapi.Update, tag string) {
	/* Extract tags */
	queryTagsMap, err := utils.GetQueryTags(update)
	if err != nil {
		log.Printf("error getting query tags: %+v", err)
		utils.SendMessage(update, "Sorry an error occured!")
		return
	}

	/* if already selected */
	if queryTagsMap[tag] {
		return
	}
	go utils.AddQueryTag(update, tag)

	/* Send current tags */
	var queryTags = make([]string, len(queryTagsMap)+1)
	i := 0
	for tag := range queryTagsMap {
		queryTags[i] = tag
		i++
	}
	queryTags[len(queryTagsMap)] = tag
	curTags := strings.Join(queryTags, ", ")
	msg := utils.SendMessage(update, fmt.Sprintf("Selected tags: %s", curTags))
	utils.AddMessageToDelete(update, msg)
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
	msg := utils.CreateAndSendInlineKeyboard(update, text, 1, tags...)
	utils.AddMessageToDelete(update, msg)
}

/* Search from name of places */
func sendAvailablePlaceNamesResponse(update *tgbotapi.Update, text string) {
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserIDString: %+v", err)
	}

	placeNamesMap, err := utils.GetPlaceNames(update, chatID)
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
	msg := utils.CreateAndSendInlineKeyboard(update, text, 1, placeNames...)
	utils.AddMessageToDelete(update, msg)
}

func queryHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.QuerySelectType:
		// Expect user to select from inline markup keyboard
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options")
			utils.AddMessageToDelete(update, msg)
			return
		}
		// Get message
		message, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry an error occured!")
			// return
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
			chatID, _, err := utils.GetChatUserIDString(update)
			if err != nil {
				log.Printf("error GetChatUserIDString: %+v", err)
			}

			placeNames, err := utils.GetPlaceNames(update, chatID)
			if err != nil {
				log.Printf("error GetPlaceNames: %+v", err)
				utils.SendMessage(update, "Sorry an error occured")
				return
			}

			// Store to delete
			msg := utils.RemoveMarkupKeyboard(update, fmt.Sprintf("You have %v recorded", len(placeNames)))
			utils.AddMessageToDelete(update, msg)
			// Get message
			messageID, err := utils.GetMessageTarget(update)
			if err != nil {
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			// Store to delete
			utils.SendMessageForceReply(update, "How many placeFs do you want?", messageID)
			// if err != nil {
			// 	log.Printf("error SetMessageForceReply: %+v", err)
			// 	utils.SendMessage(update, "Sorry an error occured!")
			// }
			// utils.AddMessageToDelete(update, msg)
			// Set state
			if err := utils.SetUserState(update, constants.QueryFewSetNum); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		case "/getAll":
			// getAll GoTo QueryAllRetrieve
			chatID, _, err := utils.GetChatUserIDString(update)
			if err != nil {
				log.Printf("error GetChatUserIDString: %+v", err)
			}

			placeNames, err := utils.GetPlaceNames(update, chatID)
			if err != nil {
				log.Printf("error GetPlaceNames: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			utils.SetQueryNum(update, len(placeNames))
			// Store to delete
			msg := utils.RemoveMarkupKeyboard(update, fmt.Sprintf("All in I see. Shall fetch your %v places", len(placeNames)))
			utils.AddMessageToDelete(update, msg)
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
			msg := utils.SendMessage(update, "Please select from the above options")
			utils.AddMessageToDelete(update, msg)
			return
		}
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry an error occured!")
			// return
		}

		message, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		switch message {
		case "/withTag":
			// Send responses
			msg := utils.RemoveMarkupKeyboard(update, "Searching for tags")
			utils.AddMessageToDelete(update, msg)

			sendAvailableTagsResponse(update, "Add the tags you'd like to search with! Press \"done\" once finished")

			msg = utils.SendMessage(update, "(Don't add any to consider all places)")
			utils.AddMessageToDelete(update, msg)

			if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		case "/withName":
			msg := utils.RemoveMarkupKeyboard(update, "Searching for places")
			utils.AddMessageToDelete(update, msg)
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
			msg := utils.SendMessage(update, "Please select from the above options")
			utils.AddMessageToDelete(update, msg)
			return
		}

		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry an error occured!")
			// return
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
			msg := utils.SendMessage(update, "comeon, send a proper number")
			utils.AddMessageToDelete(update, msg)
			return
		}
		// By here, proper number received
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry an error occured!")
			// return
		}

		// Add queryNum
		chatID, _, err := utils.GetChatUserIDString(update)
		if err != nil {
			log.Printf("error GetChatUserIDString: %+v", err)
		}

		placeNames, err := utils.GetPlaceNames(update, chatID)
		if err != nil {
			log.Printf("error GetPlaceNames: %+v", err)
			utils.SendMessage(update, "Sorry an error occured")
			return
		}

		if numQuery > len(placeNames) {
			msg := utils.SendMessage(update, fmt.Sprintf("thats too many. I'll just assume you want %v", len(placeNames)))
			utils.AddMessageToDelete(update, msg)
			utils.SetQueryNum(update, len(placeNames))
		} else {
			utils.SetQueryNum(update, numQuery)
		}

		sendAvailableTagsResponse(update, "Add the tags you'd like to search with! Press \"done\" once finished")
		msg := utils.SendMessage(update, "(Don't add any to consider all places)")
		utils.AddMessageToDelete(update, msg)
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
			msg := utils.SendMessage(update, "Please select from the above options")
			utils.AddMessageToDelete(update, msg)
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
		switch tag {
		case "/done":
			// Delete messages
			if err := utils.DeleteRecentMessages(update); err != nil {
				log.Printf("error DeleteRecentMessages: %+v", err)
				// utils.SendMessage(update, "Sorry an error occured!")
				// return
			}
			sendQueryGetImagesResponse(update, "Do you want the images too? (if there is)")
			if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		default:
			addAndSendSelectedTags(update, tag)
		}
	/* Ask whether want pics, and retrieve */
	case constants.QueryRetrieve:
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options")
			utils.AddMessageToDelete(update, msg)
			return
		}
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry an error occured!")
			// return
		}

		// Expect user to select from inline keyboard markup (yes or no to image)
		sendImage, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
		if len(sendImage) > 0 {
			queryName, _ := utils.GetQueryName(update)

			// if name != "", get and show place data. (one result)
			if len(queryName) > 0 {
				chatID, _, err := utils.GetChatUserIDString(update)
				if err != nil {
					log.Printf("error GetChatUserID: %+v", err)
					utils.SendMessage(update, "Sorry an error occured!")
					return
				}
				placeData, err := utils.GetPlace(update, queryName, chatID)
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

			tagList := make([]string, len(queryTags))
			i := 0
			for tag := range queryTags {
				tagList[i] = tag
				i++
			}
			utils.SendMessage(update, fmt.Sprintf("Searching with tag(s): %+s", strings.Join(tagList, ", ")))
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
				utils.SendMessage(update, fmt.Sprintf("Only found %v result(s) with matching tags", len(places)))
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
