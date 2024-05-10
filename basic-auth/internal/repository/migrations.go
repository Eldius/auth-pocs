package repository

import (
	"embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"log/slog"
	"sync"
)

var (
	MigrationsExecutionFailedErr = errors.New("executing migrations")
)

//go:embed migrations/*.sql
var dbMigrations embed.FS

func newPool() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		err = fmt.Errorf("opening database connection: %w", err)
		panic(err)
	}
	return db
}

func DB() *sqlx.DB {
	return sync.OnceValue(newPool)()
}

func InitDB(db *sqlx.DB) *sqlx.DB {
	log := slog.Default()
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "migrations",
	}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		err = fmt.Errorf("%w: %w", MigrationsExecutionFailedErr, err)
		log.With("error", err).Error("Error loading migration files")
		panic(err)
	}
	log.With(slog.Int("migrations_executed", n)).Info("MigrationsCompleted")

	return db
}
