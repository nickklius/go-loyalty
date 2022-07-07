package controller

import (
	"errors"
	"net/http"

	"github.com/nickklius/go-loyalty/internal/usecase"
)

func parseError(err error) (int, string) {
	if errors.Is(err, usecase.ErrDBDuplicatedEntry) {
		return http.StatusConflict, err.Error()
	}

	return http.StatusInternalServerError, err.Error()
}
