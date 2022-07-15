package handler

import (
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

func NewRouter(h *chi.Mux, l *zap.Logger, u *usecase.UseCase, c *config.Config) {
	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	newRoutes(h, l, u, c)
}
