package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/williamchand/my-wallet/models"
	"github.com/williamchand/my-wallet/wallet"
)

type walletUsecase struct {
	walletRepo     wallet.Repository
	contextTimeout time.Duration
}

// NewWalletUsecase will create new an walletUsecase object representation of wallet.Usecase interface
func NewWalletUsecase(a wallet.Repository, timeout time.Duration) wallet.Usecase {
	return &walletUsecase{
		walletRepo:     a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

func (a *walletUsecase) EnableWallet(c context.Context, authorization string) (*models.FetchWallet, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	data := JWT(authorization)
	res, err := a.walletRepo.EnableWallet(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *walletUsecase) FetchWallet(c context.Context, authorization string) (*models.FetchWallet, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	data := JWT(authorization)
	res, err := a.walletRepo.FetchWallet(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *walletUsecase) AddWallet(c context.Context, req *models.ReqTransaction, authorization string) (*models.TransactionDeposit, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	data := JWT(authorization)
	res, err := a.walletRepo.AddWallet(ctx, req, data.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *walletUsecase) WithdrawWallet(c context.Context, req *models.ReqTransaction, authorization string) (*models.TransactionWithdraw, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	data := JWT(authorization)
	res, err := a.walletRepo.WithdrawWallet(ctx, req, data.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *walletUsecase) DisableWallet(c context.Context, isDisabled bool, authorization string) (*models.WalletDisabled, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	data := JWT(authorization)
	res, err := a.walletRepo.DisableWallet(ctx, isDisabled, data.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *walletUsecase) InitWallet(c context.Context, costumer_id string) (*models.FetchWallet, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.walletRepo.InitWallet(ctx, costumer_id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//this part for translate authorization to id and owner_id I need your API Method to translate this
func JWT(token string) *models.User {
	ss := strings.Fields(token)
	newToken := ss[1]
	wallet := models.User{
		ID:   newToken,
		Name: "william-chandra",
	}
	return &wallet
}
