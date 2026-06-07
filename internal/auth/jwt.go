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
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
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

	iss, _ := claims["iss"].(string)
	if iss != JWT_ISSUER {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return AccessTokenClaims{}, apperrors.InvalidToken
	}

	isVerified, _ := claims["is_email_verified"].(bool)

	username, _ := claims["username"].(string)

	return AccessTokenClaims{
		UserID:          userID,
		Username:        username,
		IsEmailVerified: isVerified,
	}, nil
}
