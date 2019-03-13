package helper

import (
	"log"

	"github.com/joho/godotenv"
)

// DotEnvLoad load .env file
func DotEnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("DotEnv:", err)
	}
}
