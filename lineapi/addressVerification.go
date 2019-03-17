package lineapi

import (
	"crypto/sha256"
	"time"

	"../helper"
	"../mongodb"
	"github.com/globalsign/mgo/bson"
)

// GenerateVerificationCode ..
func GenerateVerificationCode(lineID string, address string) {
	configVars := helper.ConfigVars()
	verificationCode := bson.NewObjectId().String()
	verificationCodeHash := sha256.Sum256([]byte(verificationCode))
	verificationPendingAddress := mongodb.VerificationPendingAddress{
		LineID:               lineID,
		Address:              address,
		VerificationCodeHash: string(verificationCodeHash[:]),
		CreatedAt:            time.Now(),
	}
	mongodb.CreateOrUpdateVerificationPendingAddress(verificationPendingAddress, configVars.MongodbURI)
}

// VerifyAddress ..
func VerifyAddress(lineID string, verificationCode string) {
	configVars := helper.ConfigVars()

	verificationCodeHash := sha256.Sum256([]byte(verificationCode))
	verificationPendingAddress := mongodb.ReadVerificationPendingAddress(string(verificationCodeHash[:]), configVars.MongodbURI)

}
