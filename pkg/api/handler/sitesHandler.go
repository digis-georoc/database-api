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

const (
	QP_SAMPFEATUREID = "samplingfeatureID"
)

// GetSites godoc
// @Summary     Retrieve all sites
// @Description Get all sites
// @Security    ApiKeyAuth
// @Tags        sites
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.SiteResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/sites [get]
func (h *Handler) GetSites(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	query := sql.NewQuery(sql.SitesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	sites, err := repository.Query[model.Site](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetSites: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve site data")
	}
	response := model.SiteResponse{
		NumItems: len(sites),
		Data:     sites,
	}
	return c.JSON(http.StatusOK, response)
}

// GetSiteByID godoc
// @Summary     Retrieve sites by samplingfeatureID
// @Description Get sites by samplingfeatureID
// @Security    ApiKeyAuth
// @Tags        sites
// @Accept      json
// @Produce     json
// @Param       samplingfeatureID path     string true "samplingfeatureID"
// @Success     200               {object} model.Site
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /queries/sites/{samplingfeatureID} [get]
func (h *Handler) GetSiteByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	sites, err := repository.Query[model.Site](c.Request().Context(), h.db, sql.SiteByIDQuery, c.Param(QP_SAMPFEATUREID))
	if err != nil {
		logger.Errorf("Can not GetSiteByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve site data")
	}
	num := len(sites)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}
	return c.JSON(http.StatusOK, sites[0])
}

// GetGeoSettings godoc
// @Summary     Retrieve all geological settings
// @Description Get all geological settings
// @Security    ApiKeyAuth
// @Tags        sites
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.GeologicalSettingResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     500    {object} string
// @Router      /queries/sites/settings [get]
func (h *Handler) GetGeoSettings(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	query := sql.NewQuery(sql.GeoSettingsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	sites, err := repository.Query[model.GeologicalSetting](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetGeoSettings: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geological settings data")
	}
	response := model.GeologicalSettingResponse{
		NumItems: len(sites),
		Data:     sites,
	}
	return c.JSON(http.StatusOK, response)
}
