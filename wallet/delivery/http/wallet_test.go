package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	walletHttp "github.com/williamchand/my-wallet/wallet/delivery/http"
	"github.com/williamchand/my-wallet/wallet/mocks"
	"github.com/williamchand/my-wallet/models"
)

func TestFetch(t *testing.T) {
	var mockWallet models.Wallet
	err := faker.FakeData(&mockWallet)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	mockListWallet := make([]*models.Wallet, 0)
	mockListWallet = append(mockListWallet, &mockWallet)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListWallet, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/wallet?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := walletHttp.WalletHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchWallet(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/wallet?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := walletHttp.WalletHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchWallet(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockWallet models.Wallet
	err := faker.FakeData(&mockWallet)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockWallet.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(&mockWallet, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/wallet/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("wallet/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := walletHttp.WalletHandler{
		AUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockWallet := models.Wallet{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockWallet := mockWallet
	tempMockWallet.ID = 0
	mockUCase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockWallet)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*models.Wallet")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/wallet", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet")

	handler := walletHttp.WalletHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockWallet models.Wallet
	err := faker.FakeData(&mockWallet)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockWallet.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/wallet/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("wallet/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := walletHttp.WalletHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
