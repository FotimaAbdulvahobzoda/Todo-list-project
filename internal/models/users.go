package models

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
}

type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type LoginInput struct {
	Email    string
	Password string
}
