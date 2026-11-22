package models

import (
	"time"
)

const (
	PriorityLow     = 1
	PriorityMedium  = 2 //default
	PriorityHigh    = 3
	PriorityDefault = PriorityMedium
)

const (
	TaskRoleOwner    = "owner"
	TaskRoleAssignee = "assignee"
)

type Task struct {
	ID          int64
	Title       string
	Description *string
	Completed   bool
	Priority    int
	Deadline    *time.Time
	OwnerID     int64
	AssigneeID  *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTaskInput struct {
	Title       string
	Description *string
	Priority    *int
	Deadline    *time.Time
	AssigneeID  *int64
}

type UpdateTaskInput struct {
	Title       string
	Description *string
	Priority    int
	Deadline    *time.Time
}

type TaskFilter struct {
	UserID    int64
	Role      string
	Completed *bool
}
