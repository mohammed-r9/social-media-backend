package auth

import (
	"social-media-backend/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func makeUser() domain.User {
	return domain.User{
		ID:         uuid.New(),
		Email:      "mail@mail.com",
		Username:   "username",
		VerifiedAt: nil,
	}
}

type generateJWTParams struct {
	issuer string
	exp    time.Time
	key    []byte
}

func makeJWT(params generateJWTParams, user domain.User) (string, error) {
	iat := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss": params.issuer,
		"iat": iat,
		"exp": params.exp.Unix(),
		"sub": user.ID.String(),

		"is_email_verified": user.VerifiedAt != nil,
		"username":          user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(params.key)
}
