package api

import (
	"strings"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"

	// import swagger docs
	_ "gitlab.gwdg.de/fe/digis/database-api/docs"
)

// @title       DIGIS Database API
// @version     0.1.0
// @description This is the database api for the new GeoROC datamodel
// @description
// @description Note: Semicolon (;) in queries are not allowed and need to be url-encoded as per this issue: golang.org/issue/25192

// @contact.name  DIGIS Project
// @contact.url   https://www.uni-goettingen.de/de/643369.html
// @contact.email digis-info@uni-goettingen.de

// @license.name Data retrieved is licensed under CC BY-SA 4.0
// @license.url  https://creativecommons.org/licenses/by-sa/4.0/

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       DIGIS_API_ACCESSKEY
// @description                Accesskey based security scheme to secure api group "/queries/*"

// @host     api-test.georoc.eu
// @schemes  https http
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
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "docs") || strings.Contains(c.Request().URL.Path, "ping")
		},
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
	v1.GET("/ping", h.Ping)
	v1.GET("/docs/*", echoSwagger.WrapHandler)

	// accesskey queries
	queries := v1.Group("/queries")
	queries.Use(middleware.GetAccessKeyMiddleware(secStore))
	// authors
	queries.GET("/authors", h.GetAuthors)
	queries.GET("/authors/:personID", h.GetAuthorByID)
	// sites
	queries.GET("/sites", h.GetSites)
	queries.GET("/sites/:samplingfeatureID", h.GetSiteByID)
	queries.GET("/sites/settings", h.GetGeoSettings)
	// citations
	queries.GET("/citations", h.GetCitations)
	queries.GET("/citations/:citationID", h.GetCitationByID)
	// full data
	queries.GET("/fullData/:identifier", h.GetFullDataByID)
	// samples
	queries.GET("/samples", h.GetSamplesByGeoSetting)

	// GeoJSON
	geoData := v1.Group("/geodata")
	geoData.Use(middleware.GetAccessKeyMiddleware(secStore))
	// Sites as GeoJSON
	geoData.GET("/sites", h.GetGeoJSONSites)

	return e
}
