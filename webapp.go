package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	//lineChannelID := os.Getenv("LINE_CHANNEL_ID")
	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	lineAccessToken := os.Getenv("LINE_ACCESS_TOKEN")

	bot, err := linebot.New(lineChannelSecret, lineAccessToken)

	events, err := bot.ParseRequest(r)
	if err != nil {
		// Do something when something bad happened.
		log.Print("ParseRequest: ", err)
		w.WriteHeader(400)
		return
	}

	for _, event := range events {

		// var userID string
		// var groupID string
		// var RoomID string
		var targetID string

		log.Print("EventSource Type: ", event.Source.Type)
		switch event.Source.Type {
		case linebot.EventSourceTypeUser:
			//userID = event.Source.UserID
			targetID = event.Source.UserID
		case linebot.EventSourceTypeGroup:
			//groupID = event.Source.GroupID
			targetID = event.Source.GroupID
		case linebot.EventSourceTypeRoom:
			//RoomID = event.Source.RoomID
			targetID = event.Source.RoomID
		}

		eventSourceType := event.Source.Type
		replyToken := event.ReplyToken

		log.Print("Event Type: ", event.Type)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.Contains(message.Text, "メールお知らせ") || strings.Contains(message.Text, "メールおしらせ") {
					sendConfirmSetupForwardingMessage(bot, replyToken)
				}
			}
			if eventSourceType == linebot.EventSourceTypeUser {
				sendRandomReply(bot, replyToken)
			}

		case linebot.EventTypeFollow:
			// Send Introduction to user
			message := "This is replay. Your message ID = " + targetID
			if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
				log.Print(err)
			}
		case linebot.EventTypeUnfollow:
			// TODO: Delete User from database
		case linebot.EventTypeJoin:
			// TODO: Send Introduction to the group

		case linebot.EventTypeLeave:
			// TODO: Delete group from database
		case linebot.EventTypeMemberJoined:
			// Send message to Joined User
			// Default send nothing
		case linebot.EventTypeMemberLeft:
			// Send message to Left User
			// Default send nothing
		case linebot.EventTypePostback:
			// Do Nothing
		case linebot.EventTypeBeacon:
			// Do Nothing
		default:
			// Do Nothing
		}

		// if _, err := bot.PushMessage(targetID, linebot.NewTextMessage("hello")).Do(); err != nil {

		// }
	}
}

func sendConfirmSetupForwardingMessage(bot *linebot.Client, replyToken string) {
	// Send Current registered addres and confirm resetting
	var messages []linebot.SendingMessage

	addresses := []string{"a@a.a", "b@b.b", "c@c.c"}

	// Current e-mail addresses
	var textContents = "こんにちは！メールお知らせくんです。\n"
	if len(addresses) > 0 {
		textContents = textContents + "現在転送設定されているメールアドレスは\n" + strings.Join(addresses, "\n") + "\nです"
	} else {
		textContents = textContents + "現在転送設定されているメールアドレスはありません"
	}
	messages = append(messages, linebot.NewTextMessage(textContents))

	// Confirm template message
	var altText string
	if len(addresses) > 0 {
		altText = "メール転送を再設定しますか？"
	} else {
		altText = "メール転送を設定しますか？"
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

func sendRandomReply(bot *linebot.Client, replyToken string) {
	contentPatterns := []string{
		"ごめんなさい！よく分かりませんでした！",
		"「メールお知らせくん」と呼んでいただければメール転送設定が確認できます",
		"新しいメールはありません(たぶん)",
	}
	// Randomize reply
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(contentPatterns))
	message := linebot.NewTextMessage(contentPatterns[i])
	// Send messages
	if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
		log.Print(err)
	}
}
