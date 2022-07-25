package api

import (
	"github.com/labstack/echo/v4"
)

func InitializeAPI() *echo.Echo {
	e := echo.New()
	return e
}
