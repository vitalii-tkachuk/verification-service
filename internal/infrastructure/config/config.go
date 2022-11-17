package config

import (
	"fmt"
	"time"
)

// Config represents application config defined in environment variables.
type Config struct {
	Port             uint8         `default:"80" split_words:"true"`
	ShutdownTimeout  time.Duration `default:"10s" split_words:"true"`
	DatabaseUser     string        `default:"user" split_words:"true"`
	DatabasePassword string        `default:"password" split_words:"true"`
	DatabaseHost     string        `default:"localhost" split_words:"true"`
	DatabasePort     uint          `default:"5432" split_words:"true"`
	DatabaseName     string        `default:"database_name" split_words:"true"`
	DatabaseTimeout  time.Duration `default:"5s" split_words:"true"`
}

// PostgresDatabaseDsn transform database environment variables to PostgreSQL DSN connection string.
func (c Config) PostgresDatabaseDsn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
	)
}
