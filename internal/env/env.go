package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_CONNECTION string
	JWTKey              []byte
	RESEND_API_KEY      string
	BACKEND_URL         string
	FRONTEND_URL        string
	S3_KEY              string
	S3_SECRET           string
	S3_SESSION          string
	S3_URL              string
	REDIS_ADDR          string
	REDIS_PASSWORD      string
}

func New(path string) (*Config, error) {
	_ = godotenv.Load(path)

	get := func(key string) (string, error) {
		v := os.Getenv(key)
		if v == "" {
			return "", fmt.Errorf("missing env var: %s", key)
		}
		return v, nil
	}

	postgres, err := get("POSTGRES_CONNECTION")
	if err != nil {
		return nil, err
	}

	jwt, err := get("JWT_KEY")
	if err != nil {
		return nil, err
	}

	resend, err := get("RESEND_API_KEY")
	if err != nil {
		return nil, err
	}

	backendURL, err := get("BACKEND_URL")
	if err != nil {
		return nil, err
	}

	frontendURL, err := get("FRONTEND_URL")
	if err != nil {
		return nil, err
	}

	s3Key, err := get("S3_KEY")
	if err != nil {
		return nil, err
	}

	s3Secret, err := get("S3_SECRET")
	if err != nil {
		return nil, err
	}

	s3Session, err := get("S3_SESSION")
	if err != nil {
		return nil, err
	}

	s3URL, err := get("S3_URL")
	if err != nil {
		return nil, err
	}

	redisAddr, err := get("REDIS_ADDR")
	if err != nil {
		return nil, err
	}
	redisPass, err := get("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Config{
		POSTGRES_CONNECTION: postgres,
		JWTKey:              []byte(jwt),
		RESEND_API_KEY:      resend,
		BACKEND_URL:         backendURL,
		FRONTEND_URL:        frontendURL,
		S3_KEY:              s3Key,
		S3_SECRET:           s3Secret,
		S3_SESSION:          s3Session,
		S3_URL:              s3URL,
		REDIS_ADDR:          redisAddr,
		REDIS_PASSWORD:      redisPass,
	}, nil
}
