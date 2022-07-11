package usecase

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var (
	ErrDBNotFound        = errors.New("not found")
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrDBBuildQuery      = errors.New("query not valid")

	ErrDBDuplicateOrder                = errors.New("order was already registered")
	ErrDBDuplicateOrderByUserItself    = errors.New("order was already registered by user itself")
	ErrDBNotEnoughBalanceForWithdrawal = errors.New("not enough balance for withdrawal")
)

func ParseError(err error) error {
	if err == pgx.ErrNoRows {
		return ErrDBNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ErrDBDuplicatedEntry
		case pgerrcode.UndefinedColumn:
			return ErrDBBuildQuery
		}
	}

	return err
}
