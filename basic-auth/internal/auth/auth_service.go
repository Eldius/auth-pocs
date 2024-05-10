package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/model"
	"github.com/eldius/auth-pocs/basic-auth/internal/repository"
	"log/slog"
)

var (
	NotAuthorizedErr = errors.New("user lookup")
)

type ctxAuthUserInfo struct {
}

type Service interface {
	AuthenticateUser(ctx context.Context, username, password string) (*model.User, context.Context, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) Service {
	return &authService{repo: repo}
}

func (s authService) AuthenticateUser(ctx context.Context, username, password string) (*model.User, context.Context, error) {
	user, err := s.repo.FindByUser(ctx, username)
	if err != nil {
		err = fmt.Errorf("fetching user info")
		slog.With("error", err).Error("AuthHandler")
		return nil, ctx, err
	}

	if user.Pass != password {
		slog.With("error", NotAuthorizedErr).Error("AuthHandler")
		return nil, ctx, NotAuthorizedErr
	}

	usr := model.User{
		ID:   user.ID,
		User: user.User,
		Pass: user.Pass,
	}
	ctx = context.WithValue(ctx, ctxAuthUserInfo{}, usr)

	return &usr, ctx, nil
}
