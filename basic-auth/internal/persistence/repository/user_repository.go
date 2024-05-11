package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/model"
	"github.com/jmoiron/sqlx"
	"sync"
)

var (
	AuthUserNotFoundErr = errors.New("user lookup")
)

type UserRepository interface {
	FindByUser(ctx context.Context, user string) (*model.User, error)
	Save(ctx context.Context, user *model.User) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return sync.OnceValue(func() UserRepository {
		return &userRepository{db: db}
	})()
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

func (r *userRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := r.db.NamedExecContext(ctx, `insert into auth_users (user, pass) values (:user, :pass)`, user)
	if err != nil {
		err = fmt.Errorf("saving user info: %w", err)
		return nil, err
	}
	nu, err := r.FindByUser(ctx, user.User)
	if err != nil {
		err = fmt.Errorf("finding newly created user info: %w", err)
		return nil, err
	}

	return nu, nil
}
