package repository

import (
	"embed"
	"errors"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
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
	db, err := sqlx.Open("sqlite", ":memory:")
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
