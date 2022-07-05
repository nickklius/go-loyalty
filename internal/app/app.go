package app

import (
	"go.uber.org/zap"

	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/interfaces"
	"github.com/nickklius/go-loyalty/internal/storage/postgres"
)

type App struct {
	cfg  config.Config
	repo interfaces.Repository
}

func Run(cfg *config.Config, logger *zap.Logger) {
	pg, err := postgres.New(cfg.DB.DatabaseURI)
	if err != nil {
		logger.Error(err.Error())
	}
	defer pg.Close()
}
