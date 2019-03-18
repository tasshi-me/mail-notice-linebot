package mongodb

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// OnConfigureUser ..
type OnConfigureUser struct {
	LineID    string    `bson:"line_id"`
	Addresses []string  `bson:"address"`
	CreatedAt time.Time `bson:"created_at"`
}

// CreateIndexForOnConfigureUser ..
func CreateIndexForOnConfigureUser(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("OnConfigureUser")

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

// CreateOrUpdateOnConfigureUser ..
func CreateOrUpdateOnConfigureUser(onConfigureUser OnConfigureUser, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("OnConfigureUser")

	if _, err := col.Upsert(bson.M{"line_id": onConfigureUser.LineID}, &onConfigureUser); err != nil {
		log.Println(err)
	}
}

// ReadAllOnConfigureUser ..
func ReadAllOnConfigureUser(url string) []OnConfigureUser {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("OnConfigureUser")

	// Read All ConfigureUsers
	onConfigureUsers := []OnConfigureUser{}
	query := col.Find(bson.M{})
	query.All(&onConfigureUsers)

	return onConfigureUsers
}

// ReadOnConfigureUser ..
func ReadOnConfigureUser(lineID string, url string) OnConfigureUser {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("OnConfigureUser")

	// Find OnConfigureUser by LineUser.LineID
	onConfigureUser := OnConfigureUser{}
	query := col.Find(bson.M{"line_id": lineID})
	query.One(&onConfigureUser)

	return onConfigureUser
}

// DeleteAllOnConfigureUser ..
func DeleteAllOnConfigureUser(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("OnConfigureUser")

	// Remove All OnConfigureUsers
	if _, err := col.RemoveAll(bson.M{}); err != nil {
		log.Println(err)
	}
}

// DeleteOnConfigureUser ..
func DeleteOnConfigureUser(lineID string, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("OnConfigureUser")

	// Remove LineUser by LineUser.LineID
	if _, err := col.RemoveAll(bson.M{"line_id": lineID}); err != nil {
		log.Println(err)
	}
}
