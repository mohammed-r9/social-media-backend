package network

import (
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

const ctxClaimKey = "jwt_claims"

type refreshTokens struct {
	SessionID    string
	RefreshToken string
	CsrfToken    *string
}

func refreshGetTokensFromCookies(c *gin.Context) (refreshTokens, error) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return refreshTokens{}, apperrors.MissingRefreshToken
	}
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		return refreshTokens{}, apperrors.MissingSessionID
	}
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" {
		return refreshTokens{}, apperrors.MissingCsrf
	}

	return refreshTokens{
		RefreshToken: refreshToken,
		CsrfToken:    &csrfToken,
		SessionID:    sessionID,
	}, nil
}

// refreshTokenReq represents the refresh token request body.
type refreshTokenReq struct {
	SessionID    string `json:"session_id" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func refreshGetTokensFromBody(c *gin.Context) (refreshTokens, error) {
	var req refreshTokenReq

	if err := c.ShouldBindJSON(&req); err != nil {
		return refreshTokens{}, err
	}

	return refreshTokens{
		SessionID:    req.SessionID,
		RefreshToken: req.RefreshToken,
		CsrfToken:    nil,
	}, nil
}

// get the calims from the gin context
func getClaims(c *gin.Context) (auth.AccessTokenClaims, bool) {
	v, ok := c.Get(ctxClaimKey)
	if !ok {
		return auth.AccessTokenClaims{}, false
	}

	claims, ok := v.(auth.AccessTokenClaims)
	return claims, ok
}
