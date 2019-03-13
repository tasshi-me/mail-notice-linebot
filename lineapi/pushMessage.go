package lineapi

import (
	"log"
	"strconv"

	"../helper"
	"../mailmanager"
	"github.com/line/line-bot-sdk-go/linebot"
)

// SendPushNotification ..
func SendPushNotification(userMailObjects []mailmanager.UserMailObject) {
	configVars := helper.ConfigVars()

	//lineChannelID := configVars.LineAPI.ChannelID
	lineChannelSecret := configVars.LineAPI.ChannelSecret
	lineAccessToken := configVars.LineAPI.AccessToken

	bot, err := linebot.New(lineChannelSecret, lineAccessToken)
	if err != nil {
		log.Print(err)
	}

	for _, userMailObject := range userMailObjects {
		var textContents string
		textContents = "新着メールが" + strconv.Itoa(len(userMailObject.MailObjects)) + "件あります\n"
		for i, mailObject := range userMailObject.MailObjects {
			if len(userMailObject.MailObjects) > 1 {
				textContents += strconv.Itoa(i+1) + ".\n"
			}
			if len(mailObject.MailFromName) > 0 {
				textContents += "差出人: " + mailObject.MailFromName + "\n"
			} else {
				textContents += "差出人: " + mailObject.MailFromAddress + "\n"
			}
			//textContents += "宛先: " + mailObject.MailReceivedAddress + "\n"
			textContents += "件名: " + mailObject.MailSubject + "\n"
		}
		if _, err := bot.PushMessage(userMailObject.TargetLineID, linebot.NewTextMessage(textContents)).Do(); err != nil {
			log.Print(err)
		}
	}

}
