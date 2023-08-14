// Package config represents struct Config.
package config

// Config is a structure of environment variables.
type Config struct {
	PostgresPath string `env:"POSTGRES_PATH_POSITIONS"`
}
