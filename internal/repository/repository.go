package repository

import (
	"context"
	"todo_app/internal/models"
)

type Repositories struct {
	User UserRepository
	Task TaskRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Task) (*models.Task, error)
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
	Update(ctx context.Context, task models.Task) (*models.Task, error)
	Delete(ctx context.Context, taskID int64) error
	SetCompleted(ctx context.Context, taskID int64, complete bool) (*models.Task, error)
}
