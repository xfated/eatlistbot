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
	msg := utils.CreateAndSendInlineKeyboard(update, text, 2, "yes", "no")
	utils.AddMessageToDelete(update, msg)
}

func checkAnyItem(update *tgbotapi.Update) error {
	/* Check if there are any items registed */
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserIDString: %+v", err)
	}

	itemNames, err := utils.GetItemNames(update, chatID)
	if err != nil {
		log.Printf("error GetItemNames: %+v", err)
		utils.SendMessage(update, "Sorry an error occured", false)
		return err
	}
	if len(itemNames) == 0 {
		utils.SendMessage(update, "No items registered :( add some", false)
		return errors.New("no item registered")
	}
	return nil
}

/* Search from available tags to get */
func addAndSendSelectedTags(update *tgbotapi.Update, tag string) {
	/* Extract tags */
	queryTagsMap, err := utils.GetQueryTags(update)
	if err != nil {
		log.Printf("error getting query tags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
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
	msg := utils.SendMessage(update, fmt.Sprintf("Selected tags: %s", curTags), false)
	utils.AddMessageToDelete(update, msg)
}

func sendAvailableTagsResponse(update *tgbotapi.Update, text string) {
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserID: %+v", chatID)
		utils.SendMessage(update, "Sorry, an error occured!", false)
		return
	}

	tagsMap, err := utils.GetTags(update, chatID)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
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

/* Search from name of items */
func sendAvailableItemNamesResponse(update *tgbotapi.Update, text string) {
	chatID, _, err := utils.GetChatUserIDString(update)
	if err != nil {
		log.Printf("error GetChatUserIDString: %+v", err)
	}

	itemNamesMap, err := utils.GetItemNames(update, chatID)
	if err != nil {
		log.Printf("error GetItemNames: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
		return
	}

	/* Set each name as its own inline row */
	itemNames := make([]string, len(itemNamesMap))
	i := 0
	for itemName := range itemNamesMap {
		itemNames[i] = itemName
		i++
	}
	msg := utils.CreateAndSendInlineKeyboard(update, text, 1, itemNames...)
	utils.AddMessageToDelete(update, msg)
}

func queryHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.QuerySelectType:
		// Expect user to select from inline markup keyboard
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options", false)
			utils.AddMessageToDelete(update, msg)
			return
		}
		// Get message
		message, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry, an error occured!", false)
			// return
		}
		switch message {
		case "/getOne":
			// getOne markup (/withTag, /withName), GoTo QueryOneTagOrName
			sendQueryOneTagOrNameResponse(update, "How do you want to search?")
			utils.SetQueryNum(update, 1)
			if err := utils.SetUserState(update, constants.QueryOneTagOrName); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		case "/getFew":
			// getFew GoTo QueryFewSetNum. Message how many they want?
			chatID, _, err := utils.GetChatUserIDString(update)
			if err != nil {
				log.Printf("error GetChatUserIDString: %+v", err)
			}

			itemNames, err := utils.GetItemNames(update, chatID)
			if err != nil {
				log.Printf("error GetItemNames: %+v", err)
				utils.SendMessage(update, "Sorry an error occured", false)
				return
			}

			// Store to delete
			msg := utils.RemoveMarkupKeyboard(update, fmt.Sprintf("You have %v recorded", len(itemNames)), false)
			utils.AddMessageToDelete(update, msg)
			// Get message
			messageID, err := utils.GetMessageTarget(update)
			if err != nil {
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			// Store to delete
			utils.SendMessageForceReply(update, "How many items do you want?", messageID, false)
			// if err != nil {
			// 	log.Printf("error SetMessageForceReply: %+v", err)
			// 	utils.SendMessage(update, "Sorry, an error occured!", false)
			// }
			// utils.AddMessageToDelete(update, msg)
			// Set state
			if err := utils.SetUserState(update, constants.QueryFewSetNum); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		case "/getAll":
			// getAll GoTo QuerySetTags
			chatID, _, err := utils.GetChatUserIDString(update)
			if err != nil {
				log.Printf("error GetChatUserIDString: %+v", err)
			}

			itemNames, err := utils.GetItemNames(update, chatID)
			if err != nil {
				log.Printf("error GetItemNames: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			utils.SetQueryNum(update, len(itemNames))
			// Store to delete
			msg := utils.RemoveMarkupKeyboard(update, "All in I see.", false)
			utils.AddMessageToDelete(update, msg)

			sendAvailableTagsResponse(update, "Add the tags you'd like to search with! \n\nPress \"/done\" once finished")
			msg = utils.SendMessage(update, "(Don't add any to consider all items)", false)
			utils.AddMessageToDelete(update, msg)
			if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			// sendQueryGetImagesResponse(update, "Do you want the images as well?")
			// if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
			// 	log.Printf("error SetUserState: %+v", err)
			// 	utils.SendMessage(update, "Sorry, an error occured!", false)
			// 	return
			// }
		}
	/* Ask to get one using tag or name */
	case constants.QueryOneTagOrName:
		// Expect user to select from inline markup keyboard (use tag or name to search)
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options", false)
			utils.AddMessageToDelete(update, msg)
			return
		}
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry, an error occured!", false)
			// return
		}

		message, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		switch message {
		case "/withTag":
			// Send responses
			msg := utils.RemoveMarkupKeyboard(update, "Searching for tags", false)
			utils.AddMessageToDelete(update, msg)

			sendAvailableTagsResponse(update, "Add the tags you'd like to search with! \n\nPress \"/done\" once finished")

			msg = utils.SendMessage(update, "(Don't add any to consider all items)", false)
			utils.AddMessageToDelete(update, msg)

			if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		case "/withName":
			msg := utils.RemoveMarkupKeyboard(update, "Searching for items", false)
			utils.AddMessageToDelete(update, msg)
			sendAvailableItemNamesResponse(update, "Which item do you want?")
			if err := utils.SetUserState(update, constants.QueryOneSetName); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		}

	/* Ask for name to search with */
	case constants.QueryOneSetName:
		// Expect user to select from inline markup keyboard (select name of item)
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options", false)
			utils.AddMessageToDelete(update, msg)
			return
		}

		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry, an error occured!", false)
			// return
		}
		name, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		utils.SetQueryName(update, name)
		sendQueryGetImagesResponse(update, "Do you want the images too? (if there is)")
		if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}

	/* Ask how many records to get */
	case constants.QueryFewSetNum:
		// Expect user to send a number (number of records to query)
		// Get queryNum
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error GetMessage: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		numQuery, err := strconv.Atoi(message)
		if err != nil || numQuery < 0 {
			msg := utils.SendMessage(update, "comeon, send a proper number", false)
			utils.AddMessageToDelete(update, msg)
			return
		}
		// By here, proper number received
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry, an error occured!", false)
			// return
		}

		// Add queryNum
		chatID, _, err := utils.GetChatUserIDString(update)
		if err != nil {
			log.Printf("error GetChatUserIDString: %+v", err)
		}

		itemNames, err := utils.GetItemNames(update, chatID)
		if err != nil {
			log.Printf("error GetItemNames: %+v", err)
			utils.SendMessage(update, "Sorry an error occured", false)
			return
		}

		if numQuery > len(itemNames) {
			msg := utils.SendMessage(update, fmt.Sprintf("thats too many. I'll just assume you want %v", len(itemNames)), false)
			utils.AddMessageToDelete(update, msg)
			utils.SetQueryNum(update, len(itemNames))
		} else {
			utils.SetQueryNum(update, numQuery)
		}

		sendAvailableTagsResponse(update, "Add the tags you'd like to search with! \n\nPress \"/done\" once finished")
		msg := utils.SendMessage(update, "(Don't add any to consider all items)", false)
		utils.AddMessageToDelete(update, msg)
		if err := utils.SetUserState(update, constants.QuerySetTags); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}

	/* Ask for tags to search with */
	case constants.QuerySetTags:
		// Expect user to select from inline keyboard markup (tags to include)
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options", false)
			utils.AddMessageToDelete(update, msg)
			return
		}

		// preview current, inline (show tags not yet added, /done)
		tag, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		// done GoTo QueryRetrieve. Markup("yes, no"), ask with pic
		switch tag {
		case "/done":
			// Delete messages
			if err := utils.DeleteRecentMessages(update); err != nil {
				log.Printf("error DeleteRecentMessages: %+v", err)
				// utils.SendMessage(update, "Sorry, an error occured!", false)
				// return
			}
			sendQueryGetImagesResponse(update, "Do you want the images too? (if there is)")
			if err := utils.SetUserState(update, constants.QueryRetrieve); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		default:
			addAndSendSelectedTags(update, tag)
		}
	/* Ask whether want pics, and retrieve */
	case constants.QueryRetrieve:
		/* If user send a message instead */
		if update.Message != nil {
			msg := utils.SendMessage(update, "Please select from the above options", false)
			utils.AddMessageToDelete(update, msg)
			return
		}
		// Delete messages
		if err := utils.DeleteRecentMessages(update); err != nil {
			log.Printf("error DeleteRecentMessages: %+v", err)
			// utils.SendMessage(update, "Sorry, an error occured!", false)
			// return
		}

		// Expect user to select from inline keyboard markup (yes or no to image)
		sendImage, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		if len(sendImage) > 0 {
			queryName, _ := utils.GetQueryName(update)

			// if name != "", get and show item data. (one result)
			if len(queryName) > 0 {
				chatID, _, err := utils.GetChatUserIDString(update)
				if err != nil {
					log.Printf("error GetChatUserID: %+v", err)
					utils.SendMessage(update, "Sorry, an error occured!", false)
					return
				}
				itemData, err := utils.GetItem(update, queryName, chatID)
				if err != nil {
					log.Printf("error GetItem: %+v", err)
					if err := utils.SetUserState(update, constants.Idle); err != nil {
						log.Printf("error SetUserState: %+v", err)
						utils.SendMessage(update, "Sorry, an error occured!", false)
						return
					}
					utils.SendMessage(update, "Sorry, error with getting data on the item.", false)
					return
				}
				utils.SendItemDetails(update, itemData, sendImage == "yes")
				if err := utils.SetUserState(update, constants.Idle); err != nil {
					log.Printf("error SetUserState: %+v", err)
					utils.SendMessage(update, "Sorry, an error occured!", false)
					return
				}
				return
			}

			// Get number of queries to return
			queryNum, err := utils.GetQueryNum(update)
			if err != nil {
				log.Printf("error GetQueryNum: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			// Get tags for filter
			queryTags, err := utils.GetQueryTags(update)
			if err != nil {
				log.Printf("error GetQueryTags: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}

			tagList := make([]string, len(queryTags))
			i := 0
			for tag := range queryTags {
				tagList[i] = tag
				i++
			}
			utils.SendMessage(update, fmt.Sprintf("Searching with tag(s): %+s", strings.Join(tagList, ", ")), false)
			// Get matching items
			// if len(tags) == 0, get all, randomly choose QueryNum
			// if len(tags) > 0, get all, extract with matching tags. randomly select queryNum
			items, err := utils.GetItems(update, queryTags)
			if err != nil {
				log.Printf("error GetItems: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			// less than queryNum found
			if len(items) < queryNum {
				utils.SendMessage(update, fmt.Sprintf("Found %v result(s) with matching tags", len(items)), false)
				queryNum = len(items)
			}
			for _, itemData := range items[:queryNum] {
				utils.SendItemDetails(update, itemData, sendImage == "yes")
			}
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		}
	}
}
