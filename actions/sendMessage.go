package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/fabiocaruso/NotificationServer/services"
	"errors"
)

func sendMessageHandler(c buffalo.Context) error {
	user, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.Error(401, err)
	}
	devices, err := user.getDevices()
	if err != nil {
		return c.Error(500, err)
	}
	service, ok := services.Services["telegram"].(services.Service)
	if !ok {
		return c.Error(400, errors.New("Service not found!"))
	}
	for _, d := range *devices {
		if d.Name == c.Param("deviceName") {
			service.SendMessage(d, c.Param("text"))
		}
	}
	return c.Render(200, r.JSON("{'result': 'sucess'}"))
}
