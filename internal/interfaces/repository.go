package interfaces

import (
	"context"

	"github.com/nickklius/go-loyalty/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	CheckPassword(ctx context.Context, user models.User) (*models.User, error)
}
