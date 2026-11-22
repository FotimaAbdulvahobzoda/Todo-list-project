package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrTaskNotFound       = errors.New("task not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong    = errors.New("password must be at most 72 bytes")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidDescription = errors.New("description must be at most 1000 characters")
	ErrInvalidTitle       = errors.New("title is required and must be at most 100 characters.")
	ErrAlreadyFinished    = errors.New("task is already finished")
	ErrNotAvailable       = errors.New("task is not available")
	ErrAssigneeNotFound   = errors.New("assignee not found")
	ErrInvalidPriority    = errors.New("invalid priority")
	ErrInvalidRole        = errors.New("invalid role")
)
