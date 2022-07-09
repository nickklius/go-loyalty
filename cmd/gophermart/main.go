package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nickklius/go-loyalty/config"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(err.Error())
	}

	srv := &http.Server{Addr: cfg.App.RunAddress, Handler: nil}

	srv.ListenAndServe()
	//app.Run(cfg, logger)
}
