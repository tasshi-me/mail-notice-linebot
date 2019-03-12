package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"./lineapi"
	"./mailmanager"
	"./mongodb"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Getenv("ENV_LOADED")) < 1 {
		DotEnvLoad()
	}

	//mailCheck()
	//sendVerificationMail("Test User", os.Getenv("IMAP_AUTH_USER"), time.Now().String())

	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// DotEnvLoad load .env file
func DotEnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("DotEnv:", err)
	}
}

func mailCheck() {
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
