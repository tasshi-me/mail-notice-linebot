package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"./lineapi"
	"./mongodb"
	"./workers"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Getenv("ENV_LOADED")) < 1 {
		DotEnvLoad()
	}

	//mailCheck()
	//sendVerificationMail("Test User", os.Getenv("IMAP_AUTH_USER"), time.Now().String())
	// Start MailCheckWorker
	go workers.MailCheckWorker()

	port := os.Getenv("PORT")
	http.HandleFunc("/", lineapi.WebhookHandler)
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
