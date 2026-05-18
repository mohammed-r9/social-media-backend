package service

import (
	"social-media-backend/internal/domain"
)

type SessionService struct {
	repo domain.SessionsRepository
}

func NewSessionService(repo domain.SessionsRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}
