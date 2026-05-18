package network

import "social-media-backend/internal/service"

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}
