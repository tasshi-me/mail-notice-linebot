package lineapi

import (
	"log"
	"net/http"
	"net/mail"
	"strings"

	"../helper"
	"../mailmanager"

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
					SendConfirmSetupForwarding(bot, replyToken)
				case strings.Contains(message.Text, "お知らせ解除"):
					SendConfirmRevokeForwarding(bot, replyToken)
				case strings.HasPrefix(message.Text, "vcode"):
					VerifyAddress(targetID, message.Text)
				}

			}
			if eventSourceType == linebot.EventSourceTypeUser {
				SendRandomReply(bot, replyToken)
			}

		case linebot.EventTypeFollow:
			// Send Introduction to user
			SendIntroduction(bot, replyToken)
		case linebot.EventTypeUnfollow:
			// TODO: Delete User from database
		case linebot.EventTypeJoin:
			// Send Introduction to the group
			SendIntroduction(bot, replyToken)
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
	}
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
