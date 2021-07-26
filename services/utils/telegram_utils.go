package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xfated/golistbot/services/constants"
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

/* Redirect */
func RedirectToBotChat(update *tgbotapi.Update, text string, url string) {
	redirectButton := tgbotapi.NewInlineKeyboardButtonURL("Add place", url)
	row := tgbotapi.NewInlineKeyboardRow(redirectButton)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	SendInlineKeyboard(update, text, inlineKeyboard)
}

// func RedirectToChat(update *tgbotapi.Update, text string) {
// 	redirectButton := tgbotapi.NewInlineKeyboardButtonSwitch("Exit to chats", "")
// 	row := tgbotapi.NewInlineKeyboardRow(redirectButton)

// 	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
// 	SendInlineKeyboard(update, text, inlineKeyboard)
// }

/* General Logging */
func LogMessage(update *tgbotapi.Update) {
	if update.Message != nil {
		log.Printf("Message: %+v", update.Message)
	}
}

func LogUpdate(update *tgbotapi.Update) {
	log.Printf("Update: %+v", update)
}

func LogCallbackQuery(update *tgbotapi.Update) {
	if update.CallbackQuery != nil {
		log.Printf("Callback Query: %+v", update.CallbackQuery)
		// log.Printf("Callback Query Message: %+v", update.CallbackQuery.Message)
	}
}

/* Deleting */
func DeleteMessage(chatID int64, messageID int) error {
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := bot.DeleteMessage(deleteConfig); err != nil {
		return err
	}
	return nil
}

/* Sending */
func SendMessage(update *tgbotapi.Update, text string) (*tgbotapi.Message, error) {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(chatID, text)
	message, err := bot.Send(msg)

	// Debug
	log.Printf("Sent message %s", text)
	return &message, err
}

func SendMessageForceReply(update *tgbotapi.Update, text string, messageID int) (*tgbotapi.Message, error) {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}
	message, err := bot.Send(msg)
	return &message, err
}

func SendMessageTargetChat(text string, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	return err
}

func SendUnknownCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command, please use /start for commands")
	bot.Send(msg)
}

func SendPhoto(update *tgbotapi.Update, photoID string) error {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		return err
	}
	photoConfig := tgbotapi.NewPhotoShare(chatID, photoID)
	bot.Send(photoConfig)
	return nil
}

func SendPlaceDetails(update *tgbotapi.Update, placeData constants.PlaceDetails, sendImage bool) {
	placeText := ""

	if placeData.Name != "" {
		placeText = placeText + fmt.Sprintf("Name: %s\n", placeData.Name)
	}
	if placeData.Address != "" {
		placeText = placeText + fmt.Sprintf("Address: %s\n", placeData.Address)
	}
	if placeData.Images != nil {
		placeText = placeText + fmt.Sprintf("Images: %v\n", len(placeData.Images))
	}
	if placeData.Tags != nil {
		tags := make([]string, len(placeData.Tags))
		i := 0
		for tag := range placeData.Tags {
			tags[i] = tag
			i++
		}
		tagText := strings.Join(tags, ", ")
		placeText = placeText + fmt.Sprintf("Tags: %s\n", tagText)
	}
	if placeData.Notes != "" {
		placeText = placeText + fmt.Sprintf("Notes: %s", placeData.Notes)
	}
	if placeData.URL != "" {
		// placeText = placeText + fmt.Sprintf("URL: %s\n", placeData.URL)

		/* To send as inline keyboard */
		redirectButton := tgbotapi.NewInlineKeyboardButtonURL(placeData.URL, placeData.URL)
		row := tgbotapi.NewInlineKeyboardRow(redirectButton)
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
		SendInlineKeyboard(update, placeText, inlineKeyboard)
	} else {
		SendMessage(update, placeText)
	}
	if sendImage && placeData.Images != nil {
		for imageID := range placeData.Images {
			SendPhoto(update, imageID)
		}
	}
}

func SetReplyMarkupKeyboard(update *tgbotapi.Update, text string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error GetChatUserID: %+v", err)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.BaseChat.ReplyMarkup = keyboard

	messageTarget := 0
	if update.Message != nil {
		messageTarget = update.Message.MessageID
	} else {
		messageTarget, err = GetMessageTarget(update)
		if err != nil {
			log.Printf("Error GetMessageTarget: %+v", err)
		}
	}
	msg.ReplyToMessageID = messageTarget
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error setting markup keyboard: %+v", err)
	}
}

func CreateAndSendInlineKeyboard(update *tgbotapi.Update, text string, col int, buttons ...string) *tgbotapi.Message {
	var buttonList []tgbotapi.InlineKeyboardButton

	for _, button := range buttons {
		buttonList = append(buttonList, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	// create rows
	for i := 0; i < len(buttonList); i += col {
		end := i + col
		if end > len(buttonList) {
			end = len(buttonList)
		}
		rows = append(rows, buttonList[i:end])
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return SendInlineKeyboard(update, text, inlineKeyboard)
}

func SendInlineKeyboard(update *tgbotapi.Update, text string, keyboard tgbotapi.InlineKeyboardMarkup) *tgbotapi.Message {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error getting chat user ID: %+v", err)
		return nil
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.BaseChat.ReplyMarkup = keyboard
	message, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error setting markup keyboard: %+v", err)
		return nil
	}
	return &message

}

func RemoveMarkupKeyboard(update *tgbotapi.Update, text string) *tgbotapi.Message {
	chatID, _, err := GetChatUserID(update)
	if err != nil {
		log.Printf("Error GetChatUserID: %+v", err)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
	removeKeyboard.Selective = true
	msg.BaseChat.ReplyMarkup = removeKeyboard

	messageTarget := 0
	if update.Message != nil {
		messageTarget = update.Message.MessageID
	} else {
		messageTarget, err = GetMessageTarget(update)
		if err != nil {
			log.Printf("Error GetMessageTarget: %+v", err)
		}
	}
	msg.ReplyToMessageID = messageTarget
	message, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error removing markup keyboard: %+v", err)
	}
	return &message
}

/* Getting */
func GetChatUserID(update *tgbotapi.Update) (chatID int64, userID int, err error) {
	if update.Message != nil {
		chatID = update.Message.Chat.ID
		userID = update.Message.From.ID
		err = nil
		return
	}
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
		userID = update.CallbackQuery.From.ID
		err = nil
		return
	}

	chatID = 0
	userID = 0
	err = errors.New("invalid message or callback query")
	return
}

func GetChatUserIDString(update *tgbotapi.Update) (chatID, userID string, err error) {
	if update.Message != nil {
		chatID = strconv.FormatInt(update.Message.Chat.ID, 10)
		userID = strconv.Itoa(update.Message.From.ID)
		err = nil
		return
	}
	if update.CallbackQuery != nil {
		chatID = strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
		userID = strconv.Itoa(update.CallbackQuery.From.ID)
		err = nil
		return
	}

	chatID = ""
	userID = ""
	err = errors.New("invalid message or callback query")
	return
}

func GetMessage(update *tgbotapi.Update) (message string, messageID int, err error) {
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

func GetCallbackQueryMessage(update *tgbotapi.Update) (string, error) {
	if update.CallbackQuery == nil {
		return "", errors.New("invalid callback data")
	}
	return update.CallbackQuery.Data, nil
}

func GetPhotoIDs(update *tgbotapi.Update) ([]string, error) {
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
	return photoIDs, nil
}

/* Check string */
func CheckForSlash(update *tgbotapi.Update) error {
	if update.Message != nil {
		message, _, err := GetMessage(update)
		if err != nil {
			return err
		}
		if strings.Contains(message, "/") {
			SendMessage(update, "Don't use / here! I will get confused :(")
			SendMessage(update, "Please resend with a proper message")
			return errors.New("slash in message")
		}
		return nil
	}
	// if update.CallbackQuery != nil {
	// 	message, err := GetCallbackQueryMessage(update)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if strings.Contains(message, "/") {
	// 		SendMessage(update, "Please don't use / in your input! I will get confused :(")
	// 		SendMessage(update, "Please resend with a proper message")
	// 		return errors.New("slash in message")
	// 	}
	// 	return nil
	// }
	return nil
}
