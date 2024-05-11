//go:build wireinject
// +build wireinject

// wire.go
package persistence

import (
	"github.com/eldius/auth-pocs/basic-auth/internal/config"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func DB(cfg config.DBConfig) *sqlx.DB {
	wire.Build(newPool)
	return nil
}
