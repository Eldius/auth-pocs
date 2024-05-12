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

func InitDB(db *sqlx.DB, f embed.FS, root, dialect string) *sqlx.DB {
	log := slog.Default()
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: f,
		Root:       root,
	}

	n, err := migrate.Exec(db.DB, dialect, migrations, migrate.Up)
	if err != nil {
		err = fmt.Errorf("%w: %w", MigrationsExecutionFailedErr, err)
		log.With("error", err).Error("Error loading migration files")
		panic(err)
	}
	log.With(slog.Int("migrations_executed", n)).Info("MigrationsCompleted")

	return db
}
