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

	TAS_NO_DATA = -1

	ITEM_GROUP_MJ  = "mj"
	ITEM_GROUP_REE = "ree"
	ITEM_GROUP_TE  = "te"
)

// GetFullDataByID godoc
//	@Summary		Retrieve full dataset by samplingfeatureid
//	@Description	get full dataset by samplingfeatureid
//	@Security		ApiKeyAuth
//	@Tags			fulldata
//	@Accept			json
//	@Produce		json
//	@Param			samplingfeatureids	path		string	true	"Samplingfeature identifier"
//	@Success		200					{object}	model.FullData
//	@Failure		401					{object}	string
//	@Failure		404					{object}	string
//	@Failure		500					{object}	string
//	@Router			/queries/fulldata/{samplingfeatureid} [get]
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
//	@Summary		Retrieve full datasets by a list of samplingfeatureids
//	@Description	get full datasets by a list of samplingfeatureids
//	@Security		ApiKeyAuth
//	@Tags			fulldata
//	@Accept			json
//	@Produce		json
//	@Param			samplingfeatureids	query		string	true	"List of Samplingfeature identifiers"
//	@Success		200					{object}	model.FullDataResponse
//	@Failure		401					{object}	string
//	@Failure		404					{object}	string
//	@Failure		500					{object}	string
//	@Router			/queries/fulldata [get]
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

// priority maps for methods
// major elements
var methodPriosMj = map[string]int{
	"XRF":        10,
	"WET":        9,
	"EMP (EPMA)": 8,
	"AES":        7,
	"AAS":        6,
}

// Rare earth elements
var methodPriosRee = map[string]int{
	"TIMS_ID": 10,
	"ICPMS":   9,
	"AES":     8,
	"SSMS":    7,
	"INAA":    6,
	"XRF":     5,
}

// trace elements:
var methodPriosTe = map[string]int{
	"TIMS_ID": 10,
	"ICPMS":   9,
	"SSMS":    8,
	"XRF":     7,
	"AES":     6,
	"INAA":    5,
	"AAS":     4,
	"SIMS":    3,
	"WET":     2,
}

type TASData struct {
	SIO2       *float64
	NA2O       *float64
	K2O        *float64
	Itemgroups []string
}

func getTASData(results []*model.Result) (*model.DiagramData, error) {
	// aggregate results by method; first method to have all 3 values is put as TAS values
	methodsMap := map[string]TASData{}
	for _, result := range results {
		if result == nil || (*result.ItemName != TAS_SIO2 && *result.ItemName != TAS_NA2O && *result.ItemName != TAS_K2O) {
			continue
		}
		data, ok := methodsMap[*result.Method]
		if !ok {
			data = TASData{}
		}
		if *result.ItemName == TAS_SIO2 {
			// recalculate value to WT%
			factor, ok := unitFactors[*result.Unit]
			if !ok {
				return nil, fmt.Errorf("Invalid unit: %v", *result.Unit)
			}
			value := (*result.Value * factor)
			data.SIO2 = &value
			data.Itemgroups = append(data.Itemgroups, *result.ItemGroup)
		} else if *result.ItemName == TAS_K2O {
			// recalculate value to WT%
			factor, ok := unitFactors[*result.Unit]
			if !ok {
				return nil, fmt.Errorf("Invalid unit: %v", *result.Unit)
			}
			value := (*result.Value * factor)
			data.K2O = &value
			data.Itemgroups = append(data.Itemgroups, *result.ItemGroup)
		} else if *result.ItemName == TAS_NA2O {
			// recalculate value to WT%
			factor, ok := unitFactors[*result.Unit]
			if !ok {
				return nil, fmt.Errorf("Invalid unit: %v", *result.Unit)
			}
			value := (*result.Value * factor)
			data.NA2O = &value
			data.Itemgroups = append(data.Itemgroups, *result.ItemGroup)
		}
		methodsMap[*result.Method] = data
	}
	curMethod := ""
	var prioTASData *TASData = nil
	for method, data := range methodsMap {
		if isTASDataComplete(data) {
			if prioTASData == nil {
				prioTASData = &data
				curMethod = method
				continue
			}
			switch data.Itemgroups[0] {
			case ITEM_GROUP_MJ:
				if prio, ok := methodPriosMj[method]; ok {
					curPrio, ok := methodPriosMj[curMethod]
					if !ok || prio > curPrio {
						// overwrite with higher prio
						prioTASData = &data
						curMethod = method
						continue
					}
				}
			case ITEM_GROUP_REE:
				if prio, ok := methodPriosRee[method]; ok {
					curPrio, ok := methodPriosRee[curMethod]
					if !ok || prio > curPrio {
						// overwrite with higher prio
						prioTASData = &data
						curMethod = method
						continue
					}
				}
			case ITEM_GROUP_TE:
				if prio, ok := methodPriosTe[method]; ok {
					curPrio, ok := methodPriosTe[curMethod]
					if !ok || prio > curPrio {
						// overwrite with higher prio
						prioTASData = &data
						curMethod = method
						continue
					}
				}
			}
		}
	}
	values := [][]float64{}
	if prioTASData != nil && isTASDataComplete(*prioTASData) {
		values = [][]float64{{*prioTASData.SIO2, *prioTASData.K2O + *prioTASData.NA2O}}
	}
	return &model.DiagramData{
		XAxisLabel: TAS_SIO2,
		YAxisLabel: TAS_NA2O + "+" + TAS_K2O,
		Values:     values,
	}, nil
}

func isTASDataComplete(data TASData) bool {
	return data.SIO2 != nil && data.NA2O != nil && data.K2O != nil
}
