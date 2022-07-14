package handler

import (
	"errors"
	"net/http"

	"github.com/nickklius/go-loyalty/internal/usecase"
)

func parseError(err error) (int, string) {
	if errors.Is(err, usecase.ErrDBDuplicatedEntry) {
		return http.StatusConflict, err.Error()
	}

	if errors.Is(err, usecase.ErrDBDuplicateOrderByUserItself) {
		return http.StatusOK, err.Error()
	}

	if errors.Is(err, usecase.ErrDBDuplicateOrder) {
		return http.StatusConflict, err.Error()
	}

	if errors.Is(err, usecase.ErrDBNotEnoughBalanceForWithdrawal) {
		return http.StatusPaymentRequired, err.Error()
	}

	return http.StatusInternalServerError, err.Error()
}
