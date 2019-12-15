package services

import (
	"github.com/fabiocaruso/NotificationServer/models"
	"github.com/gobuffalo/buffalo"
)

var Services map[string]interface{} = map[string]interface{}{
	"telegram": Telegram{botToken: "lol"},
}

type Service interface {
	SendMessage(models.Device, string) error
	WebhookHandler(buffalo.Context) error
}
