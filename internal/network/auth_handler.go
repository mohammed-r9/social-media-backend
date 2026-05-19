package network

import (
	"errors"
	"log"
	"net/http"
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
	})
	// verification email is yet to be added here

	if err != nil {
		// maybe add more errors later?
		switch {
		case errors.Is(err, domain.ErrEmailAlreadyTaken):
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already exists",
			})
			return

		default:
			log.Printf("register user failed email=%s err=%v", req.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, createdUser)
}
