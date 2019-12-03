package actions

import (
	"github.com/gobuffalo/buffalo"
	"gopkg.in/couchbase/gocb.v1"
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	//"crypto/ed25519"
	"github.com/pascaldekloe/jwt"
)

func generateHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func validateToken(token string) (string, error) {
	claims, err := jwt.EdDSACheck([]byte(token), publicKey)
	if err != nil {
		return "", errors.New("API key denied: " + err.Error())
	}
	if !claims.Valid(time.Now()) {
		return "", errors.New("API key not valid! Please sign in again.")
	}
	return claims.ID, nil
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
	bucket, err := getBucket(nsdb, "NotificationServer")
	if err != nil {
		return &User{}, err
	}
	query := gocb.NewN1qlQuery(`
		SELECT ns.*, meta().id 
		FROM NotificationServer ns
		WHERE type = 'user' 
			AND meta().id = '` + userID + `'
		LIMIT 1`)
	result, err := bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return &User{}, err
	}
	var user User
	result.One(&user)
	if result.Metrics().ResultCount != 1 || user.ID == "" {
		return &User{}, errors.New("User not found!")
	}
	return &user, nil
}

func authHandler(c buffalo.Context) error {
	//var user User
	//TODO: get username and password from body
	username := c.Param("username")
	email := c.Param("email")
	password := c.Param("password")
	hash := generateHash(password)
	bucket, err := getBucket(nsdb, "NotificationServer")
	if err != nil {
		return c.Error(401, err)
	}
	query := gocb.NewN1qlQuery(`
		SELECT type, meta().id, username, email, ` + "`hash`" + ` 
		FROM NotificationServer ns
		WHERE type = 'user' 
			AND (username='` + username + `' 
			OR email='` + email + `') 
		LIMIT 1`)
	result, err := bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return c.Error(401, errors.New("Cannot execute query: " + err.Error()))
	}
	var user User
	result.One(&user)
	if hash == user.Hash {
		var claims jwt.Claims
		claims.ID = user.ID
		//claims.Subject = user.FirstName + " " + user.Name
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
