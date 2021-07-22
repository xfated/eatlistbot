package services

import (
	"context"
	"log"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"

	"google.golang.org/api/option"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	app    *firebase.App
	client *db.Client
)

func InitFirebase() {
	// initialize firebase app
	var err error
	ctx := context.Background()
	// Initialize the app with a custom auth variable, limiting the server's access
	ao := map[string]interface{}{"uid": "my-service-worker"}
	conf := &firebase.Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		AuthOverride: &ao,
	}
	// Fetch service account
	opt := option.WithCredentialsJSON([]byte(os.Getenv("SERVICE_ACCOUNT_JSON")))
	// Initialize app w service account
	app, err = firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("restricted_access/secret_document")
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	log.Println(data)

	// SaveData()
	log.Println("Loaded firebase")
}

func SaveData() {
	ctx := context.Background()
	ref := client.NewRef("levelOne")
	levelTwoRef := ref.Child("levelTwo")
	err := levelTwoRef.Set(ctx, "my first test")
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
}

func SetUserState(update tgbotapi.Update, state State) {
	ctx := context.Background()
	chatRef := client.NewRef(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err := chatRef.Child(strconv.Itoa(update.Message.From.ID)).Set(ctx,
		strconv.Itoa(state),
	); err != nil {
		log.Fatalln("Error setting state")
	}
}
