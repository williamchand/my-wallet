package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/williamchand/my-wallet/models"
	"github.com/williamchand/my-wallet/wallet"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mysqlWalletRepository struct {
	Conn *sql.DB
}

// NewMysqlWalletRepository will create an object that represent the wallet.Repository interface
func NewMysqlWalletRepository(Conn *sql.DB) wallet.Repository {
	return &mysqlWalletRepository{Conn}
}

func (m *mysqlWalletRepository) EnableWallet(ctx context.Context, id string) (*models.FetchWallet, error) {
	query := `UPDATE wallet set status = "enabled" updated_at=? WHERE wallet_id = ? AND status = "disabled"`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rowUpdate, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return nil, err
	}
	affect, err := rowUpdate.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect != 1 {
		return nil, models.ErrAlreadyEnabled
	}

	res, err := m.FetchWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlWalletRepository) FetchWallet(ctx context.Context, id string) (*models.FetchWallet, error) {
	query := `SELECT id, owned_by, status, updated_at, balance 
			  FROM wallet WHERE wallet_id = ? AND status = "enabled"`

	list, err := m.fetchWallet(ctx, query, id)
	if err != nil {
		return nil, err
	}

	res := &models.FetchWallet{}
	if len(list) > 0 {
		res = &models.FetchWallet{
			ID:        list[0].ID,
			OwnedBy:   list[0].OwnedBy,
			Status:    list[0].Status,
			EnabledAt: list[0].UpdatedAt,
			Balance:   list[0].Balance,
		}
	} else {
		return nil, models.ErrDisabled
	}

	return res, nil
}

func (m *mysqlWalletRepository) FetchDisabledWallet(ctx context.Context, id string) (*models.WalletDisabled, error) {
	query := `SELECT id, owned_by, status, updated_at, balance 
			  FROM wallet WHERE wallet_id = ? AND status = "disabled"`

	list, err := m.fetchWallet(ctx, query, id)
	if err != nil {
		return nil, err
	}

	res := &models.WalletDisabled{}
	if len(list) > 0 {
		res = &models.WalletDisabled{
			ID:         list[0].ID,
			OwnedBy:    list[0].OwnedBy,
			Status:     list[0].Status,
			DisabledAt: list[0].UpdatedAt,
			Balance:    list[0].Balance,
		}
	} else {
		return nil, models.ErrNotFound
	}

	return res, nil
}

func (m *mysqlWalletRepository) FetchTransactionAdd(ctx context.Context, id int64) (*models.TransactionDeposit, error) {
	query := `SELECT reference_id, wallet_id, amount, status, created_by, created_at
			  FROM transaction WHERE id = ?`
	list, err := m.fetchTransaction(ctx, query, id)
	if err != nil {
		return nil, err
	}

	res := &models.TransactionDeposit{}
	if len(list) > 0 {
		res = &models.TransactionDeposit{
			ReferenceID: list[0].ReferenceID,
			ID:          list[0].ID,
			Amount:      list[0].Amount,
			Status:      list[0].Status,
			DepositBy:   list[0].CreatedBy,
			DepositAt:   list[0].CreatedAt,
		}
	} else {
		return nil, models.ErrDisabled
	}

	return res, nil
}

func (m *mysqlWalletRepository) FetchTransactionWithdraw(ctx context.Context, id int64) (*models.TransactionWithdraw, error) {
	query := `SELECT reference_id, wallet_id, amount, status, created_by, created_at
			  FROM transaction WHERE id = ?`

	list, err := m.fetchTransaction(ctx, query, id)
	if err != nil {
		return nil, err
	}

	res := &models.TransactionWithdraw{}
	if len(list) > 0 {
		res = &models.TransactionWithdraw{
			ReferenceID: list[0].ReferenceID,
			ID:          list[0].ID,
			Amount:      list[0].Amount,
			Status:      list[0].Status,
			WithdrawnBy: list[0].CreatedBy,
			WithdrawnAt: list[0].CreatedAt,
		}
	} else {
		return nil, models.ErrDisabled
	}

	return res, nil
}

func (m *mysqlWalletRepository) AddWallet(ctx context.Context, req *models.ReqTransaction, id string) (*models.TransactionDeposit, error) {
	wallet, err := m.FetchWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	balance := wallet.Balance + req.Amount
	query := `UPDATE wallet set balance = ? WHERE wallet_id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, balance, id)
	if err != nil {
		return nil, err
	}

	query2 := `INSERT INTO transaction (reference_id, wallet_id, type, amount, status, created_by) VALUES (?,?,0,?,?,?);`

	stmt, err = m.Conn.PrepareContext(ctx, query2)
	if err != nil {
		return nil, err
	}

	rowInsert, err := stmt.ExecContext(ctx, req.ReferenceID, wallet.ID, req.Amount, "success", wallet.OwnedBy)
	if err != nil {
		return nil, err
	}

	lastID, err := rowInsert.LastInsertId()
	if err != nil {
		return nil, err
	}

	res, err := m.FetchTransactionAdd(ctx, lastID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlWalletRepository) WithdrawWallet(ctx context.Context, req *models.ReqTransaction, id string) (*models.TransactionWithdraw, error) {
	wallet, err := m.FetchWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	balance := wallet.Balance - req.Amount
	query := `UPDATE wallet set balance = ? WHERE wallet_id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, balance, id)
	if err != nil {
		return nil, err
	}

	query2 := `INSERT INTO transaction (reference_id, wallet_id, type, amount, status, created_by) VALUES (?,?,1,?,?,?);`

	stmt, err = m.Conn.PrepareContext(ctx, query2)
	if err != nil {
		return nil, err
	}

	status := "success"
	if balance < 0 {
		status = "failed"
	}

	rowInsert, err := stmt.ExecContext(ctx, req.ReferenceID, wallet.ID, req.Amount, status, wallet.OwnedBy)
	if err != nil {
		return nil, err
	}

	if balance < 0 {
		return nil, models.ErrBadParamInput
	}

	lastID, err := rowInsert.LastInsertId()
	if err != nil {
		return nil, err
	}

	res, err := m.FetchTransactionWithdraw(ctx, lastID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlWalletRepository) DisableWallet(ctx context.Context, isDisabled bool, id string) (*models.WalletDisabled, error) {
	query := `UPDATE wallet set status = "disabled" updated_at=? WHERE wallet_id = ? AND status = "enabled"`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rowUpdate, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return nil, err
	}
	affect, err := rowUpdate.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect != 1 {
		return nil, models.ErrDisabled
	}

	res, err := m.FetchDisabledWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlWalletRepository) InitWallet(ctx context.Context, customer_id string) (*models.FetchWallet, error) {
	query2 := `INSERT INTO wallet (wallet_id, owned_by, status, balance) VALUES (?,"william-chandra","enabled",0);`

	stmt, err := m.Conn.PrepareContext(ctx, query2)
	if err != nil {
		return nil, err
	}

	_, err = stmt.ExecContext(ctx, customer_id)
	if err != nil {
		return nil, err
	}

	res, err := m.FetchWallet(ctx, customer_id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlWalletRepository) fetchWallet(ctx context.Context, query string, args ...interface{}) ([]*models.Wallet, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Wallet, 0)
	for rows.Next() {
		t := new(models.Wallet)
		err = rows.Scan(
			&t.ID,
			&t.OwnedBy,
			&t.Status,
			&t.UpdatedAt,
			&t.Balance,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlWalletRepository) fetchTransaction(ctx context.Context, query string, args ...interface{}) ([]*models.Transaction, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Transaction, 0)
	for rows.Next() {
		t := new(models.Transaction)
		err = rows.Scan(
			&t.ReferenceID,
			&t.ID,
			&t.Type,
			&t.Amount,
			&t.Status,
			&t.CreatedBy,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
