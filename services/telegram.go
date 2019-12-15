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

type Telegram struct {
	botToken string
}

func (t Telegram) SendMessage(device models.Device, text string) error {
	//TODO: options validation
	options := make(map[string]string)
	for k, v := range device.Services["telegram"].(map[string]interface{}) {
		options[k] = v.(string)
	}
	bot, err := tgbotapi.NewBotAPI(options["botToken"])
	if err != nil {
		return err
	}
	chatID, err := strconv.ParseInt(options["chatId"], 10, 64)
	if err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(chatID, text))
	return nil
}

func (t Telegram) WebhookHandler(c buffalo.Context) error {
	bytes, _ := ioutil.ReadAll(c.Request().Body)
	c.Request().Body.Close()
	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)
	fmt.Println(update)
	return nil
}
