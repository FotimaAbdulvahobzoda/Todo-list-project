package task

import (
	"todo_app/internal/repository"
	"todo_app/internal/service"
)

var _ service.TaskService = (*taskService)(nil)

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) service.TaskService {
	return &taskService{taskRepo: taskRepo}
}
