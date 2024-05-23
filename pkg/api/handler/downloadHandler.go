// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/download"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	PARAM_FORMAT   = "format"
	QP_SAMPLE_LIST = "sampleids"
)

// GetDataDownload godoc
// @Summary     Retrieve download data for the given sample IDs
// @Description get the full data for a list of sample IDs as a csv or xlsx file
// @Security    ApiKeyAuth
// @Tags        download
// @Accept      json
// @Produce     file
// @Param       sampleids query    string true "List of Sample identifiers"
// @Success     200                {object} file
// @Failure     401                {object} string
// @Failure     404                {object} string
// @Failure     500                {object} string
// @Router      /download/:format [get]
func (h *Handler) GetDataDownload(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	identifierList := []int{}
	identifiers := c.QueryParam(QP_SAMPLE_LIST)
	if identifiers == "" {
		return c.String(http.StatusBadRequest, "missing identifier")
	}
	for _, id := range strings.Split(identifiers, ",") {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Can not parse given identifiers: %s", err.Error()))
		}
		identifierList = append(identifierList, idInt)
	}
	// create temp download file
	targetFormat := c.Param(PARAM_FORMAT)
	fileName := fmt.Sprintf("GEOROC_data_download_%s_%s.%s", c.Request().Header.Get("requestID"), time.Now().Format("20060102"), targetFormat)
	f, err := os.Create(fileName)
	defer cleanupDownloadFile(f, logger)
	if err != nil {
		logger.Errorf("Can not create file %s: %v", fileName, err)
		return c.String(http.StatusInternalServerError, "Failed to create download file")
	}
	if len(identifierList) == 0 {
		return c.File(fileName)
	}
	// query the full data for each given identifier
	samples, err := repository.Query[model.FullData](c.Request().Context(), h.db, sql.FullDataByMultiIdQuery, identifierList)
	if err != nil {
		logger.Errorf("Can not retrieve FullDataById: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve full data")
	}

	data, err := formatData(samples, targetFormat)
	if err != nil {
		logger.Errorf("Can not format given data as %s: %s", targetFormat, err.Error())
		return c.String(http.StatusInternalServerError, "Data formatting failed (supported formats are 'csv' and 'xlsx')")
	}

	// write the formatted data to the download file and set the response headers
	n, err := f.Write(data)
	if err != nil {
		logger.Errorf("Can not write data: stopped at %d with error %v", n, err)
		return c.String(http.StatusInternalServerError, "Failed to write data")
	}
	stats, _ := f.Stat()
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Length", strconv.FormatInt(stats.Size(), 10))

	return c.File(fileName)
}

// formatData takes a list of full sample data
// and formats it according to the current GEOROC output format in the specified data format
func formatData(samples []model.FullData, targetFormat string) ([]byte, error) {
	if targetFormat == download.CSV || targetFormat == download.XLSX {
		formatter := download.GetFormatter(targetFormat)
		return formatter.FormatData(samples)
	}
	return nil, fmt.Errorf("Invalid format '%s': must be one of 'csv' or 'xlsx'", targetFormat)
}

// cleanupDownloadFile deletes the download file and closes it
func cleanupDownloadFile(f *os.File, logger middleware.APILogger) {
	err := f.Close()
	if err != nil {
		logger.Errorf("Can not close download file: %v", err)
	}
	err = os.Remove(f.Name())
	if err != nil {
		logger.Errorf("Can not remove download file: %v", err)
	}
}
