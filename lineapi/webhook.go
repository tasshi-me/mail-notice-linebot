package lineapi

import (
	"log"
	"net/http"
	"strings"

	"github.com/tasshi-me/mail-notice-linebot/helper"

	"github.com/line/line-bot-sdk-go/linebot"
)

// WebhookHandler ..
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	configVars := helper.ConfigVars()

	//lineChannelID := configVars.LineAPI.ChannelID
	lineChannelSecret := configVars.LineAPI.ChannelSecret
	lineAccessToken := configVars.LineAPI.AccessToken

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
		log.Print("TargetID: ", targetID)

		eventSourceType := event.Source.Type
		replyToken := event.ReplyToken

		log.Print("Event Type: ", event.Type)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch {
				case strings.Contains(message.Text, "メールお知らせ"):
					fallthrough
				case strings.Contains(message.Text, "メールおしらせ"):
					SendConfirmSetupForwarding(bot, replyToken, targetID)
				case strings.Contains(message.Text, "お知らせ解除"):
					SendConfirmRevokeForwarding(bot, replyToken, targetID)
				case strings.HasPrefix(message.Text, "VC-"):
					address, err := VerifyAddress(targetID, message.Text)
					var contentText string
					if err != nil {
						contentText = err.Error()
					} else {
						contentText = "メールアドレスが確認されました\n" + address + "\n以下のメールアドレス宛にメール転送設定を行うとお知らせが来るようになります\n" + configVars.IMAP.Address
					}
					message := linebot.NewTextMessage(contentText)
					if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
						log.Print(err)
					}
				case strings.Contains(message.Text, "@"):
					PushAddressToConfigureQueue(bot, replyToken, targetID, message.Text)
				case message.Text == ".":
					FinishConfigureAddress(bot, replyToken, targetID)
				default:
					if eventSourceType == linebot.EventSourceTypeUser {
						SendRandomReply(bot, replyToken)
					}
				}
			}
		case linebot.EventTypeFollow:
			// Send Introduction to user
			SendIntroduction(bot, replyToken)
		case linebot.EventTypeUnfollow:
			RevokeRegisteredUser(bot, replyToken, targetID)
		case linebot.EventTypeJoin:
			// Send Introduction to the group
			SendIntroduction(bot, replyToken)
		case linebot.EventTypeLeave:
			RevokeRegisteredUser(bot, replyToken, targetID)
		case linebot.EventTypeMemberJoined:
			// Send message to Joined User
			// Default send nothing
		case linebot.EventTypeMemberLeft:
			// Send message to Left User
			// Default send nothing
		case linebot.EventTypePostback:
			data := event.Postback.Data
			if data == "setup=true" {
				StartConfigureAddress(bot, replyToken, targetID)
			}
			if data == "revoke=true" {
				RevokeRegisteredUser(bot, replyToken, targetID)
			}
			// Do Nothing
		case linebot.EventTypeBeacon:
			// Do Nothing
		default:
			// Do Nothing
		}
	}
}
