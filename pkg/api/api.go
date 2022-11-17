package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"
)

// @title       DIGIS Database API
// @version     0.1.0
// @description This is the database api for the new GeoROC datamodel

// @contact.name  DIGIS Project
// @contact.url   https://www.uni-goettingen.de/de/643369.html
// @contact.email digis-info@uni-goettingen.de

// @license.name Data retrieved is licensed under CC BY-SA 4.0
// @license.url  https://creativecommons.org/licenses/by-sa/4.0/

// @host     localhost:8081
// @BasePath /api/v1
func InitializeAPI(h *handler.Handler, secStore secretstore.SecretStore) *echo.Echo {
	e := echo.New()
	log := logrus.New()
	e.Use(emw.Recover())
	e.Use(emw.RequestID())
	e.Use(middleware.GetUserTrackMiddleware())
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
				"userTrack": c.Request().Header.Get(middleware.HEADER_USER_TRACKING),
			}).Info("request")

			return nil
		},
	}))

	// api/v1
	v1 := e.Group("/api/v1")
	v1.GET("/ping", func(c echo.Context) error { return c.JSON(http.StatusOK, "pong") })
	v1.GET("/openapi/*", echoSwagger.WrapHandler)

	// accesskey secured
	secured := v1.Group("/secured")
	secured.Use(middleware.GetAccessKeyMiddleware(secStore))
	secured.GET("/authors/:lastName", h.GetAuthorsByLastname)
	secured.GET("/fullData/:identifier", h.GetFullData)
	secured.GET("/sites", h.GetSites)

	return e
}
