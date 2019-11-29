package wallet

import (
	"context"

	"github.com/williamchand/my-wallet/models"
)

// Repository represent the wallet's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Wallet, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Wallet, error)
	GetByTitle(ctx context.Context, title string) (*models.Wallet, error)
	Update(ctx context.Context, ar *models.Wallet) error
	Store(ctx context.Context, a *models.Wallet) error
	Delete(ctx context.Context, id int64) error
}
