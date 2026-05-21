package network

import (
	"log"
	"net/http"
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
	Username string `json:"user_name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
		res := mapError(err)
		c.JSON(res.Status, res)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *AuthHandler) VerifyUserEmail(c *gin.Context) {
	token := c.Query("t")
	if token == "" {
		log.Fatalln("invalid email verification token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	err := h.svc.VerifyUserEmail(c.Request.Context(), token)

	if err != nil {
		log.Printf("request failed: %v", err)
		res := mapError(err)
		c.JSON(res.Status, res)
		return
	}
}
