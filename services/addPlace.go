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

func addPlaceHandler(update tgbotapi.Update, userState constants.State) {

	switch userState {
	case constants.SetName:
		// Message should contain name of place
		if err := utils.InitPlace(update); err != nil {
			log.Printf("Error creating new place: %+v", err)
			utils.SendMessage(update, "Message should be a text")
		}
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
		utils.SendMessage(update, "Start adding the details for the place")
	case constants.ReadyForNextAction:
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "/addAddress":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetAddress); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send an address to be added")
		case "/addNotes":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetNotes); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Give some additional details as notes")
		case "/addURL":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetURL); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send a URL to be added")
		case "/addImage":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetImages); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send an image to be added")
		case "/addTag":
			// Prep for next state
			if err := utils.SetUserState(update, constants.SetTags); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "Send a tag to be added. (Can be used to query your record of places)")
		case "/preview":
			// Get data and send
			placeData, err := utils.GetTempPlace(update)
			if err != nil {
				log.Printf("error getting temp place: %+v", err)
			}
			utils.SendPlaceDetails(update, placeData)
			sendTemplateReplies(update, "Select your next action")
		case "/submit":
			// Submit
			name, err := utils.AddPlaceFromTemp(update)
			if err != nil {
				log.Printf("error adding place from temp: %+v", err)
				utils.SendMessage(update, "An error occured :( please try again")
			}
			// Prep for next state
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, fmt.Sprintf("%s was added for this chat!", name))
		case "/cancel":
			// Prep for next state
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
			}
			utils.RemoveMarkupKeyboard(update, "/addplace process cancelled")
		default:
			sendTemplateReplies(update, "Please select a response from the provided options")
		}
		return
	case constants.SetAddress:
		// Message should contain address
		if err := utils.SetTempPlaceAddress(update); err != nil {
			log.Printf("Error adding address: %+v", err)
			utils.SendMessage(update, "Address should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("Address set to: %s", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetNotes:
		// Message should contain notes
		if err := utils.SetTempPlaceNotes(update); err != nil {
			log.Printf("Error adding notes: %+v", err)
			utils.SendMessage(update, "Notes should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("Notes set to: %s", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetURL:
		// Message should contain url
		if err := utils.SetTempPlaceURL(update); err != nil {
			log.Printf("Error adding url: %+v", err)
			utils.SendMessage(update, "URL should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("URL set to: %s", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetImages:
		// should be an image input
		if err := utils.AddTempPlaceImage(update); err != nil {
			log.Printf("Error adding image: %+v", err)
			utils.SendMessage(update, "Error occured. Did you send an image?")
		} else {
			utils.SendMessage(update, "Image added")
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	case constants.SetTags:
		// Message should contain text
		if err := utils.AddTempPlaceTag(update); err != nil {
			log.Printf("Error adding tag: %+v", err)
			utils.SendMessage(update, "Tag should be a text")
		} else {
			utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", update.Message.Text))
		}
		// Prep for next state
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error setting state: %+v", err)
			utils.SendMessage(update, "Sorry an error occured!")
		}
	}

	/* Create and send keyboard for tarutils.Geted response */
	sendTemplateReplies(update, "What do you want to do next?")
}
