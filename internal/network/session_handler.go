package network

import "social-media-backend/internal/service"

type Handler struct {
	svc *service.SessionService
}

func NewHandler(svc *service.SessionService) *Handler {
	return &Handler{
		svc: svc,
	}
}
