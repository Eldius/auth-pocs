package config

import (
	"github.com/eldius/auth-pocs/helper-library/persistence"
	"github.com/spf13/viper"
	"log/slog"
)

func LogLevel() slog.Level {
	if viper.GetBool("debug") {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func GetDBConfig() persistence.DBConfig {
	return persistence.DBConfig{
		Engine: viper.GetString("db.engine"),
		URL:    viper.GetString("db.url"),
	}
}
