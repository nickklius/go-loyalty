package app

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/handler"
	"github.com/nickklius/go-loyalty/internal/httpserver"
	"github.com/nickklius/go-loyalty/internal/storage/jobstorage"
	"github.com/nickklius/go-loyalty/internal/storage/postgres"
	"github.com/nickklius/go-loyalty/internal/usecase"
	"github.com/nickklius/go-loyalty/internal/usecase/repo"
	"github.com/nickklius/go-loyalty/internal/worker"
)

type App struct {
	cfg config.Config
}

func Run(cfg *config.Config, logger *zap.Logger) {
	Migrate(cfg.DB.DatabaseURI)

	pg, err := postgres.New(cfg.DB.DatabaseURI)
	if err != nil {
		logger.Error(err.Error())
	}
	defer pg.Close()

	pgRepository := repo.NewPostgresRepository(pg)
	jobStorage := jobstorage.NewJobStorage()
	jobRepository := repo.NewJobRepository(jobStorage)

	useCases := usecase.New(
		pgRepository,
		jobRepository,
	)

	w := worker.NewWorker(pgRepository, jobRepository, logger, cfg)

	go func() {
		w.Run()
	}()

	h := chi.NewRouter()
	handler.NewRouter(h, logger, useCases, cfg)

	httpServer := httpserver.New(h, httpserver.Port(cfg.App.RunAddress))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
		//w.Done()
	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify: " + err.Error())
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown: " + err.Error())
	}
}
