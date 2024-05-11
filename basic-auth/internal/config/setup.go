package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	serviceName = "basic-auth"
)

var (
	logKeys = []string{
		"host",
		"service.name",
		"level",
		"message",
		"time",
		"error",
		"source",
		"function",
		"file",
		"line",
	}
)

// Setup configures app parameters
func Setup(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mock-server-go" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".basic-auth"))
		viper.AddConfigPath(filepath.Join(home))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("basic-auth")
	}

	SetDefaults()
	MapEnvVars()

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
	//else {
	//	return err
	//}

	if err := SetupLogs(); err != nil {
		err = fmt.Errorf("failed to configure logs: %w", err)
		return err
	}

	return nil
}

// SetDefaults sets default configuration values
func SetDefaults() {
	viper.SetDefault("auth.key", "1234567890")
}

// MapEnvVars sets up environment variables mapping
func MapEnvVars() {
	viper.SetEnvPrefix("basic")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
