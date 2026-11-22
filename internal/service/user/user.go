package user

import (
	"context"
	"fmt"
	"todo_app/internal/models"
	"todo_app/internal/repository"
	"todo_app/internal/service"
)

var _ service.UserService = (*userService)(nil)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) service.UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}
