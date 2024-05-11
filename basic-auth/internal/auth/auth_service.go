package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/model"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence/repository"
	"github.com/eldius/auth-pocs/helper-library/auth/helper"
	"log/slog"
	"sync"
)

var (
	NotAuthorizedErr = errors.New("user lookup")
)

type ctxAuthUserInfo struct {
}

type Service interface {
	AuthenticateUser(ctx context.Context, username, password string) (*model.User, context.Context, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) Service {
	return sync.OnceValue(func() Service {
		return newService(repo)
	})()
}

func newService(repo repository.UserRepository) Service {
	return &authService{repo: repo}
}

func (s authService) AuthenticateUser(ctx context.Context, username, password string) (*model.User, context.Context, error) {
	user, err := s.repo.FindByUser(ctx, username)
	if err != nil {
		err = fmt.Errorf("%w: %w", NotAuthorizedErr, err)
		slog.With("error", err).Error("AuthHandler")
		return nil, ctx, err
	}

	if err := helper.ValidatePassword(user.Pass, password); err != nil {
		err = fmt.Errorf("%w: %w", NotAuthorizedErr, err)
		slog.With("error", err).Error("AuthHandler")
		return nil, ctx, err
	}

	usr := model.User{
		ID:   user.ID,
		User: user.User,
	}
	ctx = context.WithValue(ctx, ctxAuthUserInfo{}, usr)

	return &usr, ctx, nil
}

func (s authService) CreateUser(ctx context.Context, user *model.User) error {
	pass, err := helper.HashPassword(user.Pass)
	if err != nil {
		err = fmt.Errorf("generating bcrypt hash: %w", err)
		return err
	}
	user.Pass = pass

	savedUser, err := s.repo.Save(ctx, user)
	if err != nil {
		err = fmt.Errorf("saving user: %w", err)
		return err
	}
	user.ID = savedUser.ID

	return nil
}

// UserFromContext returns user information stored in context.Context
func UserFromContext(ctx context.Context) *model.User {
	u, ok := ctx.Value(ctxAuthUserInfo{}).(model.User)
	if !ok {
		return nil
	}

	return &u
}
