package wallet

import (
	"context"

	"github.com/williamchand/my-wallet/models"
)

// Usecase represent the wallet's usecases
type Usecase interface {
	EnableWallet(ctx context.Context, authorization string) (*models.FetchWallet, error)
	FetchWallet(ctx context.Context, authorization string) (*models.FetchWallet, error)
	AddWallet(ctx context.Context, req *models.ReqTransaction, authorization string) (*models.TransactionDeposit, error)
	WithdrawWallet(ctx context.Context, req *models.ReqTransaction, authorization string) (*models.TransactionWithdraw, error)
	DisableWallet(ctx context.Context, isDisabled bool, authorization string) (*models.WalletDisabled, error)
	InitWallet(ctx context.Context, customer_id string) (*models.FetchWallet, error)
}
