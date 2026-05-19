package network

import "social-media-backend/internal/service"

type SessionHandler struct {
	svc *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{
		svc: svc,
	}
}
