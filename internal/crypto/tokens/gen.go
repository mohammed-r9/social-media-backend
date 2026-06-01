package tokens

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"social-media-backend/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func GenerateAccessToken(user domain.User, key []byte) (string, error) {
	isVerified := user.VerifiedAt != nil

	token, err := jwt.NewBuilder().
		Issuer("social-media-app").
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(ACCESS_TOKEN_TTL)).
		Subject(user.ID.String()).
		Claim("is_email_verified", isVerified).
		Build()
	if err != nil {
		return "", err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

// VerifyAccessToken verifies the token and returns the claims if valid
func VerifyAccessToken(token string, key []byte) (AccessTokenClaims, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return AccessTokenClaims{}, err
	}

	if err := jwt.Validate(verifiedToken); err != nil {
		return AccessTokenClaims{}, err
	}

	userIDStr := verifiedToken.Subject()
	if userIDStr == "" {
		return AccessTokenClaims{}, errors.New("subject missing")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	isVerifiedRaw, found := verifiedToken.Get("is_email_verified")
	if !found {
		return AccessTokenClaims{}, ErrMissingClaim
	}

	isVerified, ok := isVerifiedRaw.(bool)
	if !ok {
		return AccessTokenClaims{}, errors.New("invalid type for claim is_email_verified")
	}

	return AccessTokenClaims{
		UserID:          userID,
		IsEmailVerified: isVerified,
	}, nil
}
func GenerateOpaqueToken(length int) string {
	randomByte := make([]byte, length)
	rand.Read(randomByte)
	return hex.EncodeToString(randomByte)
}

func newShortLivedToken(userID uuid.UUID, scope TokenScope, ttl time.Duration) ShortLivedToken {
	raw := GenerateOpaqueToken(32)
	hash := HashToken(raw)

	return ShortLivedToken{
		UserID:    userID,
		Scope:     scope,
		Ttl:       ttl,
		Raw:       raw,
		Hash:      hash,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func GenerateEmailVerificationToken(userID uuid.UUID) ShortLivedToken {
	return newShortLivedToken(userID, ScopeEmailVerification, EMAIL_VERIFICATION_TTL)
}

func GeneratePasswordResetToken(userID uuid.UUID) ShortLivedToken {
	return newShortLivedToken(userID, ScopePasswordReset, PASSWORD_RESET_TTL)
}

func GenerateSessionID() string {
	return GenerateOpaqueToken(24)
}
