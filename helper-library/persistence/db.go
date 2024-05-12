package persistence

import (
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sync"
)

type DBConfig struct {
	Engine string
	URL    string
}

func newPool(cfg DBConfig) *sqlx.DB {
	return sync.OnceValue(func() *sqlx.DB {
		db, err := sqlx.Open(cfg.Engine, cfg.URL)
		if err != nil {
			err = fmt.Errorf("opening database connection: %w", err)
			panic(err)
		}
		return db
	})()
}
