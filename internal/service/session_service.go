package service

import "social-media-backend/internal/repo"

type SessionService struct {
	repo repo.SessionsRepository
}

func NewSessionService(repo repo.SessionsRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}
