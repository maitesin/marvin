package config

import (
	"github.com/maitesin/marvin/internal/infra/telegram"
	"os"

	"github.com/maitesin/marvin/internal/infra/sql"
)

// Config defines the configuration of the marvin application
type Config struct {
	SQL      sql.Config
	Telegram telegram.Config
}

// NewConfig is the constructor for the marvin application configuration
func NewConfig() Config {
	return Config{
		SQL: sql.Config{
			URL:          GetEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:54321/marvin"),
			SSLMode:      GetEnvOrDefault("DATABASE_SSL_MODE", "disable"),
			BinaryParams: GetEnvOrDefault("DATABASE_BINARY_PARAMETERS", "yes"),
		},
		Telegram: telegram.Config{
			Token: GetEnvOrDefault("TELEGRAM_TOKEN", ""),
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
