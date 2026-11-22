package task

import (
	"todo_app/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.TaskRepository = (*taskRepo)(nil)

type taskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *taskRepo {
	return &taskRepo{pool: pool}
}
