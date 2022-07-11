package usecase

import (
	"context"

	"github.com/nickklius/go-loyalty/internal/entity"
)

type (
	User interface {
		CreateUser(ctx context.Context, user entity.User) error
		CheckPassword(ctx context.Context, user entity.User) (entity.User, error)
		GetUserBalance(ctx context.Context, userID string) (entity.UserBalance, error)
		WithdrawFromBalance(ctx context.Context, withdraw entity.Withdraw) error
	}

	Order interface {
		CreateOrder(ctx context.Context, order entity.Order) error
		GetOrdersByUserID(ctx context.Context, userID string) ([]entity.Order, error)
	}

	Repository interface {
		StoreUser(ctx context.Context, user entity.User) error
		CheckUser(ctx context.Context, user entity.User) (entity.User, error)
		StoreOrder(ctx context.Context, order entity.Order) error
		GetOrders(ctx context.Context, userID string) ([]entity.Order, error)
		GetBalance(ctx context.Context, userID string) (entity.UserBalance, error)
		Withdraw(ctx context.Context, withdraw entity.Withdraw) error
	}
)
