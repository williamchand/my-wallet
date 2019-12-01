package models

import (
	"time"
)

// Wallet represent the wallet model
type Wallet struct {
	ID        string    `json:"wallet_id"`
	Status    string    `json:"status"`
	Balance   int64     `json:"balance"`
	OwnedBy   string    `json:"owned_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FetchWallet struct {
	ID        string    `json:"id"`
	OwnedBy   string    `json:"owned_by"`
	Status    string    `json:"status"`
	EnabledAt time.Time `json:"enabled_at"`
	Balance   int64     `json:"balance"`
}

type WalletDisabled struct {
	ID         string    `json:"id"`
	OwnedBy    string    `json:"owned_by"`
	Status     string    `json:"status"`
	DisabledAt time.Time `json:"disabled_at"`
	Balance    int64     `json:"balance"`
}

type Transaction struct {
	ReferenceID string    `json:"reference_id"`
	ID          string    `json:"wallet_id"`
	Type        bool      `json:"type"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransactionDeposit struct {
	ReferenceID string    `json:"reference_id"`
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	DepositBy   string    `json:"deposited_by"`
	DepositAt   time.Time `json:"deposited_at"`
}

type TransactionWithdraw struct {
	ReferenceID string    `json:"reference_id"`
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	WithdrawnBy string    `json:"withdrawn_by"`
	WithdrawnAt time.Time `json:"withdrawn_at"`
}

type ReqTransaction struct {
	ReferenceID string `json:"reference_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required"`
}
