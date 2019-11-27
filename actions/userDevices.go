package actions

import "github.com/gobuffalo/buffalo"

func userDevicesHandler(c buffalo.Context) error {
	return c.Render(200, r.String("Works: " + c.Param("apikey")))
}
