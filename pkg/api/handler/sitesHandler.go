package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
)

const (
	QP_LAT_MIN  = "latMin"
	QP_LAT_MAX  = "latMax"
	QP_LONG_MIN = "lonMin"
	QP_LONG_MAX = "lonMax"
)

// GetSitesByCoords godoc
// @Summary     Retrieve all sites within given coordinates
// @Description Get all sites where the longitude and latitude values are within the given range
// @Description Omit query parameters to retrieve all sites
// @Security    ApiKeyAuth
// @Tags        sites
// @Accept      json
// @Produce     json
// @Param       latMin query     number false "Minimum (inclusive) latitude"
// @Param       latMax query     number false "Maximum (inclusive) latitude"
// @Param       lonMin query     number false "Minimum (inclusive) longitude"
// @Param       lonMax query     number false "Maximum (inclusive) longitude"
// @Success     200 {array}  model.Site
// @Failure     401 {object} string
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /queries/sites [get]
func (h *Handler) GetSitesByCoords(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	latMin := c.QueryParam(QP_LAT_MIN)
	latMax := c.QueryParam(QP_LAT_MAX)
	longMin := c.QueryParam(QP_LONG_MIN)
	longMax := c.QueryParam(QP_LONG_MAX)
	var sites []model.Site
	var err error
	if latMin != "" && latMax != "" && longMin != "" && longMax != "" {
		sites, err = h.db.GetSitesByCoords(latMin, latMax, longMin, longMax)
	} else {
		sites, err = h.db.GetSites()
	}
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
