package users

import (
	"comune/apps/api/internal/platform/db"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

var ErrEmailRequired = errors.New("email is required")
var ErrEmailTaken = errors.New("email is already taken")

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateUser(ctx context.Context, txn pgx.Tx, params NewUserParams) (User, error) {
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))
	params.FirstName = strings.TrimSpace(params.FirstName)
	params.LastName = strings.TrimSpace(params.LastName)
	params.Phone = strings.TrimSpace(params.Phone)

	if params.Email == "" {
		return User{}, ErrEmailRequired
	}

	user, err := insertUser(ctx, txn, params)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return User{}, ErrEmailTaken
		}
		return User{}, err
	}

	return user, nil
}
