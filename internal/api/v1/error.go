package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
)

func errorResponse(c echo.Context, status int, msg string) {
	err := errors.New(msg)
	var HTTPError *echo.HTTPError
	ok := errors.As(err, &HTTPError)
	if !ok {
		r := echo.NewHTTPError(status, err.Error())
		_ = c.JSON(status, r)
	}
	c.Error(errors.New("internal server error"))
}
