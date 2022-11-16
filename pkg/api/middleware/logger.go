package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	LOGGER_KEY = "custom_logger_key"
)

// APILogger is a wrapper for logrus
type APILogger struct {
	*logrus.Entry
}

// Logger middleware sets a custom logrus logger for each request
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		customLogger := logrus.WithFields(logrus.Fields{
			"RequestID":    c.Response().Header()[echo.HeaderXRequestID],
			"usertracking": c.Request().Header.Get(HEADER_USER_TRACKING),
		})
		c.Set(LOGGER_KEY, APILogger{customLogger})
		return next(c)
	}
}
