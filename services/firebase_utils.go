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

	log.Println("Loaded firebase")
}

/* User State */
func SetUserState(update tgbotapi.Update, state State) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Fatalf("Error getting chat and user data: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child(chatID).Set(ctx,
		strconv.Itoa(int(state)),
	); err != nil {
		log.Fatalln("Error setting state")
	}
}

func GetUserState(update tgbotapi.Update) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Fatalf("Error getting chat and user data: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID)
	var state State
	if err := userRef.Child(chatID).Get(ctx, &state); err != nil {
		log.Fatalf("Error getting user state: %+v", err)
	}
}

/* Address */

/* URL */

/* Images */

/* Tags */
