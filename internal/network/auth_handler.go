package network

import (
	"log"
	"net/http"
	"social-media-backend/internal/crypto/tokens/stateful"
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
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
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

	c.SetCookie(
		"refresh_token",
		sessionTokens.RefreshToken,
		int(stateful.REFRESH_TTL),
		"/",
		"",
		true,
		true,
	)
	c.SetCookie("csrf_token",
		sessionTokens.CsrfToken,
		int(stateful.REFRESH_TTL),
		"/",
		"",
		true,
		false,
	)
	c.SetCookie("session_id",
		sessionTokens.SessionID,
		int(stateful.REFRESH_TTL),
		"/",
		"",
		true,
		true,
	)

	OK(c, gin.H{
		"access_token": sessionTokens.AccessToken,
	})
}
