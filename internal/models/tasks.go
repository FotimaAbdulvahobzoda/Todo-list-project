package models

import (
	"time"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	OwnerID     int64      `json:"owner_id"`
	AssigneeID  *int64     `json:"assignee_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskInput struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  int64
}

type UpdateTaskInput struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
}
