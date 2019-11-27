package actions

import (
	"github.com/gobuffalo/buffalo"
	//"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
)

type UsersResource struct{}

func (usersr UsersResource) Show(c buffalo.Context) error {
	user, err := getUserFromRH(c.Request().Header.Get("Authorization"))
	if err != nil {
		c.Error(401, err)
	}
	query := []bson.D{
		bson.D{"$lookup", bson.D{
			{"from", "roles"},
			{"localField", "roles"},
			{"foreignField", "_id"},
			{"as", "rolesData"},
		}},
		/*bson.D{"$match", bson.D{
			{"from", "devices"},
			{"localField", "devices"},
			{"foreignField", "_id"},
			{"as", "devicesData"},
		}},*/
	}
	var cursor *mongo.Cursor
	cursor, err = nsdb.Collection("users").Aggregate(context.TODO(), mongo.Pipeline{query})
	fmt.Println("ERROR: ", err)
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}
	for _, result := range results {
		fmt.Printf("RESULT: \n", result)
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
