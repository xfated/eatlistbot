package services

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleUserInput(update tgbotapi.Update) {
	/* Check for main commands */
	message, _, err := getMessage(update)
	if err == nil {
		switch message {
		case "/start":
			sendStartInstructions(update)
			if err := setUserState(update, Idle); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
		}
	}

	/* Get user state for Targeted handling */
	userState, err := getUserState(update)
	if err != nil {
		log.Printf("error getting user state: %+v", err)
		return
	}

	/* Idle state */
	if userState == Idle {
		switch update.Message.Text {
		case "/addplace":
			if err := setUserState(update, SetName); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
			sendMessage(update, "Please enter the name of the place to begin")
		}
		return
	}

	/* Adding new place */
	if IsAddingNewPlace(userState) {
		switch userState {
		case SetName:
			// Message should contain name of place
			if err := initPlace(update); err != nil {
				log.Printf("Error creating new place: %+v", err)
				sendMessage(update, "Message should be a text")
			}
			if err := setUserState(update, ReadyForNextAction); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
			sendMessage(update, "Start adding the details for the place")
			return
		case ReadyForNextAction:
			message, _, err := getMessage(update)
			if err != nil {
				log.Printf("error getting message: %+v", err)
			}
			switch message {
			case "/addAddress":
				// Prep for next state
				if err := setUserState(update, SetAddress); err != nil {
					log.Printf("error setting state: %+v", err)
					sendMessage(update, "Sorry an error occured!")
				}
				sendMessage(update, "Send an address to be added")
			case "/addURL":
				// Prep for next state
				if err := setUserState(update, SetURL); err != nil {
					log.Printf("error setting state: %+v", err)
					sendMessage(update, "Sorry an error occured!")
				}
				sendMessage(update, "Send a URL to be added")
			case "/addImage":
				// Prep for next state
				if err := setUserState(update, SetImages); err != nil {
					log.Printf("error setting state: %+v", err)
					sendMessage(update, "Sorry an error occured!")
				}
				sendMessage(update, "Send an image to be added")
			case "/addTag":
				// Prep for next state
				if err := setUserState(update, SetTags); err != nil {
					log.Printf("error setting state: %+v", err)
					sendMessage(update, "Sorry an error occured!")
				}
				sendMessage(update, "Send a tag to be added")
			case "/preview":
				// Get data and send
				placeData, err := getTempPlace(update)
				if err != nil {
					log.Printf("error getting temp place: %+v", err)
				}
				placeText := ""
				if placeData.Name != "" {
					placeText = placeText + fmt.Sprintf("Name: %s\n", placeData.Name)
				}
				if placeData.Address != "" {
					placeText = placeText + fmt.Sprintf("Address: %s\n", placeData.Address)
				}
				if placeData.URL != "" {
					placeText = placeText + fmt.Sprintf("URL: %s\n", placeData.URL)
				}
				if placeData.Images != nil {
					placeText = placeText + fmt.Sprintf("Num images: %v\n", len(placeData.Images))
				}
				if placeData.Tags != nil {
					tags := make([]string, len(placeData.Tags))
					i := 0
					for tag := range placeData.Tags {
						tags[i] = tag
						i++
					}
					tagText := strings.Join(tags, ", ")
					placeText = placeText + fmt.Sprintf("Tags: %s\n", tagText)
				}
				sendMessage(update, placeText)
			case "/submit":
				// Submit
				name, err := addPlaceFromTemp(update)
				if err != nil {
					log.Printf("error adding place from temp: %+v", err)
					sendMessage(update, "An error occured :( please try again")
				}
				sendMessage(update, fmt.Sprintf("%s was added for this chat!", name))
				// Prep for next state
				if err := setUserState(update, Idle); err != nil {
					log.Printf("error setting state: %+v", err)
					sendMessage(update, "Sorry an error occured!")
				}
			case "/cancel":
				// Prep for next state
				if err := setUserState(update, Idle); err != nil {
					log.Printf("error setting state: %+v", err)
					sendMessage(update, "Sorry an error occured!")
				}
				sendMessage(update, "addPlace process cancelled")
			}
			return
		case SetAddress:
			// Message should contain address
			if err := setTempPlaceAddress(update); err != nil {
				log.Printf("Error adding address: %+v", err)
				sendMessage(update, "Message should be a text")
			}
			// Prep for next state
			if err := setUserState(update, ReadyForNextAction); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
		case SetURL:
			// Message should contain url
			if err := setTempPlaceURL(update); err != nil {
				log.Printf("Error adding url: %+v", err)
				sendMessage(update, "Message should be a text")
			}
			// Prep for next state
			if err := setUserState(update, ReadyForNextAction); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
		case SetImages:
			// should be an image input
			if err := addTempPlaceImage(update); err != nil {
				log.Printf("Error adding image: %+v", err)
				sendMessage(update, "error occured. did you send an image?")
			}
			// Prep for next state
			if err := setUserState(update, ReadyForNextAction); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
		case SetTags:
			// Message should contain text
			if err := addTempPlaceTag(update); err != nil {
				log.Printf("Error adding tag: %+v", err)
				sendMessage(update, "Message should be a text")
			}
			// Prep for next state
			if err := setUserState(update, ReadyForNextAction); err != nil {
				log.Printf("error setting state: %+v", err)
				sendMessage(update, "Sorry an error occured!")
			}
		}
		/* Create and send keyboard for targeted response */
		// Create buttons
		addAddressButton := tgbotapi.NewKeyboardButton("/addAddress")
		addURLButton := tgbotapi.NewKeyboardButton("/addURL")
		addImageButton := tgbotapi.NewKeyboardButton("/addImage")
		addTagButton := tgbotapi.NewKeyboardButton("/addTag")
		previewButton := tgbotapi.NewKeyboardButton("/preview")
		submitButton := tgbotapi.NewKeyboardButton("/submit")
		cancelButton := tgbotapi.NewKeyboardButton("/cancel")
		// Create rows
		row1 := tgbotapi.NewKeyboardButtonRow(addAddressButton, addURLButton)
		row2 := tgbotapi.NewKeyboardButtonRow(addImageButton, addTagButton)
		row3 := tgbotapi.NewKeyboardButtonRow(cancelButton, previewButton, submitButton)

		replyKeyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
		replyKeyboard.ResizeKeyboard = true
		replyKeyboard.OneTimeKeyboard = true
		replyKeyboard.Selective = true
		sendTextReplyMarkupKeyboard(update, "Select your next action", replyKeyboard)
		return
	}

	switch update.Message.Text {
	case "/start":
		sendStartInstructions(update)
		// case "/addName":
		// 	initPlace(update)
		// case "/addAddress":
		// 	setTempPlaceAddress(update)
		// case "/addURL":
		// 	setTempPlaceURL(update)
		// case "/addTags":
		// 	addTempPlaceTag(update)
		// case "/addPlace":
		// 	addPlaceFromTemp(update)
		// case "/deletePlace":
		// 	deletePlace(update, "addName")
		// case "/updatePlaceAddress":
		// 	updatePlaceAddress(update, "addName", "new address")
		// case "/addMoreTags":
		// 	addPlaceTag(update, "addName", "more tags")
		// case "/deleteTag":
		// 	deletePlaceTag(update, "addName", "more tags")
		// default:
		// 	LogUpdate(update)
		// 	LogMessage(update)

		// photoIDs, err := services.GetPhotoIDs(update)
		// if err != nil {
		// 	log.Printf("error sending photo: %+v", photoIDs)
		// }
		// log.Printf("photoIDs: %+v", photoIDs)
		// for _, id := range photoIDs {
		// 	services.SendPhoto(update, id)
		// }

		// if err := sendMessage(update, "received"); err != nil {
		// 	log.Printf("error sending message: %+v", err)
		// }
	}
}
