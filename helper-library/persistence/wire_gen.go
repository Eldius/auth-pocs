// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package persistence

import (
	"github.com/jmoiron/sqlx"
)

import (
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/lib/pq"
)

// Injectors from wire.go:

func DB(cfg DBConfig) *sqlx.DB {
	db := newPool(cfg)
	return db
}