package wallet

import (
	"context"

	"github.com/williamchand/my-wallet/models"
)

// Repository represent the wallet's repository contract
type Repository interface {
	EnableWallet(ctx context.Context, id string) (*models.FetchWallet, error)
	FetchWallet(ctx context.Context, id string) (*models.FetchWallet, error)
	AddWallet(ctx context.Context, req *models.ReqTransaction, id string) (*models.TransactionDeposit, error)
	WithdrawWallet(ctx context.Context, req *models.ReqTransaction, id string) (*models.TransactionWithdraw, error)
	DisableWallet(ctx context.Context, isDisabled bool, id string) (*models.WalletDisabled, error)
	InitWallet(ctx context.Context, customer_id string) (*models.FetchWallet, error)
}
