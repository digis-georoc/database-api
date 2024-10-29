// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/download"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/geometry"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	PARAM_FORMAT   = "format"
	QP_SAMPLE_LIST = "sampleids"

	CONCURRENT_TASKS = 10
	BATCH_SIZE       = 100
)

// GetDataDownloadByIDs godoc
// @Summary     Retrieve download data for the given sample IDs
// @Description get the full data for a list of sample IDs as a csv or xlsx file
// @Security    ApiKeyAuth
// @Tags        download
// @Accept      json
// @Produce     plain
// @Param       sampleids query    string true  "List of Sample identifiers"
// @Param       format    query    string true  "Desired output format: csv (default) or xlsx"
// @Response    102       {header} -      Sends back Headers while progressing the request
// @Success     200       {file}   file
// @Failure     401       {object} string
// @Failure     404       {object} string
// @Failure     500       {object} string
// @Router      /download/sampleid [get]
func (h *Handler) GetDataDownloadByIDs(c echo.Context) error {
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
	targetFormat := c.QueryParam(PARAM_FORMAT)
	if targetFormat == "" {
		targetFormat = download.CSV
	}
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
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().WriteHeader(http.StatusProcessing)
	// flush headers
	c.Response().Flush()
	// query the full data for each given identifier
	samples, err := repository.Query[model.FullData](c.Request().Context(), h.db, sql.FullDataByMultiIdQuery, identifierList)
	if err != nil {
		logger.Errorf("Can not GetDataDownloadByIDs: %v", err)
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
	c.Response().Header().Set("Content-Length", strconv.FormatInt(stats.Size(), 10))

	return c.File(fileName)
}

// GetDataDownloadByFilter godoc
// @Summary     Retrieve download data for the given filters
// @Description get the full data for the given filters as a csv or xlsx file
// @Description Filter DSL syntax:
// @Description FIELD=OPERATOR:VALUE
// @Description where FIELD is one of the accepted query params; OPERATOR is one of "lt" (<), "gt" (>), "eq" (=), "in" (IN), "lk" (LIKE), "btw" (BETWEEN)
// @Description and VALUE is an unquoted string, integer or decimal
// @Description Multiple VALUEs for an "in"-filter must be comma-separated and will be interpreted as a discunctive filter.
// @Description The OPERATORs "lt", "gt" and "btw" are only applicable to numerical values.
// @Description The OPERATOR "lk" is only applicable to string values and supports wildcards `*`(0 or more chars) and `?`(one char).
// @Description The OPERATOR "btw" accepts two comma-separated values as the inclusive lower and upper bound. Missing values are assumed as 0 and 9999999 respectively.
// @Description If no OPERATOR is specified, "eq" is assumed as the default OPERATOR.
// @Description The filters are evaluated conjunctively.
// @Description Note that applying more filters can slow down the query as more tables have to be considered in the evaluation.
// @Security    ApiKeyAuth
// @Tags        download
// @Accept      json
// @Produce     plain
// @Param       format            query    string true  "Desired output format: csv (default) or xlsx"
// @Param       limit             query    int    false "limit"
// @Param       offset            query    int    false "offset"
// @Param       setting           query    string false "tectonic setting - see /queries/sites/settings (supports Filter DSL)"
// @Param       location1         query    string false "location level 1 - see /queries/locations/l1 (supports Filter DSL)"
// @Param       location2         query    string false "location level 2 - see /queries/locations/l2 (supports Filter DSL)"
// @Param       location3         query    string false "location level 3 - see /queries/locations/l3 (supports Filter DSL)"
// @Param       latitude          query    string false "latitude (supports Filter DSL)"
// @Param       longitude         query    string false "longitude (supports Filter DSL)"
// @Param       rocktype          query    string false "rock type - see /queries/samples/rocktypes (supports Filter DSL)"
// @Param       rockclassID       query    int    false "taxonomic classifier ID - see /queries/samples/rockclasses value (supports Filter DSL)"
// @Param       mineral           query    string false "mineral - see /queries/samples/minerals (supports Filter DSL)"
// @Param       material          query    string false "material - see /queries/samples/materials (supports Filter DSL)"
// @Param       inclusiontype     query    string false "inclusion type - see /queries/samples/inclusiontypes (supports Filter DSL)"
// @Param       hostmaterial      query    string false "host material - see /queries/samples/hostmaterials (supports Filter DSL)"
// @Param       inclusionmaterial query    string false "inclusion material - see /queries/samples/inclusionmaterials (supports Filter DSL)"
// @Param       sampletech        query    string false "sampling technique - see /queries/samples/samplingtechniques (supports Filter DSL)"
// @Param       rimorcore         query    string false "rim or core - R = Rim, C = Core, I = Intermediate (supports Filter DSL)"
// @Param       chemistry         query    string false "chemical filter using the form `(TYPE,ELEMENT,MIN,MAX),...` where the filter tuples are evaluated conjunctively"
// @Param       title             query    string false "title of publication (supports Filter DSL)"
// @Param       publicationyear   query    string false "publication year (supports Filter DSL)"
// @Param       doi               query    string false "DOI (supports Filter DSL)"
// @Param       firstname         query    string false "Author first name (supports Filter DSL)"
// @Param       lastname          query    string false "Author last name (supports Filter DSL)"
// @Param       agemin            query    string false "Specimen age min (supports Filter DSL)"
// @Param       agemax            query    string false "Specimen age max (supports Filter DSL)"
// @Param       geoage            query    string false "Specimen geological age - see /queries/samples/geoages (supports Filter DSL)"
// @Param       geoageprefix      query    string false "Specimen geological age prefix - see /queries/samples/geoageprefixes (supports Filter DSL)"
// @Param       lab               query    string false "Laboratory name - see /queries/samples/organizationnames (supports Filter DSL)"
// @Param       polygon           query    string false "Coordinate-Polygon formatted as 2-dimensional json array: [[LONG,LAT],[2.4,6.3]]"
// @Param       addcoordinates    query    bool   false "Add coordinates to each sample"
// @Response    102               {header} -      Sends back Headers while progressing the request
// @Success     200               {file}   file
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /download/filtered [get]
func (h *Handler) GetDataDownloadByFilter(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	// get polygon filter
	coordData := map[string]interface{}{}
	polygonString, _, err := parseParam(c.QueryParam(QP_POLY))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not parse polygon")
	}
	if polygonString != "" {
		polygon, err := geometry.ParsePointArray(polygonString)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Can not parse polygon")
		}
		boundaryPoly, translationFactorPoly, err := geometry.CalcTranslation(polygon)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Can not calculate polygon translation - polygon too big")
		}
		coordData[KEY_POLYGON] = polygon
		coordData[KEY_TRANSLATION_FACTOR_POLY] = translationFactorPoly
		coordData[KEY_BOUNDARY_POLY] = boundaryPoly
	}
	query, err := buildSampleFilterQuery(c, coordData, nil)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	// wrap in rowcount sql
	query.WrapInSQL("select *, count(*) over () as totalCount from (", ") q")

	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)

	// prepare response and start the query
	targetFormat := c.QueryParam(PARAM_FORMAT)
	if targetFormat == "" {
		targetFormat = download.CSV
	}
	fileName := fmt.Sprintf("GEOROC_data_download_%s_%s.%s", c.Request().Header.Get("requestID"), time.Now().Format("20060102"), targetFormat)
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().WriteHeader(http.StatusProcessing)
	// flush headers
	c.Response().Flush()

	results, err := repository.Query[model.SampleByFilters](c.Request().Context(), h.db, query.GetQueryString(), query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetDataDownloadByFilter: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	identifierList := make([]int, 0, len(results))
	for _, sample := range results {
		identifierList = append(identifierList, sample.SampleID)
	}
	// create temp download file
	f, err := os.Create(fileName)
	defer cleanupDownloadFile(f, logger)
	if err != nil {
		logger.Errorf("Can not create file %s: %v", fileName, err)
		return c.String(http.StatusInternalServerError, "Failed to create download file")
	}
	if len(identifierList) == 0 {
		return c.File(fileName)
	}
	// query the full data for each given identifier concurrently
	resultChan := make(chan model.FullData)
	errChan := make(chan error)
	tasks := make([]int, len(identifierList))
	copy(tasks, identifierList)
	returnCount := 0
	readLock := sync.Mutex{}
	for i := 0; i < CONCURRENT_TASKS; i++ {
		go startWorker(c.Request().Context(), errChan, resultChan, h.db, &readLock, &tasks)
	}
	samples := make([]model.FullData, 0, len(identifierList))
Listener:
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				break Listener
			}
			returnCount++
			samples = append(samples, result)
		case err, ok := <-errChan:
			if !ok {
				break Listener
			}
			logger.Errorf("Can not retrieve FullDataById: %v", err)
			return c.String(http.StatusInternalServerError, "Can not retrieve full data")
		default:
			if returnCount >= len(identifierList) {
				close(resultChan)
				close(errChan)
				break Listener
			}
		}
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
	c.Response().Header().Set("Content-Length", strconv.FormatInt(stats.Size(), 10))
	return c.File(fileName)
}

func startWorker(context context.Context, errChan chan error, resultChan chan model.FullData, db repository.PostgresConnector, readLock *sync.Mutex, tasks *[]int) {
	for tasks != nil && len(*tasks) > 0 {
		// securely pop the first batch of items from the task list
		readLock.Lock()
		batch := BATCH_SIZE
		if len(*tasks) < BATCH_SIZE {
			batch = len(*tasks)
		}
		identifiers := (*tasks)[:batch]
		*tasks = (*tasks)[batch:]
		readLock.Unlock()
		results, err := repository.Query[model.FullData](context, db, sql.FullDataByMultiIdQuery, identifiers)
		if err != nil {
			errChan <- err
			return
		}
		for _, result := range results {
			resultChan <- result
		}
	}
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
