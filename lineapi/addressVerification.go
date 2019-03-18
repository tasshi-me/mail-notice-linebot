package lineapi

import (
	"crypto/sha256"
	"errors"
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
func VerifyAddress(lineID string, verificationCode string) error {
	configVars := helper.ConfigVars()

	verificationCodeHash := sha256.Sum256([]byte(verificationCode))
	verificationPendingAddress := mongodb.ReadVerificationPendingAddress(string(verificationCodeHash[:]), configVars.MongodbURI)
	if time.Now().Sub(verificationPendingAddress.CreatedAt) > time.Minute*5 {
		return errors.New("確認コードの有効期限が切れました")
	}
	if verificationPendingAddress.LineID != lineID {
		return errors.New("無効な確認コードです")
	}
	if verificationPendingAddress.VerificationCodeHash != string(verificationCodeHash[:]) {
		return errors.New("無効な確認コードです")
	}

	lineUser := mongodb.ReadLineUser(lineID, configVars.MongodbURI)
	lineUser.RegisteredAddresses = append(lineUser.RegisteredAddresses, verificationPendingAddress.Address)
	return nil

}
