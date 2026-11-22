package user

import (
	"context"
	"errors"
	"fmt"

	"todo_app/internal/models"
	"todo_app/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolationCode = "23505"
	userColumns         = "id, email, password_hash, first_name, last_name, created_at"
)

var _ repository.UserRepository = (*userRepo)(nil)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, first_name, last_name) VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`,
		user.Email, user.PasswordHash, user.FirstName, user.LastName,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == uniqueViolationCode {
			return 0, models.ErrEmailAlreadyExists
		}
		return 0, fmt.Errorf("create user: %w", err)
	}
	return user.ID, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := r.getOne(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := r.getOne(ctx, `SELECT `+userColumns+` FROM users WHERE lower(email) = lower($1)`, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+userColumns+` FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	dbUsers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dbModel])
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	users := make([]models.User, 0, len(dbUsers))
	for _, d := range dbUsers {
		users = append(users, d.toModel())
	}

	return users, nil
}

func (r *userRepo) getOne(ctx context.Context, query string, args ...any) (*models.User, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	user := m.toModel()
	return &user, nil
}
