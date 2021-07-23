package services

import (
	"context"
	"fmt"
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
		log.Println("Error initializing app:", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Println("Error initializing database client:", err)
	}

	log.Println("Loaded firebase")
}

/* User State */
func SetUserState(update tgbotapi.Update, state State) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf(fmt.Sprintf("Error getting chat and user data: %+v", err))
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child(chatID).Set(ctx,
		strconv.Itoa(int(state)),
	); err != nil {
		log.Println("Error setting state")
		return err
	}
	return nil
}

func GetUserState(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	var state State
	if err := userRef.Child(chatID).Get(ctx, &state); err != nil {
		log.Printf("Error getting user state: %+v", err)
		return err
	}
	return err
}

/* Name (Also init restaurant) */
func InitRestaurant(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}

	/* Set temp under userRef */
	name, _, err := GetMessage(update)
	if err != nil {
		log.Printf("Error getting message: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("restToAdd").Set(ctx, map[string]string{
		"name": name,
	}); err != nil {
		log.Printf("Error setting name: %+v", err)
		return err
	}
	return nil
}

/* Address */
func SetRestaurantAddress(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}

	/* Set temp under userRef */
	address, _, err := GetMessage(update)
	if err != nil {
		log.Printf("Error getting message: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("restToAdd").Update(ctx, map[string]interface{}{
		"address": address,
	}); err != nil {
		log.Printf("Error saving address: %+v", err)
		return err
	}

	return nil
}

/* URL */
func SetRestaurantURL(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}

	/* Set temp under userRef */
	url, _, err := GetMessage(update)
	if err != nil {
		log.Printf("Error getting message: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("restToAdd").Update(ctx, map[string]interface{}{
		"url": url,
	}); err != nil {
		log.Printf("Error saving url: %+v", err)
		return err
	}

	return nil
}

/* Images */
func AddRestaurantImage(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}

	/* Set temp under userRef */
	// TODO: FIND OUT HOW TO GET IMAGE URL
	imageUrl, _, err := GetMessage(update)
	if err != nil {
		log.Printf("Error getting message: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if _, err := userRef.Child("restToAdd").Child("images").Push(ctx, imageUrl); err != nil {
		log.Printf("Error saving image: %+v", err)
		return err
	}

	return nil
}

/* Tags */
func AddRestaurantTags(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}

	/* Set temp under userRef */
	tag, _, err := GetMessage(update)
	if err != nil {
		log.Printf("Error getting message: %+v", err)
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if _, err := userRef.Child("restToAdd").Child("tags").Push(ctx, tag); err != nil {
		log.Printf("Error saving image: %+v", err)
		return err
	}

	return nil
}

/* Get list of restaurants */

/* Add / Delete restaurant */
func getTempRestaurant(update tgbotapi.Update) (RestaurantDetails, error) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return RestaurantDetails{}, err
	}

	var RestaurantData RestaurantDetails
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("restToAdd").Get(ctx, &RestaurantData); err != nil {
		log.Printf("Error reading temp restaurant data: %+v", err)
		return RestaurantDetails{}, err
	}
	return RestaurantData, nil
}

func AddRestaurant(update tgbotapi.Update) error {
	restaurantData, err := getTempRestaurant(update)
	if err != nil {
		log.Printf("Error reading temp restaurant data: %+v", err)
		return err
	}

	ctx := context.Background()
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat and user data: %+v", err)
		return err
	}
	chatRef := client.NewRef("restaurants").Child(chatID)
	if _, err := chatRef.Push(ctx, restaurantData); err != nil {
		log.Printf("Error adding restaurant: %+v", err)
		return err
	}

	return nil
}

/* Read / Update tags */
