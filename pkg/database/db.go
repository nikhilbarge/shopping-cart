package database

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

var Db *mgo.Database

func Connect() *mgo.Database {

	// Host := []string{
	// 	"localhost:27017",
	// 	// replica set addrs...
	// }
	const (
		Database   = "shoppingcart"
		UserCollection = "user"
		CartCollection = "cart"
	)
	session, err := mgo.Dial(os.Getenv("MONGO_HOST"))
	// session, err := mgo.DialWithInfo(&mgo.DialInfo{
	// 	Addrs: Host,
	// 	// Username: Username,
	// 	// Password: Password,
	// 	// Database: Database,
	// 	// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
	// 	// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
	// 	// },
	// })
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	Db = session.DB(Database)
	return Db
}