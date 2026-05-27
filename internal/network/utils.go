package network

import (
	"social-media-backend/internal/crypto/tokens"

	"github.com/gin-gonic/gin"
)

type refreshTokens struct {
	SessionID    string
	RefreshToken string
	CsrfToken    *string
}

func refreshGetTokensFromCookies(c *gin.Context) (refreshTokens, error) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return refreshTokens{}, tokens.ErrMissingRefreshToken
	}
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		return refreshTokens{}, tokens.ErrMissingSessionID
	}
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" {
		return refreshTokens{}, tokens.ErrMissingCsrf
	}

	return refreshTokens{
		RefreshToken: refreshToken,
		CsrfToken:    &csrfToken,
		SessionID:    sessionID,
	}, nil
}

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
