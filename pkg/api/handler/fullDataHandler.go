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
	QP_IDENTIFIER      = "identifier"
	QP_IDENTIFIER_LIST = "samplingfeatureids"
)

// GetFullDataByID godoc
// @Summary     Retrieve full dataset by samplingfeatureid
// @Description get full dataset by samplingfeatureid
// @Security    ApiKeyAuth
// @Tags        fulldata
// @Accept      json
// @Produce     json
// @Param       samplingfeatureids path     string true "Samplingfeature identifier"
// @Success     200               {array}  model.FullData
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /queries/fulldata/{samplingfeatureid} [get]
func (h *Handler) GetFullDataByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	fullData := []model.FullData{}
	err := h.db.Query(sql.FullDataByIdQuery, &fullData, c.Param(QP_IDENTIFIER))
	if err != nil {
		logger.Errorf("Can not retrieve FullDataById: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve full data by id")
	}
	num := len(fullData)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{num, fullData}
	return c.JSON(http.StatusOK, response)
}

// GetFullData godoc
// @Summary     Retrieve full datasets by a list of samplingfeatureids
// @Description get full datasets by a list of samplingfeatureids
// @Security    ApiKeyAuth
// @Tags        fulldata
// @Accept      json
// @Produce     json
// @Param       samplingfeatureids query     string true "List of Samplingfeature identifiers"
// @Success     200               {array}  model.FullData
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /queries/fulldata [get]
func (h *Handler) GetFullData(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	fullData := []model.FullData{}
	err := h.db.Query(sql.FullDataByMultiIdQuery, &fullData, fmt.Sprintf("(%s)", c.QueryParam(QP_IDENTIFIER)))
	if err != nil {
		logger.Errorf("Can not retrieve FullDataById: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve full data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(fullData), fullData}
	return c.JSON(http.StatusOK, response)
}
