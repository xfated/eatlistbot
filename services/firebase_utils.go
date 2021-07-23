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
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	var state State
	if err := userRef.Child(chatID).Get(ctx, &state); err != nil {
		return err
	}
	return err
}

/* Name (Also init restaurant) */
func InitRestaurant(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	name, _, err := GetMessage(update)
	if err != nil {
		return err
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
		return err
	}

	/* Set temp under userRef */
	address, _, err := GetMessage(update)
	if err != nil {
		return err
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
		return err
	}

	/* Set temp under userRef */
	url, _, err := GetMessage(update)
	if err != nil {
		return err
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
		return err
	}

	/* Set temp under userRef */
	// TODO: FIND OUT HOW TO GET IMAGE URL
	imageUrl, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("restToAdd").Child("images").Update(ctx, map[string]interface{}{
		imageUrl: true,
	}); err != nil {
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
		return err
	}

	/* Set temp under userRef */
	tag, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("restToAdd").Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		log.Printf("Error saving image: %+v", err)
		return err
	}

	return nil
}

/* Get list of restaurants */

/* Add / Delete restaurant */
func GetTempRestaurant(update tgbotapi.Update) (RestaurantDetails, error) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserID(update)
	if err != nil {
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
	restaurantData, err := GetTempRestaurant(update)
	if err != nil {
		return err
	}

	ctx := context.Background()
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return err
	}

	/* Add restaurant to restaurant collection */
	chatRef := client.NewRef("restaurants").Child(chatID)
	if err := chatRef.Child(restaurantData.Name).Set(ctx, restaurantData); err != nil {
		log.Printf("Error adding restaurant: %+v", err)
		return err
	}

	/* Add tags to tag collection */
	for tag := range restaurantData.Tags {
		if err := updateTags(update, tag); err != nil {
			return err
		}
	}
	return nil
}

func DeleteRestaurant(update tgbotapi.Update, restaurantName string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return err
	}
	chatRef := client.NewRef("restaurants").Child(chatID)
	if err := chatRef.Child(restaurantName).Delete(ctx); err != nil {
		log.Printf("Error deletting restaurant: %+v", err)
		return err
	}
	return nil
}

/* Read / Delete / Update tags */
func GetTags(update tgbotapi.Update) ([]string, error) {
	ctx := context.Background()
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return []string{}, err
	}
	chatRef := client.NewRef("tags").Child(chatID)

	/* Retrieve tags */
	var tags map[string]bool
	if err := chatRef.Get(ctx, &tags); err != nil {
		log.Printf("Error getting tags")
		return []string{}, err
	}

	/* Get slice of tags */
	tagSlice := make([]string, 0)
	for tag, _ := range tags {
		tagSlice = append(tagSlice, tag)
	}
	return tagSlice, nil
}

func DeleteTag(update tgbotapi.Update, tag string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return err
	}

	/* Delete tag record */
	chatRef := client.NewRef("tags").Child(chatID)
	if err := chatRef.Child(tag).Delete(ctx); err != nil {
		log.Printf("Error deleting tag: %+v", err)
		return err
	}
	return nil
}

func updateTags(update tgbotapi.Update, tag string) error {
	/* If same tag won't update. Implicitly prevent double records */
	ctx := context.Background()
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return err
	}
	chatRef := client.NewRef("tags").Child(chatID)
	if err := chatRef.Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		log.Printf("Error updating tags: %+v", err)
		return err
	}

	return nil
}
