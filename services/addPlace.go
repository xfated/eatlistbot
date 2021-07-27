package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

/* Create and send template reply keyboard */
func sendTemplateReplies(update *tgbotapi.Update, text string) {
	// Create buttons
	addAddressButton := tgbotapi.NewKeyboardButton("/addAddress")
	addNotesButton := tgbotapi.NewKeyboardButton("/addNotes")
	addURLButton := tgbotapi.NewKeyboardButton("/addURL")
	addImageButton := tgbotapi.NewKeyboardButton("/addImage")
	addTagButton := tgbotapi.NewKeyboardButton("/addTag")
	removeTagButton := tgbotapi.NewKeyboardButton("/removeTag")
	previewButton := tgbotapi.NewKeyboardButton("/preview")
	submitButton := tgbotapi.NewKeyboardButton("/submit")
	cancelButton := tgbotapi.NewKeyboardButton("/cancel")
	// Create rows
	row1 := tgbotapi.NewKeyboardButtonRow(addAddressButton, addURLButton, addNotesButton)
	row2 := tgbotapi.NewKeyboardButtonRow(addImageButton, addTagButton, removeTagButton)
	row3 := tgbotapi.NewKeyboardButtonRow(cancelButton, previewButton, submitButton)

	replyKeyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.Selective = false
	utils.SetReplyMarkupKeyboard(update, text, replyKeyboard)
}

func sendExistingTagsResponse(update *tgbotapi.Update, text string) {
	chatID, err := utils.GetChatTarget(update)
	if err != nil {
		log.Printf("Error GetChatTarget: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
		return
	}
	chatIDString := strconv.FormatInt(chatID, 10)

	tagsMap, err := utils.GetTags(update, chatIDString)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
		return
	}

	/* No tags, just send done */
	if len(tagsMap) == 0 {
		utils.CreateAndSendInlineKeyboard(update, "No tags found. Just click this button when you're done!", 1, "/done")
		return
	}

	/* Get already added tags */
	curTempTags, err := utils.GetTempPlaceTags(update)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
		return
	}
	tags := make([]string, 0)
	i := 0
	for tag := range tagsMap {
		// if not inside current temp tags
		if !curTempTags[tag] {
			tags = append(tags, tag)
			i++
		}
	}

	tags = append(tags, "/done")
	utils.CreateAndSendInlineKeyboard(update, text, 1, tags...)
}

func sendAddedTagsResponse(update *tgbotapi.Update, text string) {
	tagsMap, err := utils.GetTempPlaceTags(update)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
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

// func sendDoneResponse(update *tgbotapi.Update, text string) {
// 	utils.CreateAndSendInlineKeyboard(update, text, 1, "/done", "/done")
// }

func sendConfirmSubmitResponse(update *tgbotapi.Update, text string) {
	utils.CreateAndSendInlineKeyboard(update, text, 2, "yes", "no")
	// // Create buttons
	// yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	// noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	// // Create rows
	// row := tgbotapi.NewInlineKeyboardRow(yesButton, noButton)

	// inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	// utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func addPlaceHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.AddNewSetName:
		// Expect user to send a text message (name of place)
		// Check for slash (affect firebase query)
		if err := utils.CheckForSlash(update); err != nil {
			return
		}

		if err := utils.InitPlace(update); err != nil {
			log.Printf("Error creating new place: %+v", err)
			utils.SendMessage(update, "Message should be a text")
			break
		}
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			break
		}
		utils.SendMessage(update, "You may start adding the details for the place")
	case constants.ReadyForNextAction:
		// Expect user to select reply markup (pick next action)
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "/addAddress":
			if err := utils.SetUserState(update, constants.AddNewSetAddress); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.RemoveMarkupKeyboard(update, "Send an address to be added")
		case "/addNotes":
			if err := utils.SetUserState(update, constants.AddNewSetNotes); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.RemoveMarkupKeyboard(update, "Give some additional details as notes")
		case "/addURL":
			if err := utils.SetUserState(update, constants.AddNewSetURL); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.RemoveMarkupKeyboard(update, "Send a URL to be added")
		case "/addImage":
			if err := utils.SetUserState(update, constants.AddNewSetImages); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.RemoveMarkupKeyboard(update, "Send an image to be added")
		case "/addTag":
			if err := utils.SetUserState(update, constants.AddNewSetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			/* Get message ID for targeted reply afterward */
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.SetMessageTarget(update, messageID)

			utils.RemoveMarkupKeyboard(update, "Send a tag to be added. (Can be used to query your record of places)\n"+
				"Type new or pick from existing\n\nPress \"/done\" once done!")
			sendExistingTagsResponse(update, "Existing tags:")
		case "/removeTag":
			if err := utils.SetUserState(update, constants.AddNewRemoveTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			/* Get message ID for targeted reply afterward */
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.SetMessageTarget(update, messageID)

			utils.RemoveMarkupKeyboard(update, "Select a tag to remove\n\nPress \"/done\" once done!")
			sendAddedTagsResponse(update, "Existing tags:")
		case "/preview":
			placeData, err := utils.GetTempPlace(update)
			if err != nil {
				log.Printf("error getting temp place: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.SendPlaceDetails(update, placeData, true)
			sendTemplateReplies(update, "Select your next action")
		case "/submit":
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			if err := utils.SetUserState(update, constants.ConfirmAddPlaceSubmit); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.SetMessageTarget(update, messageID)
			sendConfirmSubmitResponse(update, "Are you really ready to submit?")
		case "/cancel":
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				break
			}
			utils.RemoveMarkupKeyboard(update, "/addplace process cancelled")
		default:
			sendTemplateReplies(update, "Please select a response from the provided options")
		}
		return
	case constants.AddNewSetAddress:
		// Expect user to send a text message (address of place)
		if err := utils.SetTempPlaceAddress(update); err != nil {
			log.Printf("Error adding address: %+v", err)
			utils.SendMessage(update, "Address should be a text")
			return
		}
		utils.SendMessage(update, fmt.Sprintf("Address set to: %s", update.Message.Text))

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	case constants.AddNewSetNotes:
		// Expect user to send a text message (notes for the place)
		if err := utils.SetTempPlaceNotes(update); err != nil {
			log.Printf("Error adding notes: %+v", err)
			utils.SendMessage(update, "Notes should be a text")
			return
		}
		utils.SendMessage(update, fmt.Sprintf("Notes set to: %s", update.Message.Text))

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	case constants.AddNewSetURL:
		// Expect user to send a text message (URL for the place)
		if err := utils.SetTempPlaceURL(update); err != nil {
			log.Printf("Error adding url: %+v", err)
			utils.SendMessage(update, "URL should be a text")
			return
		}
		utils.SendMessage(update, fmt.Sprintf("URL set to: %s", update.Message.Text))

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	case constants.AddNewSetImages:
		// Expect user to send a photo
		// should be an image input
		if err := utils.AddTempPlaceImage(update); err != nil {
			log.Printf("Error adding image: %+v", err)
			utils.SendMessage(update, "Error occured. Did you send an image? Try it again")
			return
		}
		utils.SendMessage(update, "Image added")

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}
	case constants.AddNewSetTags:
		// Expect user to send a text message or Select from inline keyboard markup (set as tag for the place)
		// Check for slash (affect firebase query)
		if err := utils.CheckForSlash(update); err != nil {
			return
		}

		if update.Message != nil {
			tag, _, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			switch tag {
			case "/done",
				"done",
				"Done":
				if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
					log.Printf("error SetUserState: %+v", err)
					utils.SendMessage(update, "Sorry an error occured!")
					return
				}
				// Only continue if /done is pressed
			default:
				if err := utils.AddTempPlaceTag(update, tag); err != nil {
					log.Printf("Error adding tag: %+v", err)
					utils.SendMessage(update, "Tag should be a text")
					return
				}
				utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", tag))
				return
			}
		} else {
			// Then check if its a keyboard reply
			tag, err := utils.GetCallbackQueryMessage(update)
			if err != nil {
				log.Printf("error GetCallbackQueryMessage: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			if len(tag) > 0 {
				switch tag {
				case "/done":
					if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
						log.Printf("error SetUserState: %+v", err)
						utils.SendMessage(update, "Sorry an error occured!")
						return
					}
				default:
					if err := utils.AddTempPlaceTag(update, tag); err != nil {
						log.Printf("Error adding tag: %+v", err)
						utils.SendMessage(update, "Sorry an error occured!")
						return
					}
					utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", tag))
					// Don't continue to next action if adding tag through inline
					return
				}
			}
		}
	case constants.AddNewRemoveTags:
		// Expect user to select from inline keyboard markup (set as tag for the place)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options")
			return
		}

		tag, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
			return
		}

		switch tag {
		case "/done":
			if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		default:
			// remove tag
			if err := utils.DeleteTempPlaceTag(update, tag); err != nil {
				log.Printf("Error adding tag: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" removed", tag))
			sendAddedTagsResponse(update, "Existing tags:")
			// Don't continue to next action if removing tag through inline
			return
		}
	case constants.ConfirmAddPlaceSubmit:
		// Expect user to select from inline query (yes or no to submit)
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
			// Get target chat, where addplace was initiated
			chatID, err := utils.GetChatTarget(update)
			chatIDString := strconv.FormatInt(chatID, 10)
			if err != nil {
				log.Printf("error getting message from callback: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}

			// Submit
			name, err := utils.AddPlaceFromTemp(update, chatIDString)
			if err != nil {
				log.Printf("error adding place from temp: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			utils.RemoveMarkupKeyboard(update, fmt.Sprintf("%s has been added!", name))
			utils.SendMessage(update, "To add a new place to the chat, please initiate /addplace back in that chat")
			err = utils.SetChatTarget(update, 0)
			if err != nil {
				log.Printf("error SetChatTarget: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			return
		} else if confirm == "no" {
			if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
		}
	}

	/* Create and send keyboard for targeted response */
	sendTemplateReplies(update, "What do you want to do next?")
}
