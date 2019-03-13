package main

import (
	"log"
	"net/http"
	"os"

	"./lineapi"
	"./mongodb"
	"./workers"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Getenv("ENV_LOADED")) < 1 {
		DotEnvLoad()
	}

	// Init DB
	mongodbURL := os.Getenv("MONGODB_URI")
	mongodb.CreateIndexForLineUser(mongodbURL)

	// Start Keep-Alive Worker for Heroku
	herokuAppName := os.Getenv("HEROKU_APP_NAME")
	if len(herokuAppName) > 0 {
		appURL := "https://" + herokuAppName + ".herokuapp.com/"
		go workers.KeepAliveWorker(appURL)
	}

	// Start MailCheckWorker
	go workers.MailCheckWorker()

	// Start http server for linebot webhook
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
