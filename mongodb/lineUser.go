package mongodb

import (
	"log"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// LineUser ..
type LineUser struct {
	LineID              string   `bson:"line_id"`
	LineName            string   `bson:"line_name"`
	RegisteredAddresses []string `bson:"registered_address"`
}

// CreateIndexForLineUser ..
func CreateIndexForLineUser(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("LineUser")

	//Create Index
	index := mgo.Index{
		Key:    []string{"line_id"},
		Unique: true,
	}
	err = col.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateOrUpdateLineUser ..
func CreateOrUpdateLineUser(lineUser LineUser) {
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")

	// Insert LineUser
	col := db.C("LineUser")
	if err := col.Insert(&lineUser); err != nil {
		log.Fatalln(err)
	}
}

// ReadAllLineUsers ..
func ReadAllLineUsers() []LineUser {
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")

	// Read All LineUsers
	p := []LineUser{}
	query := db.C("LineUser").Find(bson.M{})
	query.All(&p)

	return p
}
