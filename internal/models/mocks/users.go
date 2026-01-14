package mocks

import (
	"context"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/models"
)

type UserModel struct{}

func (um *UserModel) Insert(ctx context.Context, name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (um *UserModel) Authenticate(ctx context.Context, email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (um *UserModel) Exists(ctx context.Context, id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (um *UserModel) Get(ctx context.Context, id int) (*models.User, error) {
	switch id {
	case 1:
		return &models.User{
			ID:      1,
			Name:    "agustin",
			Email:   "agustin@test.com",
			Created: time.Now(),
		}, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (um *UserModel) PasswordUpdate(ctx context.Context, id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return models.ErrInvalidCredentials
		}
		return nil
	}
	return models.ErrNoRecord
}
