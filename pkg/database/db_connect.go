package database

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Database

func Connect() *mgo.Database {
	const (
		Database   = "shoppingcart"
		UserCollection = "user"
		CartCollection = "cart"
	)
	session, err := mgo.Dial(os.Getenv("MONGO_HOST"))
	 
	if err != nil {
		panic(err)
	} 
	db = session.DB(Database)
	return db
}