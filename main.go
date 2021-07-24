// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xfated/golistbot/services"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.Message == nil {
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	switch update.Message.Text {
	case "/start":
		services.SendStartInstructions(update)
	case "/addName":
		services.InitPlace(update)
	case "/addAddress":
		services.SetTempPlaceAddress(update)
	case "/addURL":
		services.SetTempPlaceURL(update)
	case "/addTags":
		services.AddTempPlaceTag(update)
	case "/addPlace":
		services.AddPlaceFromTemp(update)
	case "/deletePlace":
		services.DeletePlace(update, "addName")
	case "/updatePlaceAddress":
		services.UpdatePlaceAddress(update, "addName", "new address")
	case "/addMoreTags":
		services.AddPlaceTag(update, "addName", "more tags")
	case "/deleteTag":
		services.DeletePlaceTag(update, "addName", "more tags")
	default:
		services.LogUpdate(update)
		services.LogMessage(update)

		photoIDs, err := services.GetPhotoIDs(update)
		if err != nil {
			log.Printf("error: %+v", photoIDs)
			return
		}
		log.Printf("photoIDs: %+v", photoIDs)
		for _, id := range photoIDs {
			services.SendPhoto(update, id)
		}

		if err := services.SendMessage(update, "received"); err != nil {
			log.Printf("error sending message: %+v", err)
		}
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Print("$PORT must be set")
	}

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	// telegram
	services.InitTelegram()
	router.POST("/"+services.TELEGRAM_BOT_TOKEN, webhookHandler)

	// firebase
	services.InitFirebase()

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
