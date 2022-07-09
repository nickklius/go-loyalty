package main

import (
	"fmt"
	"net"
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
		Addr:    cfg.RunAddress,
	}

	//fmt.Println(s.Addr, cfg.App.RunAddress)
	//err = s.ListenAndServe()
	//if err != nil {
	//	logger.Error("app - Run - httpServer.Shutdown: " + err.Error())
	//}

	fmt.Println(cfg.Secret)

	//http.ListenAndServe(cfg.App.RunAddress, nil)
	s.Addr = net.JoinHostPort("", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		return
	}

	//app.Run(cfg, logger)
}

//
