package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
)

func InitializeAPI(h *handler.Handler, config *middleware.KeycloakConfig) *echo.Echo {
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
	v1.GET("/ping", func(c echo.Context) error { return c.JSON(http.StatusOK, "pong") })
	v1.POST("/login", h.KeycloakLogin)

	// keycloak secured
	secured := v1.Group("/secured")
	secured.Use(middleware.GetAcademicCloudAuthMW(config))
	secured.GET("/authors/:lastName", h.GetAuthors)

	return e
}
