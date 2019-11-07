package grifts

import (
	"github.com/fabiocaruso/NotificationServer/notification_server/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
