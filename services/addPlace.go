package services

import (
	"fmt"
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

/* Create and send template reply keyboard */
func sendTemplateReplies(update tgbotapi.Update, text string) {
	// Create buttons
	addAddressButton := tgbotapi.NewKeyboardButton("/addAddress")
	addNotesButton := tgbotapi.NewKeyboardButton("/addNotes")
	addURLButton := tgbotapi.NewKeyboardButton("/addURL")
	addImageButton := tgbotapi.NewKeyboardButton("/addImage")
	addTagButton := tgbotapi.NewKeyboardButton("/addTag")
	previewButton := tgbotapi.NewKeyboardButton("/preview")
	submitButton := tgbotapi.NewKeyboardButton("/submit")
	cancelButton := tgbotapi.NewKeyboardButton("/cancel")
	// Create rows
	row1 := tgbotapi.NewKeyboardButtonRow(addAddressButton, addURLButton)
	row2 := tgbotapi.NewKeyboardButtonRow(addNotesButton, addImageButton, addTagButton)
	row3 := tgbotapi.NewKeyboardButtonRow(cancelButton, previewButton, submitButton)

	replyKeyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.Selective = true
	utils.SetReplyMarkupKeyboard(update, text, replyKeyboard)
}
func sendExistingTagsResponse(update tgbotapi.Update, text string) {
	tagsMap, err := utils.GetTags(update)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!")
	}

	doneButton := tgbotapi.NewInlineKeyboardButtonData("/done", "/done")
	doneRow := tgbotapi.NewInlineKeyboardRow(doneButton)

	/* No tags, just send done */
	if len(tagsMap) == 0 {
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(doneRow)
		utils.SendInlineKeyboard(update, "No tags found. Just help me click that done button thanks", inlineKeyboard)
		return
	}

	/* Set each tag as its own inline row */
	var tagButtons = make([][]tgbotapi.InlineKeyboardButton, len(tagsMap)+1)
	i := 0
	for tag := range tagsMap {
		tagButtons[i] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(tag, tag),
		)
		i++
	}
	tagButtons[len(tagsMap)] = doneRow
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tagButtons...)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func sendConfirmSubmitResponse(update tgbotapi.Update, text string) {
	// Create buttons
	yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	// Create rows
	row := tgbotapi.NewInlineKeyboardRow(yesButton, noButton)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func addPlaceHandler(update tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.SetName:
		// Expect user to send a text message (name of place)
		// Message should contain name of place
		if err := utils.InitPlace(update); err != nil {
			log.Printf("Error creating new place: %+v", err)
			utils.SendMessage(update, "Message should be a text")
		}
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
		utils.SendMessage(update, "Start adding the details for the place")
	case constants.ReadyForNextAction:
		// Expect user to pick next state from reply markup
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "/addAddress":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetAddress); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send an address to be added")
		case "/addNotes":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetNotes); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Give some additional details as notes")
		case "/addURL":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetURL); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send a URL to be added")
		case "/addImage":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetImages); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send an image to be added")
		case "/addTag":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
			}
			utils.SetMessageTarget(update, messageID)
			utils.RemoveMarkupKeyboard(update, "Send a tag to be added. (Can be used to query your record of places)\n"+
				"Type new or pick from existing\nPress done once done!")
			sendExistingTagsResponse(update, "Existing tags:")
		case "/preview":
			// Get data and send
			placeData, err := utils.GetTempPlace(update)
			if err != nil {
				log.Printf("error getting temp place: %+v", err)
			}
			utils.SendPlaceDetails(update, placeData, true)
			sendTemplateReplies(update, "Select your next action")
		case "/submit":
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
			}
			if err := utils.SetUserState(update, constants.ConfirmAddPlaceSubmit); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.SetMessageTarget(update, messageID)
			sendConfirmSubmitResponse(update, "Are you really ready to submit?")
		case "/cancel":
			// Prep for next state
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "/addplace process cancelled")
		default:
			sendTemplateReplies(update, "Please select a response from the provided options")
		}
		return
	case constants.SetAddress:
		// Expect user to send a text message (address of place)
		// Message should contain address
		if err := utils.SetTempPlaceAddress(update); err != nil {
			log.Printf("Error adding address: %+v", err)
			utils.SendMessage(update, "Address should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("Address set to: %s", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetNotes:
		// Expect user to send a text message (notes for the place)
		// Message should contain notes
		if err := utils.SetTempPlaceNotes(update); err != nil {
			log.Printf("Error adding notes: %+v", err)
			utils.SendMessage(update, "Notes should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("Notes set to: %s", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetURL:
		// Expect user to send a text message (URL for the place)
		// Message should contain url
		if err := utils.SetTempPlaceURL(update); err != nil {
			log.Printf("Error adding url: %+v", err)
			utils.SendMessage(update, "URL should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("URL set to: %s", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetImages:
		// Expect user to send a photo
		// should be an image input
		if err := utils.AddTempPlaceImage(update); err != nil {
			log.Printf("Error adding image: %+v", err)
			utils.SendMessage(update, "Error occured. Did you send an image?")
		} else {
			utils.SendMessage(update, "Image added")
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetTags:
		// Expect user to send a text message or Select from inline keyboard markup (set as tag for the place)
		// First Check if is typed.
		if update.Message != nil {
			tag, _, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
			}
			if err := utils.AddTempPlaceTag(update, tag); err != nil {
				log.Printf("Error adding tag: %+v", err)
				utils.SendMessage(update, "Tag should be a text")
			} else {
				utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", tag))
			}
			// Only continue if /done is pressed
			return
		}

		// Then check if its a keyboard reply
		tag, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error GetCallbackQueryMessage: %+v", err)
		}
		if len(tag) > 0 {
			switch tag {
			case "/done":
				// Prep for next state
				if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
					log.Printf("error SetUserState: %+v", err)
					utils.SendMessage(update, "Sorry an error occured!")
				}
			default:
				if err := utils.AddTempPlaceTag(update, tag); err != nil {
					log.Printf("Error adding tag: %+v", err)
				} else {
					utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", tag))
				}
				return
				// Don't continue to next action if adding tag through inline
			}
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
		}
		if confirm == "yes" {
			// Submit
			name, err := utils.AddPlaceFromTemp(update)
			if err != nil {
				log.Printf("error adding place from temp: %+v", err)
				utils.SendMessage(update, "An error occured :( please try again")
			}
			// Prep for next state
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, fmt.Sprintf("%s was added for this chat!", name))
			return
		} else if confirm == "no" {
			if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
		}
	}

	/* Create and send keyboard for tarutils.Geted response */
	sendTemplateReplies(update, "What do you want to do next?")
}
