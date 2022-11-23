package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_LAT_MIN  = "latMin"
	QP_LAT_MAX  = "latMax"
	QP_LONG_MIN = "lonMin"
	QP_LONG_MAX = "lonMax"
)

// GetSites godoc
// @Summary     Retrieve all sites
// @Description Get all sites
// @Security    ApiKeyAuth
// @Tags        sites
// @Accept      json
// @Produce     json
// @Success     200 {array}  model.Site
// @Failure     401 {object} string
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /queries/sites [get]
func (h *Handler) GetSites(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	sites := []model.Site{}
	err := h.db.Query(sql.SitesQuery, &sites)
	if err != nil {
		logger.Errorf("Can not GetSites: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve site data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(sites), sites}
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
// @Success     200 {array}  model.Site
// @Failure     401 {object} string
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /queries/sites/{samplingfeatureID} [get]
func (h *Handler) GetSiteByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	sites := []model.Site{}
	err := h.db.Query(sql.SiteByIDQuery, &sites, c.Param("samplingfeatureID"))
	if err != nil {
		logger.Errorf("Can not GetSiteByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve site data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(sites), sites}
	return c.JSON(http.StatusOK, response)
}

// GetGeoSettings godoc
// @Summary     Retrieve all geological settings
// @Description Get all geological settings
// @Security    ApiKeyAuth
// @Tags        sites
// @Accept      json
// @Produce     json
// @Success     200 {array}  model.Site
// @Failure     401 {object} string
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /queries/sites/settings [get]
func (h *Handler) GetGeoSettings(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	sites := []model.Site{}
	err := h.db.Query(sql.GeoSettingsQuery, &sites)
	if err != nil {
		logger.Errorf("Can not GetGeoSettings: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geological settings data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(sites), sites}
	return c.JSON(http.StatusOK, response)
}
