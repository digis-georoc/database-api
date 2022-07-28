package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	REQUEST_ID_HEADER = "X-Request-ID"
	LOGGER_KEY        = "custom_logger_key"
)

// Logger middleware sets a custom logrus logger for each request
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		customLogger := logrus.WithFields(logrus.Fields{
			"requestID": c.Request().Header[REQUEST_ID_HEADER],
		})
		c.Set(LOGGER_KEY, customLogger)
		return next(c)
	}
}
