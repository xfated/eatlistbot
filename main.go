// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xfated/eatlistbot/services"
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
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
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
