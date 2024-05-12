package persistence

import (
	"embed"
	"github.com/eldius/auth-pocs/basic-auth/internal/config"
	helperpersistence "github.com/eldius/auth-pocs/helper-library/persistence"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var DBMigrationsFS embed.FS

const (
	DBMigrationsRoot    = "migrations"
	DBMigrationsDialect = "sqlite3"
)

func InitDB() *sqlx.DB {
	return helperpersistence.InitDB(
		helperpersistence.DB(config.GetDBConfig()),
		DBMigrationsFS,
		DBMigrationsRoot,
		DBMigrationsDialect,
	)
}
