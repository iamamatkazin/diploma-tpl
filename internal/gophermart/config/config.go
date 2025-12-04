package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address        string `env:"RUN_ADDRESS"`
	Database       string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	database := flag.String("d", "", "строка с адресом подключения к БД")
	accrualAddress := flag.String("r", "", "адрес системы расчёта начислений")
	flag.Parse()

	cfg := &Config{
		Address:        *address,
		Database:       *database,
		AccrualAddress: *accrualAddress,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
