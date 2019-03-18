package workers

import (
	"log"
	"time"

	"../helper"
	"../lineapi"
	"../mailmanager"
	"../mongodb"
)

// MailCheck ..
func MailCheck() {
	configVars := helper.ConfigVars()
	mboxName := configVars.IMAP.MboxName
	dateSince := time.Now().AddDate(0, 0, -2)
	dateBefore := time.Now().AddDate(0, 0, 2)

	//messages := mailmanager.FetchMail(dateSince, dateBefore, mboxName, configVars.IMAP.ServerName, configVars.IMAP.AuthUser, configVars.IMAP.AuthPassword)
	messages := mailmanager.PopMail(dateSince, dateBefore, mboxName, configVars.IMAP.ServerName, configVars.IMAP.AuthUser, configVars.IMAP.AuthPassword)
	log.Println("fetched messages: ", len(messages))
	// for _, msg := range messages {
	// 	log.Println(msg.Envelope.Date.String() + ":" + msg.Envelope.Subject)
	// }
	if len(messages) > 0 {
		lineUsers := mongodb.ReadAllLineUsers(configVars.MongodbURI)

		userMailObjects := mailmanager.ConvertMessagesToUserMailObject(messages, lineUsers)

		if len(userMailObjects) > 0 {
			lineapi.SendPushNotification(userMailObjects)
		}
	}
}

// MailCheckWorker ..
func MailCheckWorker(interval time.Duration) {
	tic := time.NewTicker(interval)
	for {
		select {
		case <-tic.C:
			MailCheck()
		}
	}
}
