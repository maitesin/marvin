package config

import (
	"os"

	"github.com/maitesin/marvin/internal/infra/sql"
)

// Config defines the configuration of the mtga application
type Config struct {
	SQL sql.Config
}

// NewConfig is the constructor for the mtga application configuration
func NewConfig() Config {
	return Config{
		SQL: sql.Config{
			URL:          GetEnvOrDefault("DB_URL", "postgres://postgres:postgres@localhost:54321/marvin"),
			SSLMode:      GetEnvOrDefault("DB_SSL_MODE", "disable"),
			BinaryParams: GetEnvOrDefault("DB_BINARY_PARAMETERS", "yes"),
		},
	}
}

func GetEnvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}

	return defaultValue
}
