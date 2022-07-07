package usecase

import (
	"context"
	"fmt"

	"github.com/nickklius/go-loyalty/internal/entity"
)

type UserUseCase struct {
	repo UserRepo
}

func New(r UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user entity.User) error {
	err := uc.repo.Store(ctx, user)
	if err != nil {
		return fmt.Errorf("UserUseCase - Store - u.repo.Store: %w", err)
	}

	return nil
}

func (uc *UserUseCase) CheckPassword(ctx context.Context, user entity.User) (entity.User, error) {
	var u entity.User

	u, err := uc.repo.Check(ctx, user)
	if err != nil {
		return u, fmt.Errorf("UserUseCase - CheckPassword - u.repo.CheckPassword: %w", err)
	}

	return u, nil
}
