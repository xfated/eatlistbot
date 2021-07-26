// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xfated/golistbot/services"
	"github.com/xfated/golistbot/services/utils"
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

	if update.Message != nil {
		log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	}

	// Handle user input
	go services.HandleUserInput(&update)
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
	utils.InitTelegram()
	router.POST("/"+utils.TELEGRAM_BOT_TOKEN, webhookHandler)

	// firebase
	utils.InitFirebase()

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
