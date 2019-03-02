package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	lineAccessToken := os.Getenv("LINE_ACCESS_TOKEN")
	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	bot, err := linebot.New(lineAccessToken, lineChannelSecret)

	events, err := bot.ParseRequest(r)
	if err != nil {
		// Do something when something bad happened.
		log.Print(err)
		w.WriteHeader(404)
		return
	}

	for _, event := range events {

		var userID string
		var groupID string
		var RoomID string
		log.Print(event.Source.Type)
		switch event.Source.Type {
		case linebot.EventSourceTypeUser:
			userID = event.Source.UserID
		case linebot.EventSourceTypeGroup:
			groupID = event.Source.GroupID
		case linebot.EventSourceTypeRoom:
			RoomID = event.Source.RoomID
		}
		var targetID = userID
		targetID = groupID
		targetID = RoomID

		replyToken := event.ReplyToken

		log.Print(event.Type)
		switch event.Type {
		case linebot.EventTypeMessage:
			// Do Something...
			message := "This is replay. Your message ID = " + targetID
			if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {

			}
		case linebot.EventTypeFollow:
			// Do Something...
		case linebot.EventTypeUnfollow:
			// Do Something...
		case linebot.EventTypeJoin:
			// Do Something...
		case linebot.EventTypeLeave:
			// Do Something...
		case linebot.EventTypePostback:
			// Do Something...
		case linebot.EventTypeBeacon:
			// Do Something...
		default:
			// Do Something...
		}

		// if _, err := bot.PushMessage(targetID, linebot.NewTextMessage("hello")).Do(); err != nil {

		// }
	}

	// leftBtn := linebot.NewMessageAction("left", "left clicked")
	// rightBtn := linebot.NewMessageAction("right", "right clicked")

	// template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)

	// message := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)
	w.WriteHeader(200)
}
