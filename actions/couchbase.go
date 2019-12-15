package actions

import (
	"gopkg.in/couchbase/gocb.v1"
	"os"
	"errors"
)

type User struct {
	ID string 		`json:"id"`
	FirstName string	`json:"firstName"`
	Name string		`json:"name"`
	Username string		`json:"username"`
	Email string		`json:"email"`
	Hash string		`json:"hash"`
	Devices []string	`json:"devices"`
	Roles []string		`json:"roles"`
}

func connDB() (*gocb.Cluster, error) {
	host := os.Getenv("NSDB_HOST")
	//port := os.Getenv("NSDB_PORT")
	cluster, err := gocb.Connect("couchbase://" + host)
	if err != nil {
		return nil, errors.New("Failed to connect to Cochbase!")
	}
	/*cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: os.Getenv("NSDB_USR"),
		Password: os.Getenv("NSDB_PW"),
	})*/
	return cluster, nil;
}

func getBucket(cluster *gocb.Cluster, bucketName string) (*gocb.Bucket, error) {
	bucket, err := cluster.OpenBucket(bucketName, os.Getenv("NSDB_BPW"))
	if err != nil {
		return nil, errors.New("Bucket error: " + err.Error())
	}
	return bucket, nil
}
