// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package api

import (
	"net/http"
	"strings"
	"time"

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
// @version     0.5.1
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
// @name                       DIGIS-API-ACCESSKEY
// @description                Accesskey based security scheme to secure api groups "/queries/*" and "/geodata/*"

// @host     api-test.georoc.eu
// @schemes  https http
// @BasePath /api/v1
func InitializeAPI(h *handler.Handler, secStore secretstore.SecretStore) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	log := logrus.New()
	e.Use(emw.Recover())
	e.Use(emw.RequestID())
	e.Use(emw.CORSWithConfig(
		emw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead},
			AllowHeaders: []string{"*"},
		},
	))
	e.Use(middleware.GetUserTrackMiddleware())
	e.Use(middleware.Logger)
	e.Use(emw.RequestLoggerWithConfig(emw.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogMethod:    true,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "v1/docs") || strings.Contains(c.Request().URL.Path, "v1/ping") || strings.Contains(c.Request().URL.Path, "v1/alive")
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
				"duration":  time.Since(values.StartTime),
			}).Info("request")

			return nil
		},
	}))

	// api/v1
	v1 := e.Group("/api/v1")
	v1.GET("/ping", h.Ping)
	v1.GET("/alive", h.Alive)
	v1.GET("/version", h.Version)
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
	// locations
	queries.GET("/locations/l1", h.GetLocationsL1)
	queries.GET("/locations/l2", h.GetLocationsL2)
	queries.GET("/locations/l3", h.GetLocationsL3)
	// citations
	queries.GET("/citations", h.GetCitations)
	queries.GET("/citations/:citationID", h.GetCitationByID)
	// full data
	queries.GET("/fulldata/:identifier", h.GetFullDataByID)
	queries.GET("/fulldata", h.GetFullData)
	// samples
	queries.GET("/samples", h.GetSamplesFiltered)
	queries.GET("/samples/:samplingfeatureID", h.GetSampleByID)
	queries.GET("/samples/random", h.GetRandomSamples)
	queries.GET("/samples/specimentypes", h.GetSpecimenTypes)
	queries.GET("/samples/rockclasses", h.GetRockClasses)
	queries.GET("/samples/rocktypes", h.GetRockTypes)
	queries.GET("/samples/minerals", h.GetMinerals)
	queries.GET("/samples/materials", h.GetMaterials)
	queries.GET("/samples/inclusiontypes", h.GetInclusionTypes)
	queries.GET("/samples/inclusionmaterials", h.GetInclusionMaterials)
	queries.GET("/samples/hostmaterials", h.GetHostMaterials)
	queries.GET("/samples/samplingtechniques", h.GetSamplingTechniques)
	queries.GET("/samples/geoages", h.GetGeoAges)
	queries.GET("/samples/geoageprefixes", h.GetGeoAgePrefixes)
	queries.GET("/samples/organizationnames", h.GetOrganizationNames)
	// results
	queries.GET("/results/elements", h.GetElements)
	queries.GET("/results/elementtypes", h.GetElementTypes)
	// statistics
	queries.GET("/statistics", h.GetStatistics)
	// GeoJSON
	geoData := v1.Group("/geodata")
	geoData.Use(middleware.GetAccessKeyMiddleware(secStore))
	// Sites as GeoJSON
	geoData.GET("/sites", h.GetGeoJSONSites)
	geoData.GET("/samplesclustered", h.GetSamplesFilteredClustered)

	return e
}
