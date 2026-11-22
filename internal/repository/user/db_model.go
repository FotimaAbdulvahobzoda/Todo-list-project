package user

import (
	"time"

	"todo_app/internal/models"
)

type dbModel struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	CreatedAt    time.Time `db:"created_at"`
}

func (d dbModel) toModel() models.User {
	return models.User{
		ID:           d.ID,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		FirstName:    d.FirstName,
		LastName:     d.LastName,
		CreatedAt:    d.CreatedAt,
	}
}
