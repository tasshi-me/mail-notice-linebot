package lineapi

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tasshi-me/mail-notice-linebot/helper"
	"github.com/tasshi-me/mail-notice-linebot/mongodb"
)

// SendConfirmSetupForwarding //
func SendConfirmSetupForwarding(bot *linebot.Client, replyToken string, lineID string) {
	// Send Current registered addres and confirm resetting
	var messages []linebot.SendingMessage

	configVars := helper.ConfigVars()
	lineUser := mongodb.ReadLineUser(lineID, configVars.MongodbURI)
	addresses := lineUser.RegisteredAddresses

	// Current e-mail addresses
	var textContents = "こんにちは！メールお知らせくんです。\n"
	if len(addresses) > 0 {
		textContents = textContents + "現在お知らせ設定されているメールアドレスは\n" + strings.Join(addresses, "\n") + "\nです"
	} else {
		textContents = textContents + "現在お知らせ設定されているメールアドレスはありません"
	}
	messages = append(messages, linebot.NewTextMessage(textContents))

	// Confirm template message
	var altText string
	if len(addresses) > 0 {
		altText = "メールお知らせを再設定しますか？"
	} else {
		altText = "メールお知らせを設定しますか？"
	}
	leftBtn := linebot.NewPostbackAction("はい", "setup=true", "", "はい")
	rightBtn := linebot.NewPostbackAction("いいえ", "setup=false", "", "いいえ")
	template := linebot.NewConfirmTemplate(altText, leftBtn, rightBtn)
	messages = append(messages, linebot.NewTemplateMessage(altText, template))

	// Send messages
	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
		log.Print(err)
	}
}

// SendConfirmRevokeForwarding ..
func SendConfirmRevokeForwarding(bot *linebot.Client, replyToken string, lineID string) {
	// Send Current registered addres and confirm resetting
	var messages []linebot.SendingMessage

	configVars := helper.ConfigVars()
	lineUser := mongodb.ReadLineUser(lineID, configVars.MongodbURI)
	addresses := lineUser.RegisteredAddresses

	// Current e-mail addresses
	var textContents = "こんにちは！メールお知らせくんです。\n"
	if len(addresses) > 0 {
		textContents = textContents + "現在お知らせ設定されているメールアドレスは\n" + strings.Join(addresses, "\n") + "\nです"
	} else {
		textContents = textContents + "現在お知らせ設定されているメールアドレスはありません"
	}
	messages = append(messages, linebot.NewTextMessage(textContents))

	if len(addresses) <= 0 {
		return
	}

	// Confirm template message
	altText := "メールお知らせを解除しますか？"
	leftBtn := linebot.NewPostbackAction("はい", "revoke=true", "", "はい")
	rightBtn := linebot.NewPostbackAction("いいえ", "revoke=false", "", "いいえ")
	template := linebot.NewConfirmTemplate(altText, leftBtn, rightBtn)
	messages = append(messages, linebot.NewTemplateMessage(altText, template))

	// Send messages
	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
		log.Print(err)
	}
}

// SendRandomReply ..
func SendRandomReply(bot *linebot.Client, replyToken string) {
	contentPatterns := []string{
		"ごめんなさい！よく分かりませんでした！",
		"「メールお知らせくん」と呼んでいただければメールお知らせ設定が確認できます",
		"「お知らせ解除」と言っていただければメールお知らせを解除できます",
		"新しいメールはたぶんありません！",
	}
	// Randomize reply
	i := rand.Intn(len(contentPatterns))
	message := linebot.NewTextMessage(contentPatterns[i])
	// Send messages
	if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
		log.Print(err)
	}
}

// SendIntroduction ..
func SendIntroduction(bot *linebot.Client, replyToken string) {
	// Send Greeting and introduction
	var messages []linebot.SendingMessage

	// Greeting
	var textContents = "登録ありがとうございます！メールお知らせくんです。\n"
	textContents += "登録されたメールアドレスにメールが届くとお知らせします。\n"
	messages = append(messages, linebot.NewTextMessage(textContents))

	// Confirm template message
	altText := "メールお知らせを設定しますか？"
	leftBtn := linebot.NewPostbackAction("はい", "setup=true", "", "はい")
	rightBtn := linebot.NewPostbackAction("いいえ", "setup=false", "", "いいえ")
	template := linebot.NewConfirmTemplate(altText, leftBtn, rightBtn)
	messages = append(messages, linebot.NewTemplateMessage(altText, template))

	// Send messages
	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
		log.Print(err)
	}
}

// RevokeRegisteredUser ..
func RevokeRegisteredUser(bot *linebot.Client, replyToken string, lineID string) {
	configVars := helper.ConfigVars()
	mongodb.DeleteLineUser(lineID, configVars.MongodbURI)

	if len(replyToken) > 0 {
		contentText := "お知らせ設定を削除しました！"
		message := linebot.NewTextMessage(contentText)
		// Send messages
		if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
			log.Print(err)
		}
	}
}

// StartConfigureAddress ..
func StartConfigureAddress(bot *linebot.Client, replyToken string, lineID string) {
	configVars := helper.ConfigVars()
	ocUser := mongodb.ReadOnConfigureUser(lineID, configVars.MongodbURI)
	if ocUser.LineID == lineID {
		contentText := "すでに設定中です\n終了するには「.」を入力してください"
		message := linebot.NewTextMessage(contentText)
		// Send messages
		if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
			log.Print(err)
		}
		return
	}
	ocUser = mongodb.OnConfigureUser{
		LineID:    lineID,
		CreatedAt: time.Now(),
	}
	mongodb.CreateOrUpdateOnConfigureUser(ocUser, configVars.MongodbURI)
	contentText := "メールアドレスを１件ずつ入力してください\n終了するには「.」を入力してください"
	message := linebot.NewTextMessage(contentText)
	// Send messages
	if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
		log.Print(err)
	}
}

// PushAddressToConfigureQueue ..
func PushAddressToConfigureQueue(bot *linebot.Client, replyToken string, lineID string, address string) {
	configVars := helper.ConfigVars()
	ocUser := mongodb.ReadOnConfigureUser(lineID, configVars.MongodbURI)
	if ocUser.LineID == lineID {
		ocUser.Addresses = append(ocUser.Addresses, address)
		mongodb.CreateOrUpdateOnConfigureUser(ocUser, configVars.MongodbURI)
	}
}

// FinishConfigureAddress ..
func FinishConfigureAddress(bot *linebot.Client, replyToken string, lineID string) {
	configVars := helper.ConfigVars()
	ocUser := mongodb.ReadOnConfigureUser(lineID, configVars.MongodbURI)
	if ocUser.LineID != lineID {
		return
	}

	var contentText string
	if len(ocUser.Addresses) > 0 {
		contentText = "以下の" + strconv.Itoa(len(ocUser.Addresses)) + "個のメールアドレスに確認コードをお送りしました。メールを確認して確認コードを入力してください\n"
		for _, address := range ocUser.Addresses {
			verificationCode := GenerateVerificationCode(ocUser.LineID, address)
			SendVerificationMail("", address, verificationCode)
			contentText += address + "\n"
		}
	} else {
		contentText = "メールアドレスが設定されませんでした"
	}
	mongodb.DeleteOnConfigureUser(ocUser.LineID, configVars.MongodbURI)
	message := linebot.NewTextMessage(contentText)
	// Send messages
	if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
		log.Print(err)
	}
}
