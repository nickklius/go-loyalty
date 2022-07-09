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

	s := &http.Server{
		Handler: nil,
		Addr:    cfg.App.RunAddress,
	}

	//fmt.Println(s.Addr, cfg.App.RunAddress)
	//err = s.ListenAndServe()
	//if err != nil {
	//	logger.Error("app - Run - httpServer.Shutdown: " + err.Error())
	//}

	//http.ListenAndServe(cfg.App.RunAddress, nil)
	s.ListenAndServe()

	//app.Run(cfg, logger)
}

//
