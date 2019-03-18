package lineapi

import (
	"log"
	"math/rand"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

// SendConfirmSetupForwarding //
func SendConfirmSetupForwarding(bot *linebot.Client, replyToken string) {
	// Send Current registered addres and confirm resetting
	var messages []linebot.SendingMessage

	addresses := []string{"a@a.a", "b@b.b", "c@c.c"}

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
func SendConfirmRevokeForwarding(bot *linebot.Client, replyToken string) {
	// Send Current registered addres and confirm resetting
	var messages []linebot.SendingMessage

	addresses := []string{"a@a.a", "b@b.b", "c@c.c"}

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
