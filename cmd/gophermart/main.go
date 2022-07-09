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

	//s := &http.Server{
	//	//Handler: h,
	//	Addr: ":8080",
	//}

	//fmt.Println(s.Addr, cfg.App.RunAddress)
	//err = s.ListenAndServe()
	//if err != nil {
	//	logger.Error("app - Run - httpServer.Shutdown: " + err.Error())
	//}

	//http.ListenAndServe(cfg.App.RunAddress, nil)
	srv := &http.Server{Addr: cfg.App.RunAddress, Handler: nil}
	srv.ListenAndServe()

	//app.Run(cfg, logger)
}

//
