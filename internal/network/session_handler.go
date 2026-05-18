package network

import "social-media-backend/internal/service"

type SessionHandler struct {
	svc *service.SessionService
}

func NewHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{
		svc: svc,
	}
}
