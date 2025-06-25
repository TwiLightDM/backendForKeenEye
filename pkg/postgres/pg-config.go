package postgres

import "time"

type Config struct {
	MaxPoolSize    int    `mapstructure:"max_pool_size"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	MigrationsPath string `mapstructure:"migrations_path"`

	RetryConnectionAttempts int           `mapstructure:"retry_connection_attempts"`
	RetryConnectionTimeout  time.Duration `mapstructure:"retry_connection_timeout"`
}
