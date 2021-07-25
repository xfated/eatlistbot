package utils

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/xfated/golistbot/services/constants"

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
func SetUserState(update tgbotapi.Update, state constants.State) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Update(ctx, map[string]interface{}{
		"state": strconv.Itoa(int(state)),
	}); err != nil {
		log.Println("Error setting state")
		return err
	}
	return nil
}

func GetUserState(update tgbotapi.Update) (constants.State, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return 0, err
	}
	userRef := client.NewRef("users").Child(userID)
	var stateString string
	if err := userRef.Child("state").Get(ctx, &stateString); err != nil {
		return 0, err
	}

	stateInt, err := strconv.Atoi(stateString)
	if err != nil {
		return 0, err
	}
	return constants.State(stateInt), err
}

/* Name (Also init place) */
func InitPlace(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	name, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Set(ctx, map[string]string{
		"name": name,
	}); err != nil {
		return err
	}
	return nil
}

/* Address */
func SetTempPlaceAddress(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	address, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Update(ctx, map[string]interface{}{
		"address": address,
	}); err != nil {
		return err
	}

	return nil
}

func UpdatePlaceAddress(update tgbotapi.Update, placeName, address string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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

/* Notes */
func SetTempPlaceNotes(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	notes, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Update(ctx, map[string]interface{}{
		"notes": notes,
	}); err != nil {
		return err
	}

	return nil
}

func UpdatePlaceNotes(update tgbotapi.Update, placeName, notes string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	placeRef := client.NewRef("places").Child(chatID)
	if err := placeRef.Child(placeName).Update(ctx, map[string]interface{}{
		"notes": notes,
	}); err != nil {
		return err
	}
	return nil
}

/* URL */
func SetTempPlaceURL(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	url, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Update(ctx, map[string]interface{}{
		"url": url,
	}); err != nil {
		return err
	}

	return nil
}

func UpdatePlaceURL(update tgbotapi.Update, placeName, url string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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
func AddTempPlaceImage(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	imageIDs, err := GetPhotoIDs(update)
	if err != nil {
		return err
	}
	imageID := imageIDs[3] // Take largest file size
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Child("images").Update(ctx, map[string]interface{}{
		imageID: true,
	}); err != nil {
		return err
	}

	return nil
}

func AddPlaceImage(update tgbotapi.Update, placeName, imageID string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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

func DeletePlaceImage(update tgbotapi.Update, placeName, imageID string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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
func AddTempPlaceTag(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	tag, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}

	return nil
}

func AddPlaceTag(update tgbotapi.Update, placeName, tag string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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

func DeletePlaceTag(update tgbotapi.Update, placeName, tag string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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
func GetPlaces(update tgbotapi.Update, filterTags map[string]bool) ([]constants.PlaceDetails, error) {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return []constants.PlaceDetails{}, err
	}

	/* get places */
	var places map[string]constants.PlaceDetails
	userRef := client.NewRef("place").Child(chatID)
	if err := userRef.Get(ctx, &places); err != nil {
		return []constants.PlaceDetails{}, err
	}
	placesList := make([]constants.PlaceDetails, len(places))
	i := 0
	for _, placeDetails := range places {
		placesList[i] = placeDetails
		i++
	}

	/* filter if tags are present */
	if len(filterTags) > 0 {
		filteredPlaces := make([]constants.PlaceDetails, 0)
		for _, place := range places {
			consider := false
			if place.Tags != nil {
				for tag := range place.Tags {
					/* select if any tag match */
					if filterTags[tag] {
						consider = true
						break
					}
				}
			}
			if consider {
				filteredPlaces = append(filteredPlaces, place)
			}
		}
		placesList = filteredPlaces
	}
	rand.Shuffle(len(placesList), func(i, j int) { placesList[i], placesList[j] = placesList[j], placesList[i] })
	return placesList, nil
}

/* Add / Delete places */
func GetTempPlace(update tgbotapi.Update) (constants.PlaceDetails, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return constants.PlaceDetails{}, err
	}

	var PlaceData constants.PlaceDetails
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("placeToAdd").Get(ctx, &PlaceData); err != nil {
		return constants.PlaceDetails{}, err
	}
	return PlaceData, nil
}

func GetPlace(update tgbotapi.Update, name string) (constants.PlaceDetails, error) {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return constants.PlaceDetails{}, err
	}

	/* Get place */
	var placeData constants.PlaceDetails
	chatRef := client.NewRef("places").Child(chatID)
	if err := chatRef.Child(name).Get(ctx, &placeData); err != nil {
		return constants.PlaceDetails{}, err
	}

	return placeData, nil
}

func AddPlace(update tgbotapi.Update, placeData constants.PlaceDetails) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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

	/* Add name to name collection */
	nameRef := client.NewRef("placeNames").Child(chatID)
	if err := nameRef.Update(ctx, map[string]interface{}{
		placeData.Name: true,
	}); err != nil {
		return err
	}
	return nil
}

func AddPlaceFromTemp(update tgbotapi.Update) (string, error) {
	// get from user details
	placeData, err := GetTempPlace(update)
	if err != nil {
		return "", err
	}
	// Add data to place
	if err := AddPlace(update, placeData); err != nil {
		return "", err
	}
	return placeData.Name, nil
}

func DeletePlace(update tgbotapi.Update, placeName string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}
	chatRef := client.NewRef("places").Child(chatID)
	if err := chatRef.Child(placeName).Delete(ctx); err != nil {
		return err
	}
	nameRef := client.NewRef("placeNames").Child(chatID)
	if err := nameRef.Child(placeName).Delete(ctx); err != nil {
		return err
	}
	return nil
}

func GetPlaceNames(update tgbotapi.Update) (map[string]bool, error) {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return map[string]bool{}, err
	}

	var placeNames map[string]bool
	nameRef := client.NewRef("placeNames").Child(chatID)
	if err := nameRef.Get(ctx, &placeNames); err != nil {
		return map[string]bool{}, err
	}

	return placeNames, nil
}

/* Read / Delete / Update tags */
func GetTags(update tgbotapi.Update) (map[string]bool, error) {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return map[string]bool{}, err
	}
	chatRef := client.NewRef("tags").Child(chatID)

	/* Retrieve tags */
	var tags map[string]bool
	if err := chatRef.Get(ctx, &tags); err != nil {
		return map[string]bool{}, err
	}

	// /* get slice of tags */
	// tagSlice := make([]string, len(tags))
	// i := 0
	// for tag := range tags {
	// 	tagSlice[i] = tag
	// 	i++
	// }
	return tags, nil
}

func DeleteTag(update tgbotapi.Update, tag string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
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
	chatID, _, err := GetChatUserIDString(update)
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

/* Handle queries */
func ResetQuery(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Delete query */
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Delete(ctx); err != nil {
		return err
	}
	return nil
}

// message should contain name
func SetQueryName(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	name, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Update(ctx, map[string]interface{}{
		"name": name,
	}); err != nil {
		return err
	}
	return nil
}

func GetQueryName(update tgbotapi.Update) (string, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return "", err
	}

	userRef := client.NewRef("users").Child(userID)
	var name string
	if err := userRef.Child("query").Get(ctx, &name); err != nil {
		return "", err
	}
	return name, err
}

/* Message should contain query type */
func SetQueryType(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	queryType, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Update(ctx, map[string]interface{}{
		"queryType": queryType,
	}); err != nil {
		return err
	}
	return nil
}

func GetQueryType(update tgbotapi.Update) (string, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return "", err
	}

	userRef := client.NewRef("users").Child(userID)
	var queryType string
	if err := userRef.Child("queryType").Get(ctx, &queryType); err != nil {
		return "", err
	}
	return queryType, err
}

func SetQueryNum(update tgbotapi.Update, num int) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set number of queries*/
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Update(ctx, map[string]interface{}{
		"queryNum": num,
	}); err != nil {
		return err
	}
	return nil
}

func GetQueryNum(update tgbotapi.Update) (int, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return 0, err
	}

	userRef := client.NewRef("users").Child(userID)
	var queryNum int
	if err := userRef.Child("queryNum").Get(ctx, &queryNum); err != nil {
		return 0, err
	}
	return queryNum, err
}

// message should contain yes/no
func SetQueryWithPics(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Check answer whether to get pic */
	toGetPic, _, err := GetMessage(update)
	if err != nil {
		return err
	}
	var getPic bool
	if toGetPic == "yes" {
		getPic = true
	} else if toGetPic == "no" {
		getPic = false
	} else {
		return errors.New("invalid answer")
	}

	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Update(ctx, map[string]interface{}{
		"getPic": getPic,
	}); err != nil {
		return err
	}
	return nil
}

func GetQueryWithPics(update tgbotapi.Update) (string, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return "", err
	}

	var getPic string
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Get(ctx, &getPic); err != nil {
		return "", err
	}
	return getPic, nil
}

// message should contain tag
func AddQueryTag(update tgbotapi.Update, tag string) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Add tag */
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}
	return nil
}

func GetQueryTags(update tgbotapi.Update) (map[string]bool, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return map[string]bool{}, err
	}

	var tagsMap map[string]bool
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("query").Child("tags").Get(ctx, &tagsMap); err != nil {
		return map[string]bool{}, err
	}

	// var tags = make([]string, len(tagsMap))
	// i := 0
	// for tag := range tagsMap {
	// 	tags[i] = tag
	// 	i++
	// }
	return tagsMap, nil
}
