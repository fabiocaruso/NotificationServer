package actions

import (
	//"go.mongodb.org/mongo-driver/bson"
	"context"
	"log"
	"strconv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
)

type User struct {
	ID primitive.ObjectID 		`bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName string		`bson:"firstName" json:"firstName"`
	Name string			`bson:"name" json:"name"`
	Username string			`bson:"username" json:"username"`
	Hash string			`bson:"hash" json:"hash"`
	Devices []primitive.ObjectID	`bson:"devices" json:"devices"`
	Roles []primitive.ObjectID	`bson:"roles" json:"roles"`
}

func connDB(host string, port int) (*mongo.Client, bool) {
	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + strconv.Itoa(port))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, false
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, false
	}
	fmt.Println("Connected to MongoDB!")
	return client, true
}
