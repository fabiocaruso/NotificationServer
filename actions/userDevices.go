package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/fabiocaruso/NotificationServer/models"
	"github.com/fabiocaruso/NotificationServer/services"
	"gopkg.in/couchbase/gocb.v1"
	//"github.com/xeipuuv/gojsonschema"
	"github.com/Jeffail/gabs"
	"strings"
	"errors"
	"io/ioutil"
	"fmt"
)

type UserDevicesResource struct{}

func (user User) getDevices() (*[]models.Device, error) {
	//TODO: This line is crap
	userDevices := strings.Join(user.Devices[:], "', '")
	query := gocb.NewN1qlQuery(`
		SELECT ns.* FROM NotificationServer AS ns
		WHERE type = 'device'
		AND META().id IN ['` + userDevices + `']
		`)
	result, err := nsBucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return nil, err
	}
	var device models.Device
	devices := []models.Device{}
	for result.Next(&device) {
		devices = append(devices, device)
	}
	return &devices, nil
}

func (udr UserDevicesResource) Show(c buffalo.Context) error {
	user, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.Error(401, err)
	}
	devices, err := user.getDevices()
	if err != nil {
		return c.Error(500, err)
	}
	return c.Render(200, r.JSON(devices))
}

func (udr UserDevicesResource) Update(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}

func (udr UserDevicesResource) Add(c buffalo.Context) error {
	//TODO: Param checking
	//TODO: Check if a service is set and adjust the query
	
	// Check user auth token
	/*user, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.Error(401, err)
	}*/

	// Check validity of the input data
	req, _ := ioutil.ReadAll(c.Request().Body)
	/*schemaLoader := gojsonschema.NewReferenceLoader("https://github.com/fabiocaruso/NotificationServer/raw/master/validation/device.schema.json")
	documentLoader := gojsonschema.NewStringLoader(string(req))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return c.Error(400, err)
	}
	if !result.Valid() {
		errorString := ""
		for _, desc := range result.Errors() {
			errorString += desc.String() + "\n"
		}
		return c.Error(400, errors.New(errorString))
	}*/

	// Set webhooks of the services if needed
	device, _ := gabs.ParseJSON([]byte(req))
	for s, v := range device.S("services").ChildrenMap() {
		provider, _ := services.Providers[s]
		webhook, ok := provider.(services.Webhook)
		if !ok {
			return c.Error(500, errors.New("ERROR: assertion went wrong"))
		}
		err := webhook.SetWebhook(v.S("token").Data().(string))
		fmt.Println(err)
	}

	// Insert dataset in database
	/*query := gocb.NewN1qlQuery(`
	INSERT INTO NotificationServer (KEY, VALUE)
	VALUES (
		UUID(),
		` + req + `
	)
	RETURNING META().id AS ID
	`)
	params := map[string]interface{}{
		"name": c.Param("name"),
		"os": c.Param("os"),
		"telegramToken": c.Param("telegramToken"),
	}
	result, err := nsBucket.ExecuteN1qlQuery(query, params)
	if err != nil {
		return c.Error(500, err)
	}
	var row map[string]interface{}
	result.One(&row)
	fmt.Println(row)
	if result.Metrics().ResultCount != 1 || row["ID"] == "" {
		return c.Error(500, errors.New("Can't fetch device id!"))
	}
	query = gocb.NewN1qlQuery(`
	UPDATE NotificationServer
	SET devices = ARRAY_APPEND(devices, $deviceID)
	WHERE type = "user"
	AND meta().id = $userID
	`)
	params = map[string]interface{}{
		"deviceID": row["ID"],
		"userID": user.ID,
	}
	_, err = nsBucket.ExecuteN1qlQuery(query, params)
	if err != nil {
		return c.Error(500, err)
	}*/

	return c.Render(200, r.JSON("{'result': 'success'}"))
}

func (udr UserDevicesResource) Delete(c buffalo.Context) error {
	return c.Render(200, r.JSON("{'test': 'test'}"))
}
