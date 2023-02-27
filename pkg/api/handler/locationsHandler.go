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
	PARAM_LOC_LEVEL_1 = "locationl1"
	PARAM_LOC_LEVEL_2 = "locationl2"
)

// GetLocationsL1 godoc
// @Summary     Retrieve locations of first level
// @Description get top level locations
// @Security    ApiKeyAuth
// @Tags        locations
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.Location
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/locations/l1 [get]
func (h *Handler) GetLocationsL1(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	locations := []model.Location{}
	query := sql.NewQuery(sql.FirstLevelLocationNamesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &locations)
	if err != nil {
		logger.Errorf("Can not GetLocationsL1: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve first level locations")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(locations), locations}
	return c.JSON(http.StatusOK, response)
}

// GetLocationsL2 godoc
// @Summary     Retrieve locations of second level
// @Description get second level locations
// @Security    ApiKeyAuth
// @Tags        locations
// @Accept      json
// @Produce     json
// @Param       limit      query    int    false "limit"
// @Param       offset     query    int    false "offset"
// @Param       locationl1 query    string true  "Locationname Level 1"
// @Success     200        {array}  model.Location
// @Failure     401        {object} string
// @Failure     404        {object} string
// @Failure     422        {object} string
// @Failure     500        {object} string
// @Router      /queries/locations/l2 [get]
func (h *Handler) GetLocationsL2(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	locations := []model.Location{}
	query := sql.NewQuery(sql.SecondLevelLocationNamesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	locationNameL1 := c.QueryParam(PARAM_LOC_LEVEL_1)
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &locations, locationNameL1)
	if err != nil {
		logger.Errorf("Can not GetLocationsL2: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve second level locations")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(locations), locations}
	return c.JSON(http.StatusOK, response)
}

// GetLocationsL3 godoc
// @Summary     Retrieve locations of third level
// @Description get third level locations
// @Security    ApiKeyAuth
// @Tags        locations
// @Accept      json
// @Produce     json
// @Param       limit      query    int    false "limit"
// @Param       offset     query    int    false "offset"
// @Param       locationl1 query    string true  "Locationname Level 1"
// @Param       locationl2 query    string true  "Locationname Level 2"
// @Success     200        {array}  model.Location
// @Failure     401        {object} string
// @Failure     404        {object} string
// @Failure     422        {object} string
// @Failure     500        {object} string
// @Router      /queries/locations/l3 [get]
func (h *Handler) GetLocationsL3(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	locations := []model.Location{}
	query := sql.NewQuery(sql.ThirdLevelLocationNamesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	locationNameL1 := c.QueryParam(PARAM_LOC_LEVEL_1)
	locationNameL2 := c.QueryParam(PARAM_LOC_LEVEL_2)
	err = h.db.Query(query.GetQueryString(), &locations, locationNameL1, locationNameL2)
	if err != nil {
		logger.Errorf("Can not GetLocationsL3: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve third level locations")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(locations), locations}
	return c.JSON(http.StatusOK, response)
}
