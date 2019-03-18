package lineapi

import (
	"crypto/sha256"
	"errors"
	"net/mail"
	"time"

	"../helper"
	"../mailmanager"
	"../mongodb"
)

// GenerateVerificationCode ..
func GenerateVerificationCode(lineID string, address string) string {
	configVars := helper.ConfigVars()
	verificationCode := "VC-" + helper.GenerateRandomString(64)
	verificationCodeHash := sha256.Sum256([]byte(verificationCode))
	verificationPendingAddress := mongodb.VerificationPendingAddress{
		LineID:               lineID,
		Address:              address,
		VerificationCodeHash: string(verificationCodeHash[:]),
		CreatedAt:            time.Now(),
	}
	mongodb.CreateOrUpdateVerificationPendingAddress(verificationPendingAddress, configVars.MongodbURI)
	return verificationCode
}

// VerifyAddress ..
func VerifyAddress(lineID string, verificationCode string) (string, error) {
	configVars := helper.ConfigVars()

	verificationCodeHash := sha256.Sum256([]byte(verificationCode))
	verificationPendingAddress := mongodb.ReadVerificationPendingAddress(string(verificationCodeHash[:]), configVars.MongodbURI)
	if verificationPendingAddress.LineID != lineID {
		return "", errors.New("無効な確認コードです")
	}
	if verificationPendingAddress.VerificationCodeHash != string(verificationCodeHash[:]) {
		return "", errors.New("無効な確認コードです")
	}
	if time.Now().Sub(verificationPendingAddress.CreatedAt) > time.Minute*5 {
		mongodb.DeleteVerificationPendingAddress(lineID, string(verificationCodeHash[:]), configVars.MongodbURI)
		return "", errors.New("確認コードの有効期限が切れました")
	}

	lineUser := mongodb.ReadLineUser(lineID, configVars.MongodbURI)
	lineUser.RegisteredAddresses = append(lineUser.RegisteredAddresses, verificationPendingAddress.Address)
	mongodb.DeleteVerificationPendingAddress(lineID, string(verificationCodeHash[:]), configVars.MongodbURI)
	return verificationPendingAddress.Address, nil

}

// SendVerificationMail ..
func SendVerificationMail(userName, userAddress, verificationKey string) {
	configVars := helper.ConfigVars()
	from := mail.Address{Name: configVars.SMTP.SenderUsername, Address: configVars.SMTP.SenderAddress}
	to := mail.Address{Name: userName, Address: userAddress}
	subject := "LINEBOT: メールお知らせくん登録確認"
	body := "この度はメールお知らせくんのご利用ありがとうございます。\n LINEの戻って以下の確認コードを送信してください。\n 確認コード：" + verificationKey
	smptServerName := configVars.SMTP.ServerName
	smtpAuthUser := configVars.SMTP.AuthUser
	smtpAuthPassword := configVars.SMTP.AuthPassword
	mailmanager.SendMail(from, to, subject, body, smptServerName, smtpAuthUser, smtpAuthPassword)
}
