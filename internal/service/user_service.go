package service

import (
	"social-media-backend/internal/repo/postgres"
)

type UserService struct {
	repo postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type RegisterParams struct {
	Name              string
	Email             string
	Username          string
	PassowrdPlainText string
}
