package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// CheckEnv checks if the proper environment variables are present
func CheckEnv() {
	val := os.Getenv("REDIS_URL")
	if len(val) == 0 {
		log.Fatal("REDIS_URL is not defined, exiting")
	}
}
