package worker

import "errors"

var (
	ErrNoAccessToAccrual   = errors.New("problem with access accrual service")
	ErrAccrualOverloaded   = errors.New("accrual service overloaded")
	ErrOrderNotFound       = errors.New("order not found on accrual service")
	ErrOrderIsInProcessing = errors.New("order processing is not finished")
	ErrOrderIsInvalid      = errors.New("order is invalid")
)
