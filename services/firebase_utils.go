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
func setUserState(update tgbotapi.Update, state State) error {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child(chatID).Update(ctx, map[string]interface{}{
		"state": strconv.Itoa(int(state)),
	}); err != nil {
		log.Println("Error setting state")
		return err
	}
	return nil
}

func getUserState(update tgbotapi.Update) (State, error) {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return 0, err
	}
	userRef := client.NewRef("users").Child(userID)
	var stateString string
	if err := userRef.Child(chatID).Child("state").Get(ctx, &stateString); err != nil {
		return 0, err
	}

	stateInt, err := strconv.Atoi(stateString)
	if err != nil {
		return 0, err
	}
	return State(stateInt), err
}

/* Name (Also init place) */
func initPlace(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	name, _, err := getMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("placeToAdd").Set(ctx, map[string]string{
		"name": name,
	}); err != nil {
		return err
	}
	return nil
}

/* Address */
func setTempPlaceAddress(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	address, _, err := getMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("placeToAdd").Update(ctx, map[string]interface{}{
		"address": address,
	}); err != nil {
		return err
	}

	return nil
}

func updatePlaceAddress(update tgbotapi.Update, placeName, address string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Update(ctx, map[string]interface{}{
		"address": address,
	}); err != nil {
		return err
	}
	return nil
}

/* URL */
func setTempPlaceURL(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	url, _, err := getMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("placeToAdd").Update(ctx, map[string]interface{}{
		"url": url,
	}); err != nil {
		return err
	}

	return nil
}

func updatePlaceURL(update tgbotapi.Update, placeName, url string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Update(ctx, map[string]interface{}{
		"url": url,
	}); err != nil {
		return err
	}
	return nil
}

/* Images */
func addTempPlaceImage(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	imageIDs, err := getPhotoIDs(update)
	if err != nil {
		return err
	}
	imageID := imageIDs[3] // Take largest file size
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("placeToAdd").Child("images").Update(ctx, map[string]interface{}{
		imageID: true,
	}); err != nil {
		return err
	}

	return nil
}

func addPlaceImage(update tgbotapi.Update, placeName, imageID string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Child("tags").Update(ctx, map[string]interface{}{
		imageID: true,
	}); err != nil {
		return err
	}
	return nil
}

func deletePlaceImage(update tgbotapi.Update, placeName, imageID string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Child("tags").Child(imageID).Delete(ctx); err != nil {
		return err
	}
	return nil
}

/* Tags */
func addTempPlaceTag(update tgbotapi.Update) error {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	tag, _, err := getMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("placeToAdd").Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}

	return nil
}

func addPlaceTag(update tgbotapi.Update, placeName, tag string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}
	return nil
}

func deletePlaceTag(update tgbotapi.Update, placeName, tag string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Child("tags").Child(tag).Delete(ctx); err != nil {
		return err
	}
	return nil
}

/* get list of places */
func getPlaces(update tgbotapi.Update) (map[string]PlaceDetails, error) {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return map[string]PlaceDetails{}, err
	}

	/* get places */
	var Places map[string]PlaceDetails
	userRef := client.NewRef("place").Child(chatID)
	if err := userRef.Get(ctx, &Places); err != nil {
		return map[string]PlaceDetails{}, err
	}
	return Places, nil
}

/* Add / Delete places */
func getTempPlace(update tgbotapi.Update) (PlaceDetails, error) {
	ctx := context.Background()
	chatID, userID, err := getChatUserIDString(update)
	if err != nil {
		return PlaceDetails{}, err
	}

	var PlaceData PlaceDetails
	userRef := client.NewRef("users").Child(userID).Child(chatID)
	if err := userRef.Child("placeToAdd").Get(ctx, &PlaceData); err != nil {
		return PlaceDetails{}, err
	}
	return PlaceData, nil
}

func addPlace(update tgbotapi.Update, placeData PlaceDetails) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Add place to place collection */
	chatRef := client.NewRef("places").Child(chatID)
	if err := chatRef.Child(placeData.Name).Set(ctx, placeData); err != nil {
		return err
	}

	/* Add tags to tag collection */
	for tag := range placeData.Tags {
		if err := updateTags(update, tag); err != nil {
			return err
		}
	}
	return nil
}

func addPlaceFromTemp(update tgbotapi.Update) (string, error) {
	// get from user details
	placeData, err := getTempPlace(update)
	if err != nil {
		return "", err
	}
	// Add data to place
	if err := addPlace(update, placeData); err != nil {
		return "", err
	}
	return placeData.Name, nil
}

func deletePlace(update tgbotapi.Update, placeName string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}
	chatRef := client.NewRef("places").Child(chatID)
	if err := chatRef.Child(placeName).Delete(ctx); err != nil {
		return err
	}
	return nil
}

/* Read / Delete / Update tags */
func getTags(update tgbotapi.Update) ([]string, error) {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return []string{}, err
	}
	chatRef := client.NewRef("tags").Child(chatID)

	/* Retrieve tags */
	var tags map[string]bool
	if err := chatRef.Get(ctx, &tags); err != nil {
		return []string{}, err
	}

	/* get slice of tags */
	tagSlice := make([]string, 0)
	for tag, _ := range tags {
		tagSlice = append(tagSlice, tag)
	}
	return tagSlice, nil
}

func deleteTag(update tgbotapi.Update, tag string) error {
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Delete tag record */
	chatRef := client.NewRef("tags").Child(chatID)
	if err := chatRef.Child(tag).Delete(ctx); err != nil {
		return err
	}
	return nil
}

func updateTags(update tgbotapi.Update, tag string) error {
	/* If same tag won't update. Implicitly prevent double records */
	ctx := context.Background()
	chatID, _, err := getChatUserIDString(update)
	if err != nil {
		return err
	}
	chatRef := client.NewRef("tags").Child(chatID)
	if err := chatRef.Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}

	return nil
}
