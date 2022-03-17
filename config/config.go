package config

import (
	"github.com/maitesin/marvin/internal/infra/http"
	"github.com/maitesin/marvin/internal/infra/newrelic"
	"github.com/maitesin/marvin/internal/infra/telegram"
	"os"

	"github.com/maitesin/marvin/internal/infra/sql"
)

// Config defines the configuration of the marvin application
type Config struct {
	SQL      sql.Config
	HTTP     http.Config
	Telegram telegram.Config
	NewRelic newrelic.Config
}

// NewConfig is the constructor for the marvin application configuration
func NewConfig() Config {
	return Config{
		SQL: sql.Config{
			URL:          GetEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:54321/marvin"),
			SSLMode:      GetEnvOrDefault("DATABASE_SSL_MODE", "disable"),
			BinaryParams: GetEnvOrDefault("DATABASE_BINARY_PARAMETERS", "yes"),
		},
		HTTP: http.Config{
			Host: GetEnvOrDefault("HOST", "0.0.0.0"),
			Port: GetEnvOrDefault("PORT", "80"),
		},
		Telegram: telegram.Config{
			Token: GetEnvOrDefault("TELEGRAM_TOKEN", ""),
		},
		NewRelic: newrelic.Config{
			Name:    GetEnvOrDefault("NEW_RELIC_NAME", "Marvin"),
			License: GetEnvOrDefault("NEW_RELIC_LICENSE_KEY", ""),
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
