package main

import (
	"log"
	"net/http"
)

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := []byte("Hello World!")
	w.Write(data)
}

func main() {
	//logger, _ := zap.NewProduction()
	//defer logger.Sync()
	//
	//cfg, err := config.NewConfig()
	//if err != nil {
	//	logger.Error(err.Error())
	//}

	log.Print("starting app, create handler")
	handler := Handler{}
	log.Print("run listener")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}

	//app.Run(cfg, logger)

}
