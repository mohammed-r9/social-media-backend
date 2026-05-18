package service

import (
	"social-media-backend/internal/repo/postgres"
)

type SessionService struct {
	repo postgres.SessionsRepository
}

func NewSessionService(repo postgres.SessionsRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}
