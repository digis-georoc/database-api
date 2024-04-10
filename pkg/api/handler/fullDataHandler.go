// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_IDENTIFIER      = "identifier"
	QP_IDENTIFIER_LIST = "samplingfeatureids"

	TAS_SIO2 = "SIO2"
	TAS_K2O  = "K2O"
	TAS_NA2O = "NA2O"
)

// GetFullDataByID godoc
// @Summary     Retrieve full dataset by samplingfeatureid
// @Description get full dataset by samplingfeatureid
// @Security    ApiKeyAuth
// @Tags        fulldata
// @Accept      json
// @Produce     json
// @Param       samplingfeatureids path     string true "Samplingfeature identifier"
// @Success     200                {object} model.FullData
// @Failure     401                {object} string
// @Failure     404                {object} string
// @Failure     500                {object} string
// @Router      /queries/fulldata/{samplingfeatureid} [get]
func (h *Handler) GetFullDataByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	id, err := strconv.Atoi(c.Param(QP_IDENTIFIER))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not parse identifier")
	}
	identifier := []int{id}
	fullData, err := repository.Query[model.FullData](c.Request().Context(), h.db, sql.FullDataByMultiIdQuery, identifier)
	if err != nil {
		logger.Errorf("Can not retrieve FullDataById: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve full data by id")
	}
	num := len(fullData)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}

	// extract TAS diagram values
	for _, batch := range fullData[0].BatchData {
		tasData, err := getTASData(batch.Results)
		if err != nil {
			batch.TASData = nil
			continue
		}
		batch.TASData = tasData
	}

	return c.JSON(http.StatusOK, fullData[0])
}

// GetFullData godoc
// @Summary     Retrieve full datasets by a list of samplingfeatureids
// @Description get full datasets by a list of samplingfeatureids
// @Security    ApiKeyAuth
// @Tags        fulldata
// @Accept      json
// @Produce     json
// @Param       samplingfeatureids query    string true "List of Samplingfeature identifiers"
// @Success     200                {object} model.FullDataResponse
// @Failure     401                {object} string
// @Failure     404                {object} string
// @Failure     500                {object} string
// @Router      /queries/fulldata [get]
func (h *Handler) GetFullData(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	identifierList := []int{}
	identifiers := c.QueryParam(QP_IDENTIFIER_LIST)
	for _, id := range strings.Split(identifiers, ",") {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		identifierList = append(identifierList, idInt)
	}
	fullData, err := repository.Query[model.FullData](c.Request().Context(), h.db, sql.FullDataByMultiIdQuery, identifierList)
	if err != nil {
		logger.Errorf("Can not retrieve FullDataById: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve full data")
	}

	// extract TAS diagram values
	for _, fd := range fullData {
		for _, batch := range fd.BatchData {
			tasData, err := getTASData(batch.Results)
			if err != nil {
				batch.TASData = nil
				continue
			}
			batch.TASData = tasData
		}
	}

	response := model.FullDataResponse{
		NumItems: len(fullData),
		Data:     fullData,
	}
	return c.JSON(http.StatusOK, response)
}

// factors to recalculate from another unit to WT%
var unitFactors = map[string]float64{
	"PPQ": 1000 * 1000 * 1000 * 10000,
	"PPT": 1000 * 1000 * 10000,
	"PPB": 1000 * 10000,
	"PPM": 10000,
	"WT%": 1,
}

func getTASData(results []*model.Result) (*model.DiagramData, error) {
	data := &model.DiagramData{}
	xSum := 0.0
	ySum := 0.0
	for _, result := range results {
		if result == nil {
			continue
		}
		if *result.ItemName == TAS_SIO2 {
			// extract data for TAS x-value
			data.XAxisLabel = TAS_SIO2
			// recalculate to WT%
			factor, ok := unitFactors[*result.Unit]
			if !ok {
				return nil, fmt.Errorf("Invalid unit: %v", *result.Unit)
			}
			xSum += *result.Value * factor
		} else if *result.ItemName == TAS_K2O || *result.ItemName == TAS_NA2O {
			// extract data for TAS y-value
			data.YAxisLabel = TAS_NA2O + "+" + TAS_K2O
			// recalculate to WT%
			factor, ok := unitFactors[*result.Unit]
			if !ok {
				return nil, fmt.Errorf("Invalid unit: %v", *result.Unit)
			}
			ySum += *result.Value * factor
		}
	}
	data.Values = [][]float64{{xSum, ySum}}
	return data, nil
}
