package persistence

import (
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sync"
)

func newPool(driver, url string) *sqlx.DB {
	db, err := sqlx.Open(driver, url)
	if err != nil {
		err = fmt.Errorf("opening database connection: %w", err)
		panic(err)
	}
	return db
}

func DB(driver, url string) *sqlx.DB {
	return sync.OnceValue(func() *sqlx.DB {
		return newPool(driver, url)
	})()
}
