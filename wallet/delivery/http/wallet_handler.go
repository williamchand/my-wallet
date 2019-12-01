package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/williamchand/my-wallet/models"
	"github.com/williamchand/my-wallet/wallet"
)

// ResponseError represent the reseponse error struct
type Response struct {
	Status       string      `json:"status"`
	ResponseData interface{} `json:"data"`
}
type ResponseWallet struct {
	Wallet interface{} `json:"wallet"`
}
type ResponseDeposit struct {
	Deposit interface{} `json:"deposit"`
}
type ResponseWithdrawal struct {
	Withdrawal interface{} `json:"withdrawal"`
}
type ResponseError struct {
	Error interface{} `json:"error"`
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
	e.POST("/api/v1/wallet", handler.EnableWallet)
	e.GET("/api/v1/wallet", handler.FetchWallet)
	e.POST("/api/v1/wallet/deposits", handler.AddWallet)
	e.POST("/api/v1/wallet/withdrawals", handler.WithdrawWallet)
	e.PATCH("/api/v1/wallet", handler.DisableWallet)
	e.POST("/api/v1/init", handler.InitWallet)
}

// EnableWallet will enable wallet by given param
func (a *WalletHandler) EnableWallet(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := a.AUsecase.EnableWallet(ctx, authorization)

	if err != nil {
		return c.JSON(getStatusCode(err), Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}

	return c.JSON(http.StatusOK, Response{Status: "success", ResponseData: ResponseWallet{
		Wallet: res,
	}})
}

// FetchWallet will fetch the wallet based on given params
func (a *WalletHandler) FetchWallet(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := a.AUsecase.FetchWallet(ctx, authorization)

	if err != nil {
		return c.JSON(getStatusCode(err), Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}

	return c.JSON(http.StatusOK, Response{Status: "success", ResponseData: ResponseWallet{
		Wallet: res,
	}})
}

// AddWallet will deposit the wallet by given request body
func (a *WalletHandler) AddWallet(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var wallet models.ReqTransaction
	err := c.Bind(&wallet)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}

	if ok, err := isRequestValid(&wallet); !ok {
		return c.JSON(http.StatusBadRequest, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.AUsecase.AddWallet(ctx, &wallet, authorization)

	if err != nil {
		return c.JSON(getStatusCode(err), Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	return c.JSON(http.StatusOK, Response{Status: "success", ResponseData: ResponseDeposit{
		Deposit: res,
	}})
}

// WithdrawWallet will withdraw the wallet by given request body
func (a *WalletHandler) WithdrawWallet(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var wallet models.ReqTransaction
	err := c.Bind(&wallet)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}

	if ok, err := isRequestValid(&wallet); !ok {
		return c.JSON(http.StatusBadRequest, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.AUsecase.WithdrawWallet(ctx, &wallet, authorization)

	if err != nil {
		return c.JSON(getStatusCode(err), Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	return c.JSON(http.StatusOK, Response{Status: "success", ResponseData: ResponseWithdrawal{
		Withdrawal: res,
	}})
}

// DisableWallet will disable wallet by given param
func (a *WalletHandler) DisableWallet(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	isDisabled, err := strconv.ParseBool(c.FormValue("is_disabled"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	if isDisabled != true {
		return c.JSON(http.StatusBadRequest, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.AUsecase.DisableWallet(ctx, isDisabled, authorization)

	if err != nil {
		return c.JSON(getStatusCode(err), Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	return c.JSON(http.StatusOK, Response{Status: "success", ResponseData: ResponseWallet{
		Wallet: res,
	}})
}

// InitWallet will init the wallet by given request body
func (a *WalletHandler) InitWallet(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType != "application/json" {
		return c.JSON(http.StatusUnsupportedMediaType, Response{Status: "fail", ResponseData: ResponseError{
			Error: "Content-Type header is not application/json",
		}})
	}

	if c.Request().Body == nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "fail", ResponseData: ResponseError{
			Error: "Please send Body",
		}})
	}
	customer := new(models.Customer)
	err := json.NewDecoder(c.Request().Body).Decode(&customer)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.AUsecase.InitWallet(ctx, customer.ID)

	if err != nil {
		return c.JSON(getStatusCode(err), Response{Status: "fail", ResponseData: ResponseError{
			Error: err.Error(),
		}})
	}
	return c.JSON(http.StatusOK, Response{Status: "success", ResponseData: ResponseWallet{
		Wallet: res,
	}})
}

func isRequestValid(m *models.ReqTransaction) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
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
	case models.ErrAlreadyEnabled:
		return http.StatusBadRequest
	case models.ErrDisabled:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
