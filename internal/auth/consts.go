package auth

import (
	"os"
	"time"
)

const (
	PASSWORD_RESET_TTL     = time.Minute * 30
	EMAIL_VERIFICATION_TTL = time.Hour * 24
	REFRESH_TTL            = time.Hour * 24 * 30
)

var accessTokenTTL = func() time.Duration {
	if os.Getenv("ENV") == "dev" {
		return time.Hour * 24 * 365
	}
	return time.Minute * 20
}()

// returns the access token ttl dynamically according to the ENV environment variable
func AccessTokenTTL() time.Duration {
	return accessTokenTTL
}

const (
	SCOPE_PASSWORD_RESET     TokenScope = "password_reset"
	SCOPE_EMAIL_VERIFICATION TokenScope = "email_verification"
)

const JWT_ISSUER = "social-media-app"
