//go:build wireinject
// +build wireinject

// wire.go
package persistence

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func DB(cfg DBConfig) *sqlx.DB {
	wire.Build(newPool)
	return nil
}
