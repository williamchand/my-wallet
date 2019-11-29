package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/williamchand/my-wallet/wallet"
	"github.com/williamchand/my-wallet/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// WalletHandler  represent the httphandler for wallet
type WalletHandler struct {
	AUsecase wallet.Usecase
}

// NewWalletHandler will initialize the wallets/ resources endpoint
func NewWalletHandler(e *echo.Echo, us wallet.Usecase) {
	handler := &WalletHandler{
		AUsecase: us,
	}
	e.GET("/wallets", handler.FetchWallet)
	e.POST("/wallets", handler.Store)
	e.GET("/wallets/:id", handler.GetByID)
	e.DELETE("/wallets/:id", handler.Delete)
}

// FetchWallet will fetch the wallet based on given params
func (a *WalletHandler) FetchWallet(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get wallet by given id
func (a *WalletHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.AUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *models.Wallet) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the wallet by given request body
func (a *WalletHandler) Store(c echo.Context) error {
	var wallet models.Wallet
	err := c.Bind(&wallet)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&wallet); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AUsecase.Store(ctx, &wallet)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, wallet)
}

// Delete will delete wallet by given param
func (a *WalletHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	}
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
