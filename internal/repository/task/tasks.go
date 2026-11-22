package task

import (
	"context"
	"errors"
	"fmt"
	"todo_app/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	foreignKeyViolation = "23503"

	assigneeForeignKey = "tasks_assignee_id_fkey"

	taskColumns = "id, title, description, completed, priority, deadline, owner_id, assignee_id, created_at, updated_at"
)

func (t *taskRepo) CreateTask(ctx context.Context, tasks models.Task) (*models.Task, error) {
	created, err := t.getOne(ctx,
		`INSERT INTO tasks(title, description, completed, priority, deadline, owner_id, assignee_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING `+taskColumns,
		tasks.Title, tasks.Description, tasks.Completed, tasks.Priority, tasks.Deadline, tasks.OwnerID, tasks.AssigneeID)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return created, nil
}

func (t *taskRepo) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	task, err := t.getOne(ctx, `SELECT `+taskColumns+` FROM tasks WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return task, nil
}

func (t *taskRepo) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	args := []any{filter.UserID}

	var where string
	switch filter.Role {
	case models.TaskRoleOwner:
		where = "owner_id = $1"
	case models.TaskRoleAssignee:
		where = "assignee_id = $1"
	default:
		where = "(owner_id = $1 OR assignee_id = $1)"
	}

	if filter.Completed != nil {
		args = append(args, *filter.Completed)
		where += fmt.Sprintf(" AND completed = $%d", len(args))
	}

	query := `SELECT ` + taskColumns + ` FROM tasks WHERE ` + where +
		` ORDER BY completed, deadline NULLS LAST, id`

	rows, err := t.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	dbTasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	tasks := make([]models.Task, 0, len(dbTasks))
	for _, m := range dbTasks {
		tasks = append(tasks, m.toModel())
	}
	return tasks, nil
}

func (t *taskRepo) Update(ctx context.Context, task models.Task) (*models.Task, error) {
	updated, err := t.getOne(ctx,
		`UPDATE tasks
	 SET title = $2, description = $3, priority = $4, deadline = $5, assignee_id = $6, updated_at = NOW() 
	 WHERE id = $1 
	 RETURNING `+taskColumns,
		task.ID, task.Title, task.Description, int16(task.Priority), task.Deadline, task.AssigneeID,
	)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	return updated, nil
}

func (t *taskRepo) SetCompleted(ctx context.Context, taskID int64, completed bool) (*models.Task, error) {
	updated, err := t.getOne(ctx,
		`UPDATE tasks SET completed = $2, updated_at = NOW() WHERE id = $1 RETURNING `+taskColumns,
		taskID, completed)
	if err != nil {
		return nil, fmt.Errorf("set task complete: %w", err)
	}
	return updated, nil
}

func (t *taskRepo) Delete(ctx context.Context, taskID int64) error {
	deleted, err := t.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if deleted.RowsAffected() == 0 {
		return models.ErrTaskNotFound
	}
	return nil
}

func (t *taskRepo) getOne(ctx context.Context, query string, args ...any) (*models.Task, error) {
	rows, err := t.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, mapError(err)
	}
	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return nil, mapError(err)
	}
	task := m.toModel()
	return &task, nil
}

func mapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrTaskNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation && pgErr.ConstraintName == assigneeForeignKey {
		return models.ErrAssigneeNotFound
	}
	return err
}
