package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/model"
	"github.com/jmoiron/sqlx"
)

var (
	AuthUserNotFoundErr = errors.New("user lookup")
)

type UserRepository interface {
	FindByUser(ctx context.Context, user string) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUser(ctx context.Context, user string) (*model.User, error) {
	rows, err := r.db.NamedQueryContext(ctx, `select * from auth_users where user = :user`, map[string]interface{}{
		"user": user,
	})
	if err != nil {
		err = fmt.Errorf("user lookup: %w", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		return nil, AuthUserNotFoundErr
	}
	var u model.User
	if err := rows.StructScan(&u); err != nil {
		err = fmt.Errorf("parsisng user info: %w", err)
		return nil, err
	}

	return &u, nil
}
