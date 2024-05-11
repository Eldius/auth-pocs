package auth

import (
	"context"
	"github.com/eldius/auth-pocs/basic-auth/internal/model"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthService_CreateUser(t *testing.T) {
	t.Run("given a valid user", func(t *testing.T) {
		db := persistence.DB()
		persistence.InitDB(db)
		r := persistence.NewUserRepository(db)
		s := NewAuthService(r)

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
		t.Logf("old pass: %s", user.Pass)
		assert.Nil(t, err)
		assert.Equal(t, plainUser, user.User)
		assert.NotEqual(t, plainPassword, user.Pass)
	})
}
