package services

import (
	"github.com/fabiocaruso/NotificationServer/models"
	"github.com/gobuffalo/buffalo"
)

var Services map[string]interface{} = map[string]interface{}{
	"telegram": Telegram{botToken: "lol"},
}

type Service interface {
	SendMessage([]models.Device, string) error
}

// Optional Webhook interface
type Webhook interface {
	WebhookHandler(buffalo.Context) error
	SetWebhook(string) error
}
