package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	HTTP   HTTP
	Log    Log
	PG     PG
	Redis  Redis
	JWT    JWT
	Hasher Hasher
	TestPG TestPG
}

type (
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}
	Log struct {
		Level  string `env-required:"true" env:"LOG_LEVEL"`
		Output string `env-required:"true" env:"LOG_OUTPUT"`
	}
	PG struct {
		MaxPoolSize int    `env-required:"true" env:"PG_MAX_POOL_SIZE"`
		Url         string `env-required:"true" env:"PG_URL"`
	}
	Redis struct {
		Url         string `env-required:"true" env:"REDIS_URL"`
		MaxPoolSize int    `env-required:"true" env:"REDIS_MAX_POOL_SIZE"`
	}
	JWT struct {
		SignKey  string        `env-required:"true" env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `env-required:"true" env:"TOKEN_TTL"`
	}
	Hasher struct {
		Salt string `env-required:"true" env:"HASH_SALT"`
	}
	TestPG struct {
		Url string `env-required:"true" env:"TEST_PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadEnv(c); err != nil {
		return nil, fmt.Errorf("error reading config env: %w", err)
	}
	return c, nil
}
