package auth

import (
	"context"
	"github.com/eldius/auth-pocs/basic-auth/internal/config"
	"github.com/eldius/auth-pocs/basic-auth/internal/model"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence/repository"
	helperpersistence "github.com/eldius/auth-pocs/helper-library/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthService_CreateUser(t *testing.T) {
	t.Run("given a valid user should hash it's password", func(t *testing.T) {
		db := helperpersistence.DB(config.GetDBConfig())
		helperpersistence.InitDB(db, persistence.DBMigrationsFS, persistence.DBMigrationsRoot, persistence.DBMigrationsDialect)
		r := repository.NewUserRepository(db)
		s := newService(r)

		t.Cleanup(func() {
			_ = db.Close()
		})

		ctx := context.Background()

		plainPassword := "12345"
		plainUser := "admin"

		user := model.User{
			User: plainUser,
			Pass: plainPassword,
		}

		err := s.CreateUser(ctx, &user)
		assert.Nil(t, err)
		assert.Equal(t, plainUser, user.User)
		assert.NotEqual(t, plainPassword, user.Pass)
	})
}

func TestAuthService_AuthenticateUser(t *testing.T) {
	t.Run("given a valid user should validate it's password with success", func(t *testing.T) {
		db := helperpersistence.DB(config.GetDBConfig())
		helperpersistence.InitDB(db, persistence.DBMigrationsFS, persistence.DBMigrationsRoot, persistence.DBMigrationsDialect)
		r := repository.NewUserRepository(db)
		s := newService(r)

		t.Cleanup(func() {
			_ = db.Close()
		})

		ctx := context.Background()

		plainPassword := "12345"
		plainUser := "admin"

		user := model.User{
			User: plainUser,
			Pass: plainPassword,
		}

		err := s.CreateUser(ctx, &user)
		assert.Nil(t, err)
		assert.Equal(t, plainUser, user.User)
		assert.NotEqual(t, plainPassword, user.Pass)

		authUser, ctx, err := s.AuthenticateUser(context.Background(), plainUser, plainPassword)
		assert.Nil(t, err)
		assert.Equal(t, plainUser, authUser.User)
		assert.Empty(t, authUser.Pass)

		ctxUser := UserFromContext(ctx)

		assert.Equal(t, user.User, ctxUser.User)
		assert.Equal(t, user.ID, ctxUser.ID)
		assert.Empty(t, ctxUser.Pass)
	})

}
