package usecase

import (
	"context"

	"github.com/nickklius/go-loyalty/internal/entity"
)

type (
	User interface {
		CreateUser(ctx context.Context, user entity.User) error
		CheckPassword(ctx context.Context, user entity.User) (entity.User, error)
	}

	Order interface {
		CreateOrder(ctx context.Context, order entity.Order) error
	}

	Repository interface {
		StoreUser(ctx context.Context, user entity.User) error
		CheckUser(ctx context.Context, user entity.User) (entity.User, error)
		StoreOrder(ctx context.Context, order entity.Order) error
	}
)
