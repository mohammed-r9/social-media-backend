package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Config = &envConfig{}

func LoadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Printf("No .env file found: %v\n", err)
	}

	Config = &envConfig{
		POSTGRES_CONNECTION: mustGetEnv("POSTGRES_CONNECTION"),
		RESEND_API_KEY:      mustGetEnv("RESEND_API_KEY"),
		BACKEND_URL:         mustGetEnv("BACKEND_URL"),
		FRONTEND_URL:        mustGetEnv("FRONTEND_URL"),
		S3_KEY:              mustGetEnv("S3_KEY"),
		S3_SECRET:           mustGetEnv("S3_SECRET"),
		S3_SESSION:          mustGetEnv("S3_SESSION"),
		S3_URL:              mustGetEnv("S3_URL"),

		JWT_KEY: []byte(mustGetEnv("JWT_KEY")),
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("error reading environment variable: %v\n", key)
	}
	return value
}

type envConfig struct {
	POSTGRES_CONNECTION string
	// App
	JWT_KEY        []byte
	RESEND_API_KEY string
	BACKEND_URL    string
	FRONTEND_URL   string
	S3_KEY         string
	S3_SECRET      string
	S3_SESSION     string
	S3_URL         string
}
