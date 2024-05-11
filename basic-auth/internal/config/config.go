package config

import (
	"github.com/spf13/viper"
	"log/slog"
)

func LogLevel() slog.Level {
	if viper.GetBool("debug") {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Engine: "sqlite",
		URL:    ":memory:",
	}
}

type DBConfig struct {
	Engine string
	URL    string
}
