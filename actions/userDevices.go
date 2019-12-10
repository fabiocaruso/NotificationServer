package actions

import (
	"github.com/gobuffalo/buffalo"
	"gopkg.in/couchbase/gocb.v1"
	"strings"
	//"errors"
)

type UserDevicesResource struct{}

func (udr UserDevicesResource) Show(c buffalo.Context) error {
	user, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.Error(401, err)
	}
	userDevices := strings.Join(user.Devices[:], "', '")
	query := gocb.NewN1qlQuery(`
		SELECT ns.name FROM NotificationServer AS ns
		WHERE type = 'device'
		AND META().id IN ['` + userDevices + `']
		`)
	result, err := nsBucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return c.Error(500, err)
	}
	var device Device
	devices := []Device{}
	for result.Next(&device) {
		devices = append(devices, device)
	}
	return c.Render(200, r.JSON(devices))
}

func (udr UserDevicesResource) Update(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}

func (udr UserDevicesResource) Add(c buffalo.Context) error {
	//TODO: Param checking
	//TODO: Check if a service is set and adjust the query
	//TODO: Add device to users devices list in database
	_, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.Error(401, err)
	}
	query := gocb.NewN1qlQuery(`
		INSERT INTO NotificationServer (KEY, VALUE)
		VALUES (
	    		UUID(),
	        	{
				"type": "device",
		    		"name": $name,
		        	"os": $os,
			    	"services": {
					"telegram": {
						"botToken": $telegramToken
					}
				}
			}
		)
		`)
	params := map[string]interface{}{
		"name": c.Param("name"),
		"os": c.Param("os"),
		"telegramToken": c.Param("telegramToken"),
	}
	_, err = nsBucket.ExecuteN1qlQuery(query, params)
	if err != nil {
		return c.Error(500, err)
	}
	return c.Render(200, r.JSON("{'result': 'success'}"))
}

func (udr UserDevicesResource) Delete(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}
