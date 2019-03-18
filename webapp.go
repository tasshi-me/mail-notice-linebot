package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"./helper"
	"./lineapi"
	"./mongodb"
	"./workers"
)

func main() {
	configVars := helper.ConfigVars()
	if len(configVars.EnvLoaded) < 1 {
		helper.DotEnvLoad()
	}

	// Init rand
	rand.Seed(time.Now().UnixNano())

	// Init DB
	mongodbURL := configVars.MongodbURI
	mongodb.CreateIndexForLineUser(mongodbURL)
	mongodb.CreateIndexForOnConfigureUser(mongodbURL)
	mongodb.CreateIndexForVerificationPendingAddress(mongodbURL)

	// Start Keep-Alive Worker for Heroku
	herokuAppName := configVars.HerokuAppName
	if len(herokuAppName) > 0 {
		interval := 20 * time.Minute
		appURL := "https://" + herokuAppName + ".herokuapp.com/"
		go workers.KeepAliveWorker(interval, appURL)
		log.Println("Start Keep-Alive Worker for Heroku")
	}

	// Start MailCheckWorker
	interval := 5 * time.Minute
	go workers.MailCheckWorker(interval)
	log.Println("Start MailCheck Worker")

	// Start http server for linebot webhook
	port := configVars.Port
	http.HandleFunc("/", lineapi.WebhookHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
