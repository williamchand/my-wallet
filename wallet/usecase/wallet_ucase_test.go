package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/williamchand/my-wallet/wallet/mocks"
	ucase "github.com/williamchand/my-wallet/wallet/usecase"
	_authorMock "github.com/williamchand/my-wallet/author/mocks"
	"github.com/williamchand/my-wallet/models"
)

func TestFetch(t *testing.T) {
	mockWalletRepo := new(mocks.Repository)
	mockWallet := &models.Wallet{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]*models.Wallet, 0)
	mockListArtilce = append(mockListArtilce, mockWallet)

	t.Run("success", func(t *testing.T) {
		mockWalletRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		mockAuthor := &models.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(_authorMock.Repository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockWalletRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockWalletRepo := new(mocks.Repository)
	mockWallet := models.Wallet{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := &models.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockWalletRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockWallet, nil).Once()
		mockAuthorrepo := new(_authorMock.Repository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockWallet.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockWalletRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected")).Once()

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockWallet.ID)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockWalletRepo := new(mocks.Repository)
	mockWallet := models.Wallet{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockWallet := mockWallet
		tempMockWallet.ID = 0
		mockWalletRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(nil, models.ErrNotFound).Once()
		mockWalletRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.Wallet")).Return(nil).Once()

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockWallet)

		assert.NoError(t, err)
		assert.Equal(t, mockWallet.Title, tempMockWallet.Title)
		mockWalletRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingWallet := mockWallet
		mockWalletRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(&existingWallet, nil).Once()
		mockAuthor := &models.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(_authorMock.Repository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &mockWallet)

		assert.Error(t, err)
		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockWalletRepo := new(mocks.Repository)
	mockWallet := models.Wallet{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockWalletRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockWallet, nil).Once()

		mockWalletRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockWallet.ID)

		assert.NoError(t, err)
		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("wallet-is-not-exist", func(t *testing.T) {
		mockWalletRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, nil).Once()

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockWallet.ID)

		assert.Error(t, err)
		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockWalletRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected Error")).Once()

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockWallet.ID)

		assert.Error(t, err)
		mockWalletRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockWalletRepo := new(mocks.Repository)
	mockWallet := models.Wallet{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockWalletRepo.On("Update", mock.Anything, &mockWallet).Once().Return(nil)

		mockAuthorrepo := new(_authorMock.Repository)
		u := ucase.NewWalletUsecase(mockWalletRepo, mockAuthorrepo, time.Second*2)

		err := u.Update(context.TODO(), &mockWallet)
		assert.NoError(t, err)
		mockWalletRepo.AssertExpectations(t)
	})
}
