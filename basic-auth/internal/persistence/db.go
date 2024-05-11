package persistence

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/config"
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sync"
)

func newPool(cfg config.DBConfig) *sqlx.DB {
	return sync.OnceValue(func() *sqlx.DB {
		db, err := sqlx.Open(cfg.Engine, cfg.URL)
		if err != nil {
			err = fmt.Errorf("opening database connection: %w", err)
			panic(err)
		}
		return db
	})()
}

//
//func DB(cfg config.DBConfig) *sqlx.DB {
//	return sync.OnceValue(func() *sqlx.DB {
//		return newPool(cfg.Engine, cfg.URL)
//	})()
//}
