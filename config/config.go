package config

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App     `yaml:"app"`
		Accrual `yaml:"accrual"`
		DB      `yaml:"db"`
		Auth    `yaml:"auth"`
	}

	App struct {
		RunAddress string `env-required:"true" yaml:"run-address" env:"RUN_ADDRESS" `
	}

	Accrual struct {
		AccrualAddress string `env-required:"true" yaml:"accrual-address" env:"ACCRUAL_SYSTEM_ADDRESS"`
	}

	DB struct {
		DatabaseURI string `env-required:"true" yaml:"database-uri" env:"DATABASE_URI"`
	}

	Auth struct {
		Secret         string `env-required:"true" yaml:"secret"`
		AccessLifeTime int64  `env-required:"true" yaml:"access-life-time"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	flag.StringVar(&cfg.App.RunAddress, "a", cfg.App.RunAddress, "Service run address")
	flag.StringVar(&cfg.Accrual.AccrualAddress, "r", cfg.Accrual.AccrualAddress, "Accrual system address")
	flag.StringVar(&cfg.DB.DatabaseURI, "d", cfg.DB.DatabaseURI, "Database URI")
	flag.Parse()

	return cfg, nil
}
