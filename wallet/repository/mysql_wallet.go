package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/williamchand/my-wallet/wallet"
	"github.com/williamchand/my-wallet/models"
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

func (m *mysqlWalletRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Wallet, error) {
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
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = models.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlWalletRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Wallet, string, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM wallet WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, err
}
func (m *mysqlWalletRepository) GetByID(ctx context.Context, id int64) (res *models.Wallet, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM wallet WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *mysqlWalletRepository) GetByTitle(ctx context.Context, title string) (res *models.Wallet, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM wallet WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}
	return
}

func (m *mysqlWalletRepository) Store(ctx context.Context, a *models.Wallet) error {
	query := `INSERT  wallet SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func (m *mysqlWalletRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM wallet WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {

		return err
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		return err
	}

	return nil
}
func (m *mysqlWalletRepository) Update(ctx context.Context, ar *models.Wallet) error {
	query := `UPDATE wallet set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
