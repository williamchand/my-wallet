package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	walletRepo "github.com/williamchand/my-wallet/wallet/repository"
	"github.com/williamchand/my-wallet/models"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockWallets := []models.Wallet{
		models.Wallet{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: models.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		models.Wallet{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: models.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(mockWallets[0].ID, mockWallets[0].Title, mockWallets[0].Content,
			mockWallets[0].Author.ID, mockWallets[0].UpdatedAt, mockWallets[0].CreatedAt).
		AddRow(mockWallets[1].ID, mockWallets[1].Title, mockWallets[1].Content,
			mockWallets[1].Author.ID, mockWallets[1].UpdatedAt, mockWallets[1].CreatedAt)

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM wallet WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := walletRepo.NewMysqlWalletRepository(db)
	cursor := walletRepo.EncodeCursor(mockWallets[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM wallet WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := walletRepo.NewMysqlWalletRepository(db)

	num := int64(5)
	anWallet, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anWallet)
}

func TestStore(t *testing.T) {
	now := time.Now()
	ar := &models.Wallet{
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: models.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT  wallet SET title=\\? , content=\\? , author_id=\\?, updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.CreatedAt, ar.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := walletRepo.NewMysqlWalletRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestGetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM wallet WHERE title = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := walletRepo.NewMysqlWalletRepository(db)

	title := "title 1"
	anWallet, err := a.GetByTitle(context.TODO(), title)
	assert.NoError(t, err)
	assert.NotNil(t, anWallet)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM wallet WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := walletRepo.NewMysqlWalletRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	ar := &models.Wallet{
		ID:        12,
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: models.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE wallet set title=\\?, content=\\?, author_id=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := walletRepo.NewMysqlWalletRepository(db)

	err = a.Update(context.TODO(), ar)
	assert.NoError(t, err)
}
