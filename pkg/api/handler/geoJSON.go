// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

// GetGeoJSONSites godoc
//	@Summary		Retrieve site data as GeoJSON
//	@Description	get site data in GeoJSON format
//	@Security		ApiKeyAuth
//	@Tags			geodata
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"limit"
//	@Param			offset	query		int	false	"offset"
//	@Success		200		{object}	model.GeoJSONFeatureCollection
//	@Failure		401		{object}	string
//	@Failure		404		{object}	string
//	@Failure		422		{object}	string
//	@Failure		500		{object}	string
//	@Router			/geodata/sites [get]
func (h *Handler) GetGeoJSONSites(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	query := sql.NewQuery(sql.GeoJSONQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	sites, err := repository.Query[model.GeoJSONSite](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GeoJSONSites: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geoJSON site data")
	}
	response := model.GeoJSONFeatureCollection{
		Type:           model.GEOJSONTYPE_FEATURECOLLECTION,
		NumberMatched:  len(sites),
		NumberReturned: len(sites),
		Features:       buildFeatures(sites),
	}
	return c.JSON(http.StatusOK, response)
}

// buildFeatures takes a list of query results and populates GeoJSONFeatures with the data
func buildFeatures(sites []model.GeoJSONSite) []model.GeoJSONFeature {
	featureList := []model.GeoJSONFeature{}
	for i, result := range sites {
		feature := model.GeoJSONFeature{
			Type: model.GEOJSONTYPE_FEATURE,
			ID:   fmt.Sprintf("%d", i),
			Geometry: model.Geometry{
				Type:        model.GEOJSON_GEOMETRY_POINT,
				Coordinates: []interface{}{result.Longitude, result.Latitude},
			},
			Properties: map[string]interface{}{
				"latitude":               result.Latitude,
				"longitude":              result.Longitude,
				"num_samplingfeatureids": result.NumSamplingFeatureIDs,
				"samplingfeatureids":     result.SamplingFeatureIDs,
				"setting":                result.Setting,
				"loc1":                   result.Loc1,
				"loc2":                   result.Loc2,
				"loc3":                   result.Loc3,
				"land_or_sea":            result.LandOrSea,
			},
		}
		featureList = append(featureList, feature)
	}
	return featureList
}
