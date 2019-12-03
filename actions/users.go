package actions

import (
	"github.com/gobuffalo/buffalo"
	"errors"
)

type UsersResource struct{}

func (usersr UsersResource) Show(c buffalo.Context) error {
	user, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.Error(401, errors.New("ERROR: " + err.Error()))
	}
	return c.Render(200, r.JSON(user))
}

func (usersr UsersResource) Update(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}

func (usersr UsersResource) Add(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}

func (usersr UsersResource) Delete(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}
