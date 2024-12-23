package resp

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/usual2970/certimate/internal/domain"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Succ(e echo.Context, data interface{}) error {
	rs := &Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	return e.JSON(http.StatusOK, rs)
}

func Err(e echo.Context, err error) error {
	xerr, ok := err.(*domain.XError)
	code := 100
	if ok {
		code = xerr.GetCode()
	}

	rs := &Response{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	}
	return e.JSON(http.StatusOK, rs)
}
