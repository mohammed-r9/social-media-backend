package network

import (
	"errors"
	"log"
	"net/http"
	"social-media-backend/internal/crypto/tokens"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// need to update it later
		log.Printf("error binding request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	createdUser, err := h.svc.Register(c.Request.Context(), service.RegisterParams{
		Name:              req.Name,
		Email:             req.Email,
		PassowrdPlainText: req.Password,
		Username:          req.Username,
	})
	// verification email is yet to be added here

	if err != nil {
		log.Printf("request failed: %v", err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *AuthHandler) VerifyUserEmail(c *gin.Context) {
	token := c.Query("t")
	if token == "" {
		// need to update it later
		log.Println("invalid email verification token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	err := h.svc.VerifyUserEmail(c.Request.Context(), token)

	if err != nil {
		log.Printf("request failed: %v", err)
		_ = c.Error(err)
		return
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	authMode := c.GetHeader("X-Auth-Mode") // SHOULD ONLY BE USED TO FORMAT THE RESPONSE
	if authMode == "" {
		_ = c.Error(errMissingAuthModeHeader)
		return
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	sessionTokens, err := h.svc.Login(c.Request.Context(), service.LoginParams{
		Email:      req.Email,
		Password:   req.Password,
		DeviceName: "temp",
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	switch authMode {
	case "cookie":
		c.SetCookie(
			"refresh_token",
			sessionTokens.RefreshToken,
			int(tokens.REFRESH_TTL),
			"/",
			"",
			true,
			true,
		)
		c.SetCookie("csrf_token",
			sessionTokens.CsrfToken,
			int(tokens.REFRESH_TTL),
			"/",
			"",
			true,
			false,
		)
		c.SetCookie("session_id",
			sessionTokens.SessionID,
			int(tokens.REFRESH_TTL),
			"/",
			"",
			true,
			true,
		)

		OK(c, gin.H{
			"access_token": sessionTokens.AccessToken,
		})

	case "token":
		OK(c, gin.H{
			"access_token":  sessionTokens.AccessToken,
			"refresh_token": sessionTokens.RefreshToken,
			"session_id":    sessionTokens.SessionID,
		})
	default:
		_ = c.Error(errInvalidAuthModeHeader)
		return
	}
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	var opTokens refreshTokens

	// web
	opTokens, err := refreshGetTokensFromCookies(c)
	if err != nil {
		// mobile
		opTokens, err = refreshGetTokensFromBody(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}

	accessToken, err := h.svc.RefreshAccessToken(c.Request.Context(), service.RefreshParams{
		RefreshToken: opTokens.RefreshToken,
		SessionID:    opTokens.SessionID,
		CsrfToken:    opTokens.CsrfToken,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, gin.H{
		"access_token": accessToken,
	})
}

type forgotPassowrdRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPassowrdRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.svc.AskForPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			_ = c.Error(err)
			return
		}
	}

	OK(c, gin.H{
		"message": "If the account exists, a password reset email has been sent.",
	})
}

type resetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.svc.ResetUserPassword(c.Request.Context(), service.ResetUserPasswordParams{
		Token:       req.Token,
		NewPassword: req.Password,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, gin.H{
		"message": "Password resetted successfully",
	})
}
