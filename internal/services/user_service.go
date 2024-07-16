package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{Repo: repo}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.Repo.GetAllUsers(ctx)
}

func (s *UserService) SignUp(ctx context.Context, user models.User) error {
	return s.Repo.SignUp(ctx, user)
}
