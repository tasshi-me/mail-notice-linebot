package mongodb

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// VerificationPendingAddress ..
type VerificationPendingAddress struct {
	LineID               string    `bson:"line_id"`
	Address              string    `bson:"address"`
	VerificationCodeHash string    `bson:"verification_code_hash"`
	CreatedAt            time.Time `bson:"created_at"`
}

// CreateIndexForVerificationPendingAddress ..
func CreateIndexForVerificationPendingAddress(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("VerificationPendingAddress")

	//Create Index
	indexes := []mgo.Index{
		{
			Key:    []string{"line_id", "address"},
			Unique: true,
		}, {
			Key:    []string{"verification_code_hash"},
			Unique: true,
		},
	}
	for _, index := range indexes {
		err = col.EnsureIndex(index)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CreateOrUpdateConfiguringUser ..
func CreateOrUpdateConfiguringUser(verificationPendingAddress VerificationPendingAddress, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("ConfiguringUser")

	if _, err := col.Upsert(bson.M{"line_id": configuringUser.LineID}, &configuringUser); err != nil {
		log.Println(err)
	}
}

// ReadAllConfiguringUsers ..
func ReadAllConfiguringUsers(url string) []VerificationPendingAddress {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("ConfiguringUser")

	// Read All LineUsers
	lineUser := []VerificationPendingAddress{}
	query := col.Find(bson.M{})
	query.All(&lineUser)

	return lineUser
}

// ReadConfiguringUser ..
func ReadConfiguringUser(lineID string, url string) VerificationPendingAddress {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("ConfiguringUser")

	// Find ConfiguringUser by ConfiguringUser.LineID
	lineUser := VerificationPendingAddress{}
	query := col.Find(bson.M{"line_id": lineID})
	query.One(&lineUser)

	return lineUser
}

// DeleteAllConfiguringUsers ..
func DeleteAllConfiguringUsers(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("ConfiguringUser")

	// Remove All ConfiguringUser
	if _, err := col.RemoveAll(bson.M{}); err != nil {
		log.Println(err)
	}
}

// DeleteConfiguringUser ..
func DeleteConfiguringUser(lineID string, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("ConfiguringUser")

	// Remove ConfiguringUser by ConfiguringUser.LineID
	if _, err := col.RemoveAll(bson.M{"line_id": lineID}); err != nil {
		log.Println(err)
	}
}
