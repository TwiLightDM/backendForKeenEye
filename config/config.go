package config

import (
	"backendForKeenEye/pkg/postgres"
	"fmt"

	"github.com/spf13/viper"
	"os"
	"time"
)

type (
	Config struct {
		Http       `mapstructure:"http"`
		Encryption `mapstructure:"encryption"`
		Postgres   postgres.Config `mapstructure:"pg"`
		JWT        `mapstructure:"jwt"`
	}

	Postgres struct {
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

	Http struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	Encryption struct {
		Salt int `mapstructure:"salt_length"`
	}

	JWT struct {
		Key         string        `mapstructure:"key"`
		AccessTime  time.Duration `mapstructure:"access_time"`
		RefreshTime time.Duration `mapstructure:"refresh_time"`
	}
)

func NewConfig() (*Config, error) {
	cfg := Config{}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	for _, k := range v.AllKeys() {
		anyValue := v.Get(k)
		str, ok := anyValue.(string)
		if !ok {
			continue
		}

		replaced := os.ExpandEnv(str)
		v.Set(k, replaced)
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling file: %w", err))
	}

	return &cfg, nil
}
