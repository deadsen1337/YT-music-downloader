package utils

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
