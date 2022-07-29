package api

import (
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/src/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/src/middleware"
)

func InitializeAPI(h *handler.Handler) *echo.Echo {
	e := echo.New()
	log := logrus.New()
	e.Use(emw.Recover())
	e.Use(emw.RequestID())
	e.Use(middleware.Logger)
	e.Use(emw.RequestLoggerWithConfig(emw.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogMethod:    true,
		LogValuesFunc: func(c echo.Context, values emw.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"method":    values.Method,
				"URI":       values.URI,
				"status":    values.Status,
				"error":     values.Error,
				"timestamp": values.StartTime,
				"requestID": values.RequestID,
			}).Info("request")

			return nil
		},
	}))

	// api/v1
	v1 := e.Group("/api/v1")
	v1.GET("/authors/:lastName", h.GetAuthors)

	return e
}
