package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"todo_app/internal/models"
	"unicode/utf8"
)

const (
	maxTitleLength       = 100
	maxDescriptionLength = 1000
)

func (s *taskService) CreateTask(ctx context.Context, userID int64, input models.CreateTaskInput) (*models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	task, err := buildTask(input)
	if err != nil {
		return nil, err
	}

	task.OwnerID = userID

	created, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, models.ErrAssigneeNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return created, err
}

func (s *taskService) GetTask(ctx context.Context, userID, taskID int64) (*models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.getVisibleTask(ctx, userID, taskID)
}

func (s *taskService) ListTask(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	switch filter.Role {
	case "", models.TaskRoleOwner, models.TaskRoleAssignee:
	default:
		return nil, models.ErrInvalidRole
	}

	tasks, err := s.taskRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) UpdateTask(ctx context.Context, userID, taskID int64, input models.CreateTaskInput) (*models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	existing, err := s.getVisibleTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}
	if existing.OwnerID != userID {
		return nil, models.ErrForbidden
	}

	task, err := buildTask(input)
	if err != nil {
		return nil, err
	}
	task.ID = taskID

	updated, err := s.taskRepo.Update(ctx, task)
	if err != nil {
		if errors.Is(err, models.ErrAssigneeNotFound) || errors.Is(err, models.ErrTaskNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return updated, nil
}

func (s *taskService) SetComplete(ctx context.Context, userID, taskID int64, completed bool) (*models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if _, err := s.getVisibleTask(ctx, userID, taskID); err != nil {
		return nil, err
	}

	updated, err := s.taskRepo.SetCompleted(ctx, taskID, completed)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to set task status: %w", err)
	}

	return updated, nil
}

func (s *taskService) DeleteTask(ctx context.Context, userID, taskID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	existing, err := s.getVisibleTask(ctx, userID, taskID)
	if err != nil {
		return err
	}

	if existing.OwnerID != userID {
		return models.ErrForbidden
	}

	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (s *taskService) getVisibleTask(ctx context.Context, userID, taskID int64) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	isOwner := task.OwnerID == userID
	isAssignee := task.AssigneeID != nil && *task.AssigneeID == userID
	if !isOwner && !isAssignee {
		return nil, models.ErrTaskNotFound
	}

	return task, nil

}

func buildTask(input models.CreateTaskInput) (models.Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" || utf8.RuneCountInString(title) > maxTitleLength {
		return models.Task{}, models.ErrInvalidTitle
	}

	var description *string
	if input.Description != nil {
		text := strings.TrimSpace(*input.Description)
		if utf8.RuneCountInString(text) > maxDescriptionLength {
			return models.Task{}, models.ErrInvalidDescription
		}
		if text != "" {
			description = &text
		}
	}

	priority := models.PriorityDefault
	if input.Priority != nil {
		priority = *input.Priority
	}
	if priority < models.PriorityLow || priority > models.PriorityHigh {
		return models.Task{}, models.ErrInvalidPriority
	}

	return models.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		Deadline:    input.Deadline,
		AssigneeID:  input.AssigneeID,
	}, nil
}
