package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/glenntam/todoken/internal/model"

	"github.com/alexedwards/argon2id"
)

type Services struct {
	model *model.Queries
}

func NewServices(models *model.Queries) {
	return &Services{
		model: models,
	}
}

func getCurrentDate() *time.Time {
	return &time.Now().Format("2006-01-02 15:04:05")
}

func (s *Services) Register(ctx context.Context, email, password string) (*model.User, error) {
	// TODO: adjust DefaultParams for production use
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}
	err = s.model.CreateUser(r.Context(), model.CreateUserParams{
		Email:  email,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, err
	}
}
