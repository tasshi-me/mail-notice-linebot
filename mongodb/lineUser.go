package mongodb

import (
	"log"

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
func CreateOrUpdateLineUser(lineUser LineUser, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("LineUser")

	lu := ReadLineUser(lineUser.LineID, url)
	if len(lu.LineID) > 0 {
		// Update LineUser
		if err := col.Update(bson.M{"line_id": lineUser.LineID}, &lineUser); err != nil {
			log.Println(err)
		}
	} else {
		// Insert LineUser
		if err := col.Insert(&lineUser); err != nil {
			log.Println(err)
		}
	}
}

// ReadAllLineUsers ..
func ReadAllLineUsers(url string) []LineUser {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("LineUser")

	// Read All LineUsers
	lineUser := []LineUser{}
	query := col.Find(bson.M{})
	query.All(&lineUser)

	return lineUser
}

// ReadLineUser ..
func ReadLineUser(lineID string, url string) LineUser {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("LineUser")

	// Find LineUser by LineUser.LineID
	lineUser := LineUser{}
	query := col.Find(bson.M{"line_id": lineID})
	query.One(&lineUser)

	return lineUser
}
