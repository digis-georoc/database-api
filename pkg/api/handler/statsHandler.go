package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

// GetCitations godoc
// @Summary     Retrieve data statistics
// @Description get statistics
// @Security    ApiKeyAuth
// @Tags        stats
// @Accept      json
// @Produce     json
// @Success     200    {object} model.Statistics
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/statistics [get]
func (h *Handler) GetStatistics(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	stats := []model.Statistics{}
	query := sql.NewQuery(sql.CountCitationsQuery)
	err := h.db.Query(query.GetQueryString(), &stats)
	if err != nil {
		logger.Errorf("Can not GetStatistics Citations: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	query = sql.NewQuery(sql.CountAnalysesQuery)
	err = h.db.Query(query.GetQueryString(), &stats)
	if err != nil {
		logger.Errorf("Can not GetStatistics Analyses: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	query = sql.NewQuery(sql.CountSamplesQuery)
	err = h.db.Query(query.GetQueryString(), &stats)
	if err != nil {
		logger.Errorf("Can not GetStatistics Samples: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	query = sql.NewQuery(sql.CountResultsQuery)
	err = h.db.Query(query.GetQueryString(), &stats)
	if err != nil {
		logger.Errorf("Can not GetStatistics Results: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	return c.JSON(http.StatusOK, stats)
}
