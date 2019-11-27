package actions

import (
	"github.com/gobuffalo/buffalo"
	"time"
	"errors"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/pascaldekloe/jwt"
	//"crypto/ed25519"
)

type tokenPayload struct {
	ID primitive.ObjectID
}

func generateHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func validateToken(token string) (*primitive.ObjectID, error) {
	claims, err := jwt.EdDSACheck([]byte(token), publicKey)
	if err != nil {
		return &primitive.NilObjectID, errors.New("API key denied: " + err.Error())
	}
	if !claims.Valid(time.Now()) {
		return &primitive.NilObjectID, errors.New("API key expired! Please sign in again.")
	}
	objectID, err := primitive.ObjectIDFromHex(claims.ID)
	return &objectID, err
}

func getUserFromRH(header string) (*User, error) {
	token := header[7:] // Trim Bearer
	if len(token) == 0 {
		return &User{}, errors.New("No API key set!")
	}
	userID, err := validateToken(token)
	if err != nil {
		return &User{}, err
	}
	//return c.Render(200, r.JSON(userid))
	var user User
	users := nsdb.Collection("users")
	err = users.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return &User{}, errors.New("user not found!")
	}
	return &user, nil
}

func authHandler(c buffalo.Context) error {
	var user User
	username := c.Param("username")
	password := c.Param("password")
	hash := generateHash(password)
	err := nsdb.Collection("users").FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil && (user.ID == primitive.ObjectID{}) {
		return c.Error(401, errors.New("Something went wrong!"))
	}
	if hash == user.Hash {
		var claims jwt.Claims
		claims.ID = user.ID.Hex()
		claims.Subject = user.FirstName + " " + user.Name
		now := time.Now().Round(time.Second)
		claims.Issued = jwt.NewNumericTime(now)
		claims.Expires = jwt.NewNumericTime(now.Add(30000*time.Hour))
		token, err := claims.EdDSASign(privateKey)
		if err != nil {
			return c.Error(401, errors.New("Could not sign API key! " + err.Error()))
		}
		return c.Render(200, r.JSON("SUCCESS! Your API Token is: " + string(token)))
	} else {
		return c.Error(401, errors.New("Your password is wrong!"))
	}
}
