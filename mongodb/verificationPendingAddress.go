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

// CreateOrUpdateVerificationPendingAddress ..
func CreateOrUpdateVerificationPendingAddress(verificationPendingAddress VerificationPendingAddress, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("VerificationPendingAddress")

	if _, err := col.Upsert(bson.M{"line_id": verificationPendingAddress.LineID, "address": verificationPendingAddress.Address}, &verificationPendingAddress); err != nil {
		log.Println(err)
	}
}

// ReadAllVerificationPendingAddress ..
func ReadAllVerificationPendingAddress(url string) []VerificationPendingAddress {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("VerificationPendingAddress")

	// Read All VerificationPendingAddress
	verificationPendingAddresses := []VerificationPendingAddress{}
	query := col.Find(bson.M{})
	query.All(&verificationPendingAddresses)

	return verificationPendingAddresses
}

// ReadVerificationPendingAddress ..
func ReadVerificationPendingAddress(verificationCodeHash string, url string) VerificationPendingAddress {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("VerificationPendingAddress")

	// Find VerificationPendingAddress by VerificationPendingAddress.VerificationCodeHash
	verificationPendingAddresses := VerificationPendingAddress{}
	query := col.Find(bson.M{"verification_code_hash": verificationCodeHash})
	query.One(&verificationPendingAddresses)

	return verificationPendingAddresses
}

// DeleteAllVerificationPendingAddress ..
func DeleteAllVerificationPendingAddress(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("VerificationPendingAddress")

	// Remove All VerificationPendingAddress
	if _, err := col.RemoveAll(bson.M{}); err != nil {
		log.Println(err)
	}
}

// DeleteVerificationPendingAddress ..
func DeleteVerificationPendingAddress(lineID string, url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("mgo.Dial: ", err)
	}
	defer session.Close()

	db := session.DB("")
	col := db.C("VerificationPendingAddress")

	// Remove VerificationPendingAddress by VerificationPendingAddress.LineID
	if _, err := col.RemoveAll(bson.M{"line_id": lineID}); err != nil {
		log.Println(err)
	}
}
