// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xfated/eatlistbot/firebase"
	"github.com/xfated/eatlistbot/telegram"
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

	// to monitor changes run: heroku logs --tail
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)

	// Reply message (just for lels)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	telegram.SendMessage(msg)
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
	telegram.InitTelegram()
	router.POST("/"+telegram.TELEGRAM_BOT_TOKEN, webhookHandler)

	// firebase
	firebase.InitFirebase()

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
