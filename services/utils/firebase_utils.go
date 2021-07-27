package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

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
	ao := map[string]interface{}{"uid": "togolistbot"}
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
func SetUserState(update *tgbotapi.Update, state constants.State) error {
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

func GetUserState(update *tgbotapi.Update) (constants.State, error) {
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

/* ########## Name (Init item) ##########*/
func InitItem(update *tgbotapi.Update) error {
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
	if err := userRef.Child("itemToAdd").Set(ctx, map[string]string{
		"name": name,
	}); err != nil {
		return err
	}
	return nil
}

/* ########## Address ##########*/
func SetTempItemAddress(update *tgbotapi.Update) error {
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
	if err := userRef.Child("itemToAdd").Update(ctx, map[string]interface{}{
		"address": address,
	}); err != nil {
		return err
	}

	return nil
}

func UpdateItemAddress(update *tgbotapi.Update, itemName, address string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Update(ctx, map[string]interface{}{
		"address": address,
	}); err != nil {
		return err
	}
	return nil
}

/* ########## Notes ##########*/
func SetTempItemNotes(update *tgbotapi.Update) error {
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
	if err := userRef.Child("itemToAdd").Update(ctx, map[string]interface{}{
		"notes": notes,
	}); err != nil {
		return err
	}

	return nil
}

func UpdateItemNotes(update *tgbotapi.Update, itemName, notes string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Update(ctx, map[string]interface{}{
		"notes": notes,
	}); err != nil {
		return err
	}
	return nil
}

/* ########## URL ##########*/
func SetTempItemURL(update *tgbotapi.Update) error {
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
	if err := userRef.Child("itemToAdd").Update(ctx, map[string]interface{}{
		"url": url,
	}); err != nil {
		return err
	}

	return nil
}

func UpdateItemURL(update *tgbotapi.Update, itemName, url string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Update(ctx, map[string]interface{}{
		"url": url,
	}); err != nil {
		return err
	}
	return nil
}

/* ########## Images ##########*/
func AddTempItemImage(update *tgbotapi.Update) error {
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
	imageID := imageIDs[len(imageIDs)-1] // Take largest file size
	if err != nil {
		return err
	}
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("itemToAdd").Child("images").Update(ctx, map[string]interface{}{
		imageID: true,
	}); err != nil {
		return err
	}

	return nil
}

func AddItemImage(update *tgbotapi.Update, itemName, imageID string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Child("tags").Update(ctx, map[string]interface{}{
		imageID: true,
	}); err != nil {
		return err
	}
	return nil
}

func DeleteItemImage(update *tgbotapi.Update, itemName, imageID string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Child("tags").Child(imageID).Delete(ctx); err != nil {
		return err
	}
	return nil
}

/* ########## Tags ##########*/
func AddTempItemTag(update *tgbotapi.Update, tag string) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("itemToAdd").Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}

	return nil
}

func GetTempItemTags(update *tgbotapi.Update) (map[string]bool, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return nil, err
	}

	var tagsMap map[string]bool

	/* Set temp under userRef */
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("itemToAdd").Child("tags").Get(ctx, &tagsMap); err != nil {
		return nil, err
	}

	return tagsMap, nil
}

func DeleteTempItemTag(update *tgbotapi.Update, tag string) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set temp under userRef */
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("itemToAdd").Child("tags").Child(tag).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func AddItemTag(update *tgbotapi.Update, itemName, tag string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Child("tags").Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}
	return nil
}

func DeleteItemTag(update *tgbotapi.Update, itemName, tag string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	itemRef := client.NewRef("items").Child(chatID)
	if err := itemRef.Child(itemName).Child("tags").Child(tag).Delete(ctx); err != nil {
		return err
	}
	return nil
}

/* get list of items */
func GetItems(update *tgbotapi.Update, filterTags map[string]bool) ([]constants.ItemDetails, error) {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return []constants.ItemDetails{}, err
	}

	/* get items */
	var items map[string]constants.ItemDetails
	userRef := client.NewRef("items").Child(chatID)
	if err := userRef.Get(ctx, &items); err != nil {
		return []constants.ItemDetails{}, err
	}
	itemsList := make([]constants.ItemDetails, len(items))
	i := 0
	for _, itemDetails := range items {
		itemsList[i] = itemDetails
		i++
	}

	// To delete any tags that have not been used
	tagsList := make([]string, len(filterTags))
	tagUsed := make([]bool, len(filterTags))
	i = 0
	for tag := range filterTags {
		tagsList[i] = tag
		tagUsed[i] = false
		i++
	}

	/* filter if tags are present */
	if len(filterTags) > 0 {
		filteredItems := make([]constants.ItemDetails, 0)
		// For each item
		for _, item := range items {
			consider := false
			// If have tags
			if item.Tags != nil {
				// Check if match any filter
				for idx, tag := range tagsList {
					/* select if any tag match */
					if item.Tags[tag] {
						consider = true
						tagUsed[idx] = true
						break
					}
				}
			}
			if consider {
				filteredItems = append(filteredItems, item)
			}
		}
		itemsList = filteredItems
	}

	unusedTags := make([]string, 0)
	for idx, tag := range tagsList {
		if !tagUsed[idx] {
			unusedTags = append(unusedTags, tag)
		}
	}
	if len(unusedTags) > 0 {
		tagsString := strings.Join(unusedTags, ", ")
		SendMessage(update, fmt.Sprintf("There is no item with these tags: %s.\nSo imma delete them", tagsString), false)
		for _, tag := range unusedTags {
			DeleteTag(update, tag)
		}
	}

	/* Shuffle for random */
	rand.Shuffle(len(itemsList), func(i, j int) { itemsList[i], itemsList[j] = itemsList[j], itemsList[i] })

	// DEBUG
	// log.Printf("filterTags: %+v", filterTags)
	// log.Printf("itemsList: %+v", itemsList)

	return itemsList, nil
}

/* ########## Add Item ##########*/
func SetChatTarget(update *tgbotapi.Update, chatID int64) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set target */
	chatTargetRef := client.NewRef("users").Child(userID).Child("target").Child("chat")
	if err := chatTargetRef.Set(ctx, chatID); err != nil {
		return err
	}
	return nil
}

func GetChatTarget(update *tgbotapi.Update) (int64, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return 0, err
	}

	/* Get target */
	var target int64
	chatTargetRef := client.NewRef("users").Child(userID).Child("target").Child("chat")
	if err := chatTargetRef.Get(ctx, &target); err != nil {
		return 0, err
	}
	return target, nil
}

func GetTempItem(update *tgbotapi.Update) (constants.ItemDetails, error) {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return constants.ItemDetails{}, err
	}

	var ItemData constants.ItemDetails
	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("itemToAdd").Get(ctx, &ItemData); err != nil {
		return constants.ItemDetails{}, err
	}
	return ItemData, nil
}

func GetItem(update *tgbotapi.Update, name string, chatID string) (constants.ItemDetails, error) {
	ctx := context.Background()

	/* Get item */
	var itemData constants.ItemDetails
	chatRef := client.NewRef("items").Child(chatID)
	if err := chatRef.Child(name).Get(ctx, &itemData); err != nil {
		return constants.ItemDetails{}, err
	}

	return itemData, nil
}

func AddItem(update *tgbotapi.Update, itemData constants.ItemDetails, chatID string) error {
	ctx := context.Background()

	/* Add item to item collection */
	chatRef := client.NewRef("items").Child(chatID)
	if err := chatRef.Child(itemData.Name).Set(ctx, itemData); err != nil {
		return err
	}
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}

	err = SendMessageTargetChat(fmt.Sprintf("%s has been added/edited", itemData.Name), chatIDInt, false)
	if err != nil {
		log.Printf("error SendMessageTargetChat: %+v", err)
	}

	/* Add tags to tag collection */
	for tag := range itemData.Tags {
		if err := updateTags(update, chatID, tag); err != nil {
			return err
		}
	}

	/* Add name to name collection */
	nameRef := client.NewRef("itemNames").Child(chatID)
	if err := nameRef.Update(ctx, map[string]interface{}{
		itemData.Name: true,
	}); err != nil {
		return err
	}
	return nil
}

func AddItemFromTemp(update *tgbotapi.Update, chatID string) (string, error) {
	// get from user details
	itemData, err := GetTempItem(update)
	if err != nil {
		return "", err
	}
	// Add data to item
	if err := AddItem(update, itemData, chatID); err != nil {
		return "", err
	}
	return itemData.Name, nil
}

/* ########## Delete Item ##########*/
func SetMessageTarget(update *tgbotapi.Update, messageID int) error {
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

func GetMessageTarget(update *tgbotapi.Update) (int, error) {
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

func GetItemNames(update *tgbotapi.Update, chatID string) (map[string]bool, error) {
	ctx := context.Background()

	var itemNames map[string]bool
	nameRef := client.NewRef("itemNames").Child(chatID)
	if err := nameRef.Get(ctx, &itemNames); err != nil {
		return map[string]bool{}, err
	}

	return itemNames, nil
}

/* Read / Delete / Update tags */
func GetTags(update *tgbotapi.Update, chatID string) (map[string]bool, error) {
	ctx := context.Background()
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

func DeleteTag(update *tgbotapi.Update, tag string) error {
	/* Only can delete from within chat */
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

func updateTags(update *tgbotapi.Update, chatID string, tag string) error {
	/* Can update when adding item from bot privat chat */
	/* If same tag won't update. Implicitly prevent double records */
	ctx := context.Background()

	chatRef := client.NewRef("tags").Child(chatID)
	if err := chatRef.Update(ctx, map[string]interface{}{
		tag: true,
	}); err != nil {
		return err
	}

	return nil
}

/* ########## Query ##########*/
func ResetQuery(update *tgbotapi.Update) error {
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
func SetQueryName(update *tgbotapi.Update, name string) error {
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

func GetQueryName(update *tgbotapi.Update) (string, error) {
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

func SetQueryNum(update *tgbotapi.Update, num int) error {
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

func GetQueryNum(update *tgbotapi.Update) (int, error) {
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
func AddQueryTag(update *tgbotapi.Update, tag string) error {
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

func GetQueryTags(update *tgbotapi.Update) (map[string]bool, error) {
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

func AddMessageToDelete(update *tgbotapi.Update, message *tgbotapi.Message) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.New("nil message")
	}
	messageID := strconv.Itoa(message.MessageID)

	/* Add recent message data */
	recentDeleteRef := client.NewRef("deleteRecord").Child(chatID)
	if err := recentDeleteRef.Child("messages").Update(ctx, map[string]interface{}{
		messageID: true,
	}); err != nil {
		return err
	}
	return nil
}

func ResetMessagesToDelete(update *tgbotapi.Update) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Delete recent message data */
	recentDeleteRef := client.NewRef("deleteRecord").Child(chatID)
	if err := recentDeleteRef.Child("messages").Delete(ctx); err != nil {
		return err
	}
	return nil
}

func DeleteRecentMessages(update *tgbotapi.Update) error {
	ctx := context.Background()
	chatIDString, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return err
	}

	/* Get recent message data */
	recentDeleteRef := client.NewRef("deleteRecord").Child(chatIDString)
	var messageIDs map[string]bool
	if err := recentDeleteRef.Child("messages").Get(ctx, &messageIDs); err != nil {
		return err
	}
	// Delete one by one
	for messageIDString := range messageIDs {
		messageID, err := strconv.Atoi(messageIDString)
		if err != nil {
			return err
		}

		if err := DeleteMessage(chatID, messageID); err != nil {
			log.Printf("Error DeleteMessage: %+v", err)
			// return err
		}
	}
	if err := ResetMessagesToDelete(update); err != nil {
		return err
	}
	return nil
}

/* ########## Delete Item ##########*/
func SetItemTarget(update *tgbotapi.Update, name string) error {
	ctx := context.Background()
	chatID, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	/* Set target */
	itemTargetRef := client.NewRef("users").Child(userID).Child("target").Child("item")
	if err := itemTargetRef.Update(ctx, map[string]interface{}{
		chatID: name,
	}); err != nil {
		return err
	}
	return nil
}

func GetItemTarget(update *tgbotapi.Update) (string, error) {
	ctx := context.Background()
	chatID, userID, err := GetChatUserIDString(update)
	if err != nil {
		return "", err
	}

	/* Get target */
	var target string
	itemTargetRef := client.NewRef("users").Child(userID).Child("target").Child("item")
	if err := itemTargetRef.Child(chatID).Get(ctx, &target); err != nil {
		return "", err
	}
	return target, nil
}

func DeleteItem(update *tgbotapi.Update, itemName string) error {
	ctx := context.Background()
	chatID, _, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}
	chatRef := client.NewRef("items").Child(chatID)
	if err := chatRef.Child(itemName).Delete(ctx); err != nil {
		return err
	}
	nameRef := client.NewRef("itemNames").Child(chatID)
	if err := nameRef.Child(itemName).Delete(ctx); err != nil {
		return err
	}
	return nil
}

/* ########## Edit Item ##########*/
func AddItemToTemp(update *tgbotapi.Update, itemData constants.ItemDetails) error {
	ctx := context.Background()
	_, userID, err := GetChatUserIDString(update)
	if err != nil {
		return err
	}

	userRef := client.NewRef("users").Child(userID)
	if err := userRef.Child("itemToAdd").Set(ctx, itemData); err != nil {
		return err
	}
	return nil
}

func CopyItemToTempItem(update *tgbotapi.Update, itemName string, chatID string) error {
	itemData, err := GetItem(update, itemName, chatID)
	if err != nil {
		return err
	}
	if err := AddItemToTemp(update, itemData); err != nil {
		return err
	}
	return nil
}

/* ########## Feedback ##########*/
func AddFeedback(update *tgbotapi.Update) {
	ctx := context.Background()

	currentTime := time.Now()
	date := currentTime.Format("01-02-2006")
	feedback := update.Message.Text
	user := update.Message.From
	username := user.UserName
	userid := user.ID
	chatid := update.Message.Chat.ID

	feedbackRef := client.NewRef("feedback").Child(date)
	if _, err := feedbackRef.Push(ctx, map[string]interface{}{
		"username": username,
		"userid":   userid,
		"chatid":   chatid,
		"feedback": feedback,
	}); err != nil {
		log.Printf("error push feedback: %+v", err)
	}
}
