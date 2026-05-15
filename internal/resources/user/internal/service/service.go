package service

import repo "social-media-backend/internal/resources/user/internal/repository"

type Service struct {
	repo repo.UserRepository
}

func NewService(repo repo.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}
