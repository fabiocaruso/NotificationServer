package grifts

import (
	"github.com/fabiocaruso/NotificationServer/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
