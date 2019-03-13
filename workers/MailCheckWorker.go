package workers

import (
	"log"
	"os"
	"time"

	"../lineapi"
	"../mailmanager"
	"../mongodb"
)

// MailCheck ..
func MailCheck() {
	mboxName := os.Getenv("IMAP_MBOX_NAME")
	dateSince := time.Now().AddDate(0, 0, -1)
	dateBefore := time.Now().AddDate(0, 0, 1)

	messages := mailmanager.FetchMail(dateSince, dateBefore, mboxName, os.Getenv("IMAP_SERVER_NAME"), os.Getenv("IMAP_AUTH_USER"), os.Getenv("IMAP_AUTH_PASSWORD"))
	//messages := mailmanager.PopMail(dateSince, dateBefore, mboxName, os.Getenv("IMAP_SERVER_NAME"), os.Getenv("IMAP_AUTH_USER"), os.Getenv("IMAP_AUTH_PASSWORD"))
	log.Println("fetched messages: ", len(messages))
	// for _, msg := range messages {
	// 	log.Println(msg.Envelope.Date.String() + ":" + msg.Envelope.Subject)
	// }
	if len(messages) > 0 {
		lineUsers := mongodb.ReadAllLineUsers(os.Getenv("MONGODB_URI"))

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
