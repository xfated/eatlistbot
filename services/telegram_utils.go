package services

import (
	"errors"
	"log"
	"os"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	baseURL            = "https://togolist-bot.herokuapp.com/"
	bot                *tgbotapi.BotAPI
)

/* Init */
func InitTelegram() {
	var err error

	// Init bot
	bot, err = tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Println(err)
	}

	// Set webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(baseURL + bot.Token))
	if err != nil {
		log.Println("Problem setting Webhook", err.Error())
	}

	log.Println("Loaded telegram bot")
}

/* General Logging */
func LogMessage(update tgbotapi.Update) {
	if update.Message != nil {
		log.Printf("Message: %+v", update.Message)
	}
}

func LogUpdate(update tgbotapi.Update) {
	log.Printf("Update: %+v", update)
}

/* Sending */
func SendMessage(update tgbotapi.Update, text string) error {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat & user id: %+v", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
	return nil
}

func SendStartInstructions(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Haha you sent start")
	bot.Send(msg)
	SetUserState(update, Idle)

	// Debug
	LogMessage(update)
}

func SendUnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command, please use /start for commands")
	bot.Send(msg)
}

func SendPhoto(update tgbotapi.Update, photoID string) error {
	// chatID, _, err := GetChatUserID(update)
	// if err != nil {
	// 	log.Printf("Error getting chat & user id: %+v", err)
	// 	return err
	// }
	var chatID int64 = -572862838
	tgbotapi.NewPhotoShare(chatID, photoID)
	log.Printf("shared photo: %v to chatID: %v", photoID, chatID)
	return nil
}

/* Getting */
func GetChatUserID(update tgbotapi.Update) (chatID int64, userID int, err error) {
	if update.Message == nil {
		chatID = 0
		userID = 0
		err = errors.New("invalid message")
		return
	}

	chatID = update.Message.Chat.ID
	userID = update.Message.From.ID
	err = nil
	return
}

func GetChatUserIDString(update tgbotapi.Update) (chatID, userID string, err error) {
	if update.Message == nil {
		chatID = ""
		userID = ""
		err = errors.New("invalid message")
		return
	}

	chatID = strconv.FormatInt(update.Message.Chat.ID, 10)
	userID = strconv.Itoa(update.Message.From.ID)
	err = nil
	return
}

func GetMessage(update tgbotapi.Update) (message string, messageID int, err error) {
	if update.Message == nil {
		message = ""
		messageID = 0
		err = errors.New("invalid message")
		return
	}

	message = update.Message.Text
	messageID = update.Message.MessageID
	err = nil
	return
}

func GetPhotoIDs(update tgbotapi.Update) ([]string, error) {
	if update.Message == nil {
		return []string{}, errors.New("invalid message")
	}

	if update.Message.Photo == nil {
		return []string{}, errors.New("no photo")
	}

	photoIDs := make([]string, 0)
	for _, photo := range *update.Message.Photo {
		photoIDs = append(photoIDs, photo.FileID)
	}
	log.Printf("Photo: %+v", update.Message.Photo)
	return photoIDs, nil
}

/* Confirm details and add to list */
