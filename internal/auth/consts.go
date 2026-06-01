package auth

import "time"

const (
	// token scopes
	ScopePasswordReset     TokenScope = "password_reset"
	ScopeEmailVerification TokenScope = "email_verification"

	// TTLs
	PASSWORD_RESET_TTL     = time.Minute * 30
	EMAIL_VERIFICATION_TTL = time.Hour * 2
	REFRESH_TTL            = time.Hour * 24 * 30
	ACCESS_TOKEN_TTL       = time.Minute * 15

	ISSUER = "social-media-app"
)
