package service

import (
	"context"
	"todo_app/internal/models"
)

type Service struct {
	AuthService
	UserService
	TaskService
}

type AuthService interface {
	Register(ctx context.Context, firstName, lastName, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (models.TokenPair, error)
	Authenticate(ctx context.Context, token string) (int64, error)
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, userID int64, input models.CreateTaskInput) (*models.Task, error)
	GetTask(ctx context.Context, userID, taskID int64) (*models.Task, error)
	ListTask(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
	UpdateTask(ctx context.Context, userID, taskID int64, input models.CreateTaskInput) (*models.Task, error)
	SetComplete(ctx context.Context, userID, taskID int64, completed bool) (*models.Task, error)
	DeleteTask(ctx context.Context, userID, taskID int64) error
}
