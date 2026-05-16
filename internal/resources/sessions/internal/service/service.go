package service

import "social-media-backend/internal/resources/sessions/internal/repository"

type Service struct {
	repo repository.SessionsRepository
}

func NewService(repo repository.SessionsRepository) *Service {
	return &Service{
		repo: repo,
	}
}
