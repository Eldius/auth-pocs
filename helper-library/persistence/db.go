package persistence

import (
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"sync"
)

type DBConfig struct {
	Engine string
	URL    string
}

func newPool(cfg DBConfig) *sqlx.DB {
	slog.With(slog.String("engine", cfg.Engine)).Info("connecting to database")
	return sync.OnceValue(func() *sqlx.DB {
		db, err := sqlx.Open(cfg.Engine, cfg.URL)
		if err != nil {
			err = fmt.Errorf("opening database connection: %w", err)
			panic(err)
		}
		return db
	})()
}
