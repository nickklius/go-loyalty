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

	UserRepo interface {
		Store(ctx context.Context, user entity.User) error
		Check(ctx context.Context, user entity.User) (entity.User, error)
	}
)
