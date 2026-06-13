package auth

import (
	"time"

	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(user domain.User, key []byte) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"iss": JWT_ISSUER,
		"iat": now.Unix(),
		"exp": now.Add(ACCESS_TOKEN_TTL).Unix(),
		"sub": user.ID.String(),

		"is_email_verified": user.VerifiedAt != nil,
		"username":          user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}

// VerifyAccessToken verifies the token and returns the claims if valid
func VerifyAccessToken(tokenString string, key []byte) (AccessTokenClaims, error) {
	var claims AccessTokenClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Header["alg"] != "HS256" {
				return nil, apperrors.InvalidToken
			}
			return key, nil
		},
	)

	if err != nil {
		return AccessTokenClaims{}, err
	}

	if !token.Valid {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	if claims.Issuer != JWT_ISSUER {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	if claims.Subject == "" {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return AccessTokenClaims{}, apperrors.ExpiredToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	claims.UserID = userID

	return claims, nil
}
