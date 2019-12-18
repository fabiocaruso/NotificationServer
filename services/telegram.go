package services

import (
	"github.com/fabiocaruso/NotificationServer/models"
	"github.com/gobuffalo/buffalo"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// Required struct, could also be an empty struct
type Telegram struct {
	botToken string
}

func getOptions(device models.Device) map[string]string {
	//TODO: options validation
	options := make(map[string]string)
	for k, v := range device.Services["telegram"].(map[string]interface{}) {
		options[k] = v.(string)
	}
	return options
}

// Required
func (t Telegram) SendMessage(devices []models.Device, text string) error {
	for _, device := range devices {
		options := getOptions(device)
		bot, err := tgbotapi.NewBotAPI(options["botToken"])
		if err != nil {
			return err
		}
		chatID, err := strconv.ParseInt(options["chatId"], 10, 64)
		if err != nil {
			return err
		}
		bot.Send(tgbotapi.NewMessage(chatID, text))
	}
	return nil
}

// Required if Webhook is needed
func (t Telegram) SetWebhook(device models.Device, url string) error {
	options := getOptions(device)
	bot, err := tgbotapi.NewBotAPI(options["botToken"])
	if err != nil {
		return err
	}
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("http://the-domain.tld/" + options["botToken"]))
	if err != nil {
		return err
	}
	return nil
}

// Required if Webhook is needed
func (t Telegram) WebhookHandler(c buffalo.Context) error {
	bytes, _ := ioutil.ReadAll(c.Request().Body)
	c.Request().Body.Close()
	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)
	fmt.Println(update)
	return nil
}
