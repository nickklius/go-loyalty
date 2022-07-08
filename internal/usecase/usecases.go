package usecase

import (
	"context"
	"fmt"

	"github.com/nickklius/go-loyalty/internal/entity"
)

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) CreateUser(ctx context.Context, user entity.User) error {
	err := uc.repo.StoreUser(ctx, user)
	if err != nil {
		return fmt.Errorf("UseCase - CreateUser - u.repo.StoreUser: %w", err)
	}

	return nil
}

func (uc *UseCase) CheckPassword(ctx context.Context, user entity.User) (entity.User, error) {
	var u entity.User

	u, err := uc.repo.CheckUser(ctx, user)
	if err != nil {
		return u, fmt.Errorf("UseCase - CheckPassword - u.repo.CheckUser: %w", err)
	}

	return u, nil
}

func (uc *UseCase) CreateOrder(ctx context.Context, order entity.Order) error {
	err := uc.repo.StoreOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("UseCase - CreateOrder - u.repo.StoreOrder: %w", err)
	}
	return nil
}
