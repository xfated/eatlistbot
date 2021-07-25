package utils

import (
	"context"
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

/* ########## User State ##########*/
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

/* ########## Name (Init place) ##########*/
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

/* ########## Address ##########*/
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

/* ########## Notes ##########*/
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

/* ########## URL ##########*/
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

/* ########## Images ##########*/
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

/* ########## Tags ##########*/
func AddTempPlaceTag(update tgbotapi.Update, tag string) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
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
	userRef := client.NewRef("places").Child(chatID)
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

	// DEBUG
	// log.Printf("filterTags: %+v", filterTags)
	// log.Printf("placesList: %+v", placesList)

	return placesList, nil
}

/* ########## Add Place ##########*/
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

/* ########## Delete Place ##########*/
func SetMessageTarget(update tgbotapi.Update, messageID int) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set target */
	messageTargetRef := client.NewRef("users").Child(userID).Child("target").Child("message")
	if err := messageTargetRef.Update(ctx, map[string]interface{}{
		chatID: messageID,
	}); err != nil {
		return err
	}
	return nil
}

func GetMessageTarget(update tgbotapi.Update) (int, error) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserIDString(update)
	if err != nil {
		return 0, err
	}

	/* Get target */
	var target int
	messageTargetRef := client.NewRef("users").Child(userID).Child("target").Child("message")
	if err := messageTargetRef.Child(chatID).Get(ctx, &target); err != nil {
		return 0, err
	}
	return target, nil
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

/* ########## Query ##########*/
func ResetQuery(update tgbotapi.Update) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Delete query */
	queryRef := client.NewRef("users").Child(userID).Child("query")
	if err := queryRef.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// message should contain name
func SetQueryName(update tgbotapi.Update, name string) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	queryRef := client.NewRef("users").Child(userID).Child("query")
	if err := queryRef.Update(ctx, map[string]interface{}{
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

	queryRef := client.NewRef("users").Child(userID).Child("query").Child("name")
	var name string
	if err := queryRef.Get(ctx, &name); err != nil {
		return "", err
	}
	return name, err
}

func SetQueryNum(update tgbotapi.Update, num int) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set number of queries*/
	queryRef := client.NewRef("users").Child(userID).Child("query")
	if err := queryRef.Update(ctx, map[string]interface{}{
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

	queryRef := client.NewRef("users").Child(userID).Child("query")
	var queryNum int
	if err := queryRef.Child("queryNum").Get(ctx, &queryNum); err != nil {
		return 0, err
	}
	return queryNum, err
}

// message should contain tag
func AddQueryTag(update tgbotapi.Update, tag string) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Add tag */
	queryRef := client.NewRef("users").Child(userID).Child("query")
	if err := queryRef.Child("tags").Update(ctx, map[string]interface{}{
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
	queryRef := client.NewRef("users").Child(userID).Child("query")
	if err := queryRef.Child("tags").Get(ctx, &tagsMap); err != nil {
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

/* ########## Delete Place ##########*/
func SetPlaceTarget(update tgbotapi.Update, name string) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set target */
	placeTargetRef := client.NewRef("users").Child(userID).Child("target")
	if err := placeTargetRef.Update(ctx, map[string]interface{}{
		chatID: name,
	}); err != nil {
		return err
	}
	return nil
}

func GetPlaceTarget(update tgbotapi.Update) (string, error) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserIDString(update)
	if err != nil {
		return "", err
	}

	/* Get target */
	var target string
	placeTargetRef := client.NewRef("users").Child(userID).Child("target")
	if err := placeTargetRef.Child(chatID).Get(ctx, &target); err != nil {
		return "", err
	}
	return target, nil
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
