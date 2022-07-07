package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/app"
)

func main() {
	log.Print("shall we work?")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("config loaded, prepared to start app")

	app.Run(cfg, logger)
}
