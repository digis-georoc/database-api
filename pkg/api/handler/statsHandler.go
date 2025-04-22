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

// GetCitations godoc
//	@Summary		Retrieve data statistics
//	@Description	get statistics
//	@Security		ApiKeyAuth
//	@Tags			stats
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Statistics
//	@Failure		401	{object}	string
//	@Failure		404	{object}	string
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/queries/statistics [get]
func (h *Handler) GetStatistics(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	statisticsResponse := model.Statistics{}
	query := sql.NewQuery(sql.CountCitationsQuery)
	citations, err := repository.Query[struct{ NumCitations int }](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil || len(citations) == 0 {
		logger.Errorf("Can not GetStatistics Citations: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	statisticsResponse.NumCitations = citations[0].NumCitations
	query = sql.NewQuery(sql.CountAnalysesQuery)
	analyses, err := repository.Query[struct{ NumAnalyses int }](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil || len(analyses) == 0 {
		logger.Errorf("Can not GetStatistics Analyses: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	statisticsResponse.NumAnalyses = analyses[0].NumAnalyses
	query = sql.NewQuery(sql.CountSamplesQuery)
	samples, err := repository.Query[struct{ NumSamples int }](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil || len(samples) == 0 {
		logger.Errorf("Can not GetStatistics Samples: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	statisticsResponse.NumSamples = samples[0].NumSamples
	query = sql.NewQuery(sql.CountResultsQuery)
	results, err := repository.Query[struct{ NumResults int }](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil || len(results) == 0 {
		logger.Errorf("Can not GetStatistics Results: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	statisticsResponse.NumResults = results[0].NumResults
	query = sql.NewQuery(sql.LatestTimestampQuery)
	dateResult, err := repository.Query[struct{ LatestDate string }](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil || len(dateResult) == 0 {
		logger.Errorf("Can not GetStatistics Results: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve statistics data")
	}
	statisticsResponse.LatestDate = dateResult[0].LatestDate
	return c.JSON(http.StatusOK, statisticsResponse)
}
