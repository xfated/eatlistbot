package services

import (
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleUserInput(update tgbotapi.Update) {
	userState, err := getUserState(update)
	if err != nil {
		log.Printf("error getting user state: %+v", err)
		return
	}
	if userState == Idle {
		switch update.Message.Text {
		case "/start":
			sendStartInstructions(update)
		}
	}

	switch update.Message.Text {
	case "/start":
		sendStartInstructions(update)
	case "/addName":
		InitPlace(update)
	case "/addAddress":
		setTempPlaceAddress(update)
	case "/addURL":
		setTempPlaceURL(update)
	case "/addTags":
		addTempPlaceTag(update)
	case "/addPlace":
		addPlaceFromTemp(update)
	case "/deletePlace":
		deletePlace(update, "addName")
	case "/updatePlaceAddress":
		updatePlaceAddress(update, "addName", "new address")
	case "/addMoreTags":
		addPlaceTag(update, "addName", "more tags")
	case "/deleteTag":
		deletePlaceTag(update, "addName", "more tags")
	default:
		LogUpdate(update)
		LogMessage(update)

		// photoIDs, err := services.GetPhotoIDs(update)
		// if err != nil {
		// 	log.Printf("error sending photo: %+v", photoIDs)
		// }
		// log.Printf("photoIDs: %+v", photoIDs)
		// for _, id := range photoIDs {
		// 	services.SendPhoto(update, id)
		// }

		if err := sendMessage(update, "received"); err != nil {
			log.Printf("error sending message: %+v", err)
		}
	}
}
