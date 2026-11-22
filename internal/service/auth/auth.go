package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"todo_app/internal/models"
	"todo_app/internal/repository"
	"todo_app/internal/service"
	"todo_app/internal/service/token"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 8
	maxPasswordBytes  = 72
)

type authService struct {
	users  repository.UserRepository
	tokens token.Manager
}

func NewAuthService(users repository.UserRepository, tokens token.Manager) service.AuthService {
	return &authService{
		users:  users,
		tokens: tokens,
	}
}

func (s *authService) Register(ctx context.Context, firstName, lastName, email, password string) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.ToLower(strings.TrimSpace(email))

	if firstName == "" || lastName == "" || email == "" || password == "" {
		return 0, models.ErrInvalidInput
	}

	if utf8.RuneCountInString(password) < minPasswordLength {
		return 0, models.ErrPasswordTooShort
	}

	if len(password) > maxPasswordBytes {
		return 0, models.ErrPasswordTooLong
	}

	// We never store the password itself, only its bcrypt hash.
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	userID, err := s.users.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrEmailAlreadyExists) {
			return 0, err
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (models.TokenPair, error) {
	if err := ctx.Err(); err != nil {
		return models.TokenPair{}, err
	}

	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return models.TokenPair{}, models.ErrInvalidCredentials
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			// Same error as for a wrong password, so nobody can find out which emails are registered.
			return models.TokenPair{}, models.ErrInvalidCredentials
		}

		return models.TokenPair{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.TokenPair{}, models.ErrInvalidCredentials
	}

	accessToken, err := s.tokens.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Authenticate(ctx context.Context, tokenValue string) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	if strings.TrimSpace(tokenValue) == "" {
		return 0, models.ErrInvalidToken
	}

	userID, err := s.tokens.Validate(ctx, tokenValue)
	if err != nil {
		return 0, fmt.Errorf("failed to validate token: %w", err)
	}

	return userID, nil
}
