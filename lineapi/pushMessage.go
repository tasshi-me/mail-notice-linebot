package lineapi

import (
	"log"
	"os"
	"strconv"

	"../mailmanager"
	"github.com/line/line-bot-sdk-go/linebot"
)

// SendPushNotification ..
func SendPushNotification(userMailObjects []mailmanager.UserMailObject) {
	//lineChannelID := os.Getenv("LINE_CHANNEL_ID")
	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	lineAccessToken := os.Getenv("LINE_ACCESS_TOKEN")

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
