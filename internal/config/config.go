package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address    string `env:"RUN_ADDRESS"`
	Database   string `env:"DATABASE_URI"`
	AccAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey  string
	Algorithm  string
}

func New() (*Config, error) {
	address := flag.String("a", "localhost:8081", "адрес эндпоинта HTTP-сервера")
	database := flag.String("d", "host=localhost user=postgres password=postgres dbname=diplom sslmode=disable", "строка с адресом подключения к БД")
	accAddress := flag.String("r", "http://localhost:8080", "адрес системы расчёта начислений")
	flag.Parse()

	cfg := &Config{
		Address:    *address,
		Database:   *database,
		AccAddress: *accAddress,
		SecretKey:  "SecretKey",
		Algorithm:  "HS256",
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
