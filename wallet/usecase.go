package wallet

import (
	"context"

	"github.com/williamchand/my-wallet/models"
)

// Usecase represent the wallet's usecases
type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*models.Wallet, string, error)
	GetByID(ctx context.Context, id int64) (*models.Wallet, error)
	Update(ctx context.Context, ar *models.Wallet) error
	GetByTitle(ctx context.Context, title string) (*models.Wallet, error)
	Store(context.Context, *models.Wallet) error
	Delete(ctx context.Context, id int64) error
}
