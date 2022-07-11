package usecase

import (
	"context"

	"github.com/nickklius/go-loyalty/internal/entity"
)

type (
	Repository interface {
		StoreUser(ctx context.Context, user entity.User) error
		CheckUser(ctx context.Context, user entity.User) (entity.User, error)
		StoreOrder(ctx context.Context, order entity.Order) error
		GetOrders(ctx context.Context, userID string) ([]entity.Order, error)
		GetBalance(ctx context.Context, userID string) (entity.UserBalance, error)
		Withdraw(ctx context.Context, withdraw entity.Withdraw) error
		GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdraw, error)
		UpdateOrderStatus(ctx context.Context, order entity.Order) error
	}

	JobRepository interface {
		AddJob(ctx context.Context, job entity.Job) error
		GetJobs() ([]entity.Job, error)
		DeleteJob(job entity.Job) error
	}
)
