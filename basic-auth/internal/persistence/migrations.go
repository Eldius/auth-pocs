package persistence

import (
	"embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"log/slog"
)

var (
	MigrationsExecutionFailedErr = errors.New("executing migrations")
)

//go:embed migrations/*.sql
var dbMigrations embed.FS

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
