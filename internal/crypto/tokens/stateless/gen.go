package stateless

import (
	"social-media-backend/internal/domain"
	"social-media-backend/internal/env"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwt"
)

func GenerateAccessToken(user domain.User) (string, error) {
	isVerified := user.VerifiedAt != nil

	token, err := jwt.NewBuilder().
		Issuer("social-media-app").
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(ACCESS_TOKEN_TTL)).
		Subject(user.ID.String()).
		Claim("is_email_verified", isVerified).Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), env.Config.JWT_KEY))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

// VerifyAccessToken verifies the token and returns the claims if valid
func VerifyAccessToken(token string) (AccessTokenClaims, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256(), env.Config.JWT_KEY))
	if err != nil {
		return AccessTokenClaims{}, err
	}
	if err := jwt.Validate(verifiedToken); err != nil {
		return AccessTokenClaims{}, err
	}

	userIDStr, found := verifiedToken.Subject()
	if !found {
		return AccessTokenClaims{}, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	isVerified, found := verifiedToken.Field("is_email_verified")
	if !found {
		return AccessTokenClaims{}, ErrMissingClaim
	}

	return AccessTokenClaims{
		UserID:          userID,
		IsEmailVerified: isVerified.(bool),
	}, nil
}
