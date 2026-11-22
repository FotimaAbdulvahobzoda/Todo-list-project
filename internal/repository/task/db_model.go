package task

import (
	"time"
	"todo_app/internal/models"
)

type dbModel struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Completed   bool       `db:"completed"`
	Priority    int        `db:"priority"`
	Deadline    *time.Time `db:"deadline"`
	OwnerID     int64      `db:"owner_id"`
	AssigneeID  *int64     `db:"assignee_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (d dbModel) toModel() models.Task {
	return models.Task{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Completed:   d.Completed,
		Priority:    d.Priority,
		Deadline:    d.Deadline,
		OwnerID:     d.OwnerID,
		AssigneeID:  d.AssigneeID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
