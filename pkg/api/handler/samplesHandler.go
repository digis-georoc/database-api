package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_SETTING = "setting"
	QP_LOC1    = "location1"
	QP_LOC2    = "location2"
	QP_LOC3    = "location3"
	QP_LAT     = "latitude"
	QP_LONG    = "longitude"

	QP_ROCKTYPE  = "rocktype"
	QP_ROCKCLASS = "rockclass"
	QP_MINERAL   = "mineral"

	QP_MATERIAL = "material"
	QP_INCTYPE  = "inclusiontype"
	QP_SAMPTECH = "sampletech"

	QP_ELEM     = "element"
	QP_ELEMTYPE = "elementtype"
	QP_VALUE    = "value"

	QP_TITLE        = "title"
	QP_PUBYEAR      = "publicationyear"
	QP_DOI          = "doi"
	QP_AUTHOR_FIRST = "firstname"
	QP_AUTHOR_LAST  = "lastname"

	QP_AGE_MIN        = "agemin"
	QP_AGE_MAX        = "agemax"
	QP_GEO_AGE        = "geoage"
	QP_GEO_AGE_PREFIX = "geoageprefix"

	QP_LAB = "lab"

	QP_POLY = "polygon"

	QP_ADD_COORDINATES = "addcoordinates"

	QP_NUM_CLUSTERS      = "numClusters"
	QP_MAX_DISTANCE      = "maxDistance"
	DEFAULT_NUM_CLUSTERS = "7"
	DEFAULT_MAX_DISTANCE = "100"
)

// GetSampleByID godoc
// @Summary     Retrieve sample by samplingfeatureid
// @Description get sample by samplingfeatureid
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       samplingfeatureID path     string true "Sample ID"
// @Success     200               {array}  model.Sample
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /queries/samples/{samplingfeatureID} [get]
func (h *Handler) GetSampleByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	samples := []model.Sample{}
	query := sql.NewQuery(sql.GetSampleByIDQuery)
	err := h.db.Query(query.GetQueryString(), &samples, c.Param(QP_SAMPFEATUREID))
	if err != nil {
		logger.Errorf("Can not GetSampleByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	num := len(samples)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{num, samples}
	return c.JSON(http.StatusOK, response)
}

// GetSamplesFiltered godoc
// @Summary     Retrieve all samplingfeatureIDs filtered by a variety of fields
// @Description Get all samplingfeatureIDs matching the current filters
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
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit           query    int    false "limit"
// @Param       offset          query    int    false "offset"
// @Param       setting         query    string false "tectonic setting - see /queries/sites/settings"
// @Param       location1       query    string false "location level 1 - see /queries/locations/l1"
// @Param       location2       query    string false "location level 2 - see /queries/locations/l2"
// @Param       location3       query    string false "location level 3 - see /queries/locations/l3"
// @Param       latitude        query    string false "latitude"
// @Param       longitude       query    string false "longitude"
// @Param       rocktype        query    string false "rock type - see /queries/samples/rocktypes"
// @Param       rockclass       query    string false "taxonomic classifier name - see /queries/samples/rockclasses"
// @Param       mineral         query    string false "mineral - see /queries/samples/minerals"
// @Param       material        query    string false "material - see /queries/samples/materials"
// @Param       inclusiontype   query    string false "inclusion type - see /queries/samples/inclusiontypes"
// @Param       sampletech      query    string false "sampling technique - see /queries/samples/samplingtechniques"
// @Param       element         query    string false "chemical element - see /queries/samples/elements"
// @Param       elementtype     query    string false "element type - see /queries/samples/elementtypes"
// @Param       value           query    string false "measured value"
// @Param       title           query    string false "title of publication"
// @Param       publicationyear query    string false "publication year"
// @Param       doi             query    string false "DOI"
// @Param       firstname       query    string false "Author first name"
// @Param       lastname        query    string false "Author last name"
// @Param       agemin          query    string false "Specimen age min"
// @Param       agemax          query    string false "Specimen age max"
// @Param       geoage          query    string false "Specimen geological age - see /queries/samples/geoages"
// @Param       geoageprefix    query    string false "Specimen geological age prefix - see /queries/samples/geoageprefixes"
// @Param       lab             query    string false "Laboratory name - see /queries/samples/organizationnames"
// @Param       polygon         query    string false "Coordinate-Polygon formatted as 2-dimensional json array: [[LONG,LAT],[2.4,6.3]]"
// @Success     200             {array}  model.SampleByFilterResponse
// @Failure     401             {object} string
// @Failure     404             {object} string
// @Failure     422             {object} string
// @Failure     500             {object} string
// @Router      /queries/samples [get]
func (h *Handler) GetSamplesFiltered(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	specimen := []model.SampleByFilters{}

	query, err := buildSampleFilterQuery(c)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)

	err = h.db.Query(query.GetQueryString(), &specimen, query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetSamplesFiltered: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	response := model.SampleByFilterResponse{NumItems: len(specimen), Data: specimen}
	return c.JSON(http.StatusOK, response)
}

// GetSamplesFilteredClustered godoc
// @Summary     Retrieve all samplingfeatureIDs filtered by a variety of fields and clustered
// @Description Get all samplingfeatureIDs matching the current filters clustered
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
// @Tags        geodata
// @Accept      json
// @Produce     json
// @Param       limit           query    int    false "limit"
// @Param       offset          query    int    false "offset"
// @Param       setting         query    string false "tectonic setting - see /queries/sites/settings"
// @Param       location1       query    string false "location level 1 - see /queries/locations/l1"
// @Param       location2       query    string false "location level 2 - see /queries/locations/l2"
// @Param       location3       query    string false "location level 3 - see /queries/locations/l3"
// @Param       latitude        query    string false "latitude"
// @Param       longitude       query    string false "longitude"
// @Param       rocktype        query    string false "rock type - see /queries/samples/rocktypes"
// @Param       rockclass       query    string false "taxonomic classifier name - see /queries/samples/rockclasses"
// @Param       mineral         query    string false "mineral - see /queries/samples/minerals"
// @Param       material        query    string false "material - see /queries/samples/materials"
// @Param       inclusiontype   query    string false "inclusion type - see /queries/samples/inclusiontypes"
// @Param       sampletech      query    string false "sampling technique - see /queries/samples/samplingtechniques"
// @Param       element         query    string false "chemical element - see /queries/samples/elements"
// @Param       elementtype     query    string false "element type - see /queries/samples/elementtypes"
// @Param       value           query    string false "measured value"
// @Param       title           query    string false "title of publication"
// @Param       publicationyear query    string false "publication year"
// @Param       doi             query    string false "DOI"
// @Param       firstname       query    string false "Author first name"
// @Param       lastname        query    string false "Author last name"
// @Param       agemin          query    string false "Specimen age min"
// @Param       agemax          query    string false "Specimen age max"
// @Param       geoage          query    string false "Specimen geological age - see /queries/samples/geoages"
// @Param       geoageprefix    query    string false "Specimen geological age prefix - see /queries/samples/geoageprefixes"
// @Param       lab             query    string false "Laboratory name - see /queries/samples/organizationnames"
// @Param       polygon         query    string false "Coordinate-Polygon formatted as 2-dimensional json array: [[LONG,LAT],[2.4,6.3]]"
// @Param       numClusters     query    int    false "Number of clusters for k-means clustering. Can be more or less depending on maxDistance"
// @Param       maxDistance     query    int    false "Max distance of points in cluster"
// @Success     200             {array}  []model.ClusteredSample{}
// @Failure     401             {object} string
// @Failure     404             {object} string
// @Failure     422             {object} string
// @Failure     500             {object} string
// @Router      /geodata/samples [get]
func (h *Handler) GetSamplesFilteredClustered(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	clusterData := []model.ClusteredSample{}

	query, err := buildSampleFilterQuery(c)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	numClusters := c.QueryParam(QP_NUM_CLUSTERS)
	if numClusters == "" {
		numClusters = DEFAULT_NUM_CLUSTERS
	}
	maxDistance := c.QueryParam(QP_MAX_DISTANCE)
	if maxDistance == "" {
		maxDistance = DEFAULT_MAX_DISTANCE
	}

	// wrap query in clustering postGIS-sql with parameters
	params := map[string]interface{}{
		"numClusters": numClusters,
		"maxDistance": maxDistance,
	}
	query.WrapInSQLParametrized(sql.GetSamplesClusteredWrapperPrefix, sql.GetSamplesClusteredWrapperPostfix, params)

	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)

	err = h.db.Query(query.GetQueryString(), &clusterData, query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetSamplesFilteredClustered: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	response, err := parseClusterToGeoJSON(clusterData)
	if err != nil {
		logger.Errorf("Can not parse cluster data: %v", err)
		return c.String(http.StatusInternalServerError, "Can not parse cluster data")
	}
	return c.JSON(http.StatusOK, response)
}

// GetSpecimenTypes godoc
// @Summary     Retrieve specimen types
// @Description get specimen types
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.Specimen
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/specimentypes [get]
func (h *Handler) GetSpecimenTypes(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	specimentypes := []model.Specimen{}
	query := sql.NewQuery(sql.GetSpecimenTypesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &specimentypes)
	if err != nil {
		logger.Errorf("Can not GetSpecimenTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve specimentype data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(specimentypes), specimentypes}
	return c.JSON(http.StatusOK, response)
}

// GetRockClasses godoc
// @Summary     Retrieve rock classes
// @Description get rock classes
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.TaxonomicClassifier
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/rockclasses [get]
func (h *Handler) GetRockClasses(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	rockclasses := []model.TaxonomicClassifier{}
	query := sql.NewQuery(sql.RockClassQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &rockclasses)
	if err != nil {
		logger.Errorf("Can not GetRockClasses: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve rock class data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(rockclasses), rockclasses}
	return c.JSON(http.StatusOK, response)
}

// GetRockTypes godoc
// @Summary     Retrieve rock types
// @Description get rock types
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.TaxonomicClassifier
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/rocktypes [get]
func (h *Handler) GetRockTypes(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	rocktypes := []model.TaxonomicClassifier{}
	query := sql.NewQuery(sql.RockTypeQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &rocktypes)
	if err != nil {
		logger.Errorf("Can not GetRockTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve rock type data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(rocktypes), rocktypes}
	return c.JSON(http.StatusOK, response)
}

// GetMinerals godoc
// @Summary     Retrieve minerals
// @Description get minerals
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.TaxonomicClassifier
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/minerals [get]
func (h *Handler) GetMinerals(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	minerals := []model.TaxonomicClassifier{}
	query := sql.NewQuery(sql.MineralQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &minerals)
	if err != nil {
		logger.Errorf("Can not GetMinerals: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve mineral data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(minerals), minerals}
	return c.JSON(http.StatusOK, response)
}

// GetMaterials godoc
// @Summary     Retrieve materials
// @Description get materials
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.Material
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/materials [get]
func (h *Handler) GetMaterials(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	materials := []model.Material{}
	query := sql.NewQuery(sql.MaterialsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &materials)
	if err != nil {
		logger.Errorf("Can not GetMaterials: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve material data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(materials), materials}
	return c.JSON(http.StatusOK, response)
}

// GetInclusionTypes godoc
// @Summary     Retrieve inclusion types
// @Description get inclusion types
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.InclusionType
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/inclusiontypes [get]
func (h *Handler) GetInclusionTypes(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	inclusionTypes := []model.InclusionType{}
	query := sql.NewQuery(sql.InclusionTypesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &inclusionTypes)
	if err != nil {
		logger.Errorf("Can not GetInclusionTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve inclusion type data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(inclusionTypes), inclusionTypes}
	return c.JSON(http.StatusOK, response)
}

// GetSamplingTechniques godoc
// @Summary     Retrieve sampling techniques
// @Description get sampling techniques
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.SamplingTechnique
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/samplingtechniques [get]
func (h *Handler) GetSamplingTechniques(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	samplingtechniques := []model.SamplingTechnique{}
	query := sql.NewQuery(sql.SamplingTechniquesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &samplingtechniques)
	if err != nil {
		logger.Errorf("Can not GetSamplingTechniques: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sampling technique data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(samplingtechniques), samplingtechniques}
	return c.JSON(http.StatusOK, response)
}

// GetRandomSamples godoc
// @Summary     Retrieve a random set of specimen
// @Description get random specimen
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit query    int false "limit"
// @Success     200   {array}  model.Specimen
// @Failure     401   {object} string
// @Failure     404   {object} string
// @Failure     422   {object} string
// @Failure     500   {object} string
// @Router      /queries/samples/random [get]
func (h *Handler) GetRandomSamples(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	randomSpecimen := []model.Specimen{}
	query := sql.NewQuery(sql.GetRandomSpecimensQuery)
	limit := c.QueryParam(QP_LIMIT)
	err := h.db.Query(query.GetQueryString(), &randomSpecimen, limit)
	if err != nil {
		logger.Errorf("Can not GetRandomSamples: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve random data sample")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(randomSpecimen), randomSpecimen}
	return c.JSON(http.StatusOK, response)
}

// GetGeoAges godoc
// @Summary     Retrieve geological ages
// @Description get geological ages
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.GeoAge
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/geoages [get]
func (h *Handler) GetGeoAges(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	geoAges := []model.GeoAge{}
	query := sql.NewQuery(sql.GetGeoAgesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &geoAges)
	if err != nil {
		logger.Errorf("Can not GetGeoAges: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geological age data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(geoAges), geoAges}
	return c.JSON(http.StatusOK, response)
}

// GetGeoAgePrefixes godoc
// @Summary     Retrieve geological age prefixes
// @Description get geological age prefixes
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.GeoAgePrefix
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/geoageprefixes [get]
func (h *Handler) GetGeoAgePrefixes(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	geoAgePrefixes := []model.GeoAgePrefix{}
	query := sql.NewQuery(sql.GetGeoAgePrefixesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &geoAgePrefixes)
	if err != nil {
		logger.Errorf("Can not GetGeoAgePrefixes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geological age prefix data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(geoAgePrefixes), geoAgePrefixes}
	return c.JSON(http.StatusOK, response)
}

// GetOrganizationNames godoc
// @Summary     Retrieve organization names
// @Description get organization names
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.Organization
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/organizationnames [get]
func (h *Handler) GetOrganizationNames(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	organizations := []model.Organization{}
	query := sql.NewQuery(sql.GetOrganizationNamesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &organizations)
	if err != nil {
		logger.Errorf("Can not GetOrganizationNames: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve organization name data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(organizations), organizations}
	return c.JSON(http.StatusOK, response)
}

// buildSampleFilterQuery constructs a query using filter params from the request
func buildSampleFilterQuery(c echo.Context) (*sql.Query, error) {
	query := sql.NewQuery(sql.GetSamplingfeatureIdsByFilterBaseQuery)

	// add optional search filters
	junctor := sql.OpWhere // junctor to connect a new filter clause to the query: can be "WHERE" or "AND/OR"
	// location filters
	setting, opSetting, err := parseParam(c.QueryParam(QP_SETTING))
	if err != nil {
		return nil, err
	}
	location1, opLoc1, err := parseParam(c.QueryParam(QP_LOC1))
	if err != nil {
		return nil, err
	}
	location2, opLoc2, err := parseParam(c.QueryParam(QP_LOC2))
	if err != nil {
		return nil, err
	}
	location3, opLoc3, err := parseParam(c.QueryParam(QP_LOC3))
	if err != nil {
		return nil, err
	}
	lat, opLat, err := parseParam(c.QueryParam(QP_LAT))
	if err != nil {
		return nil, err
	}
	long, opLong, err := parseParam(c.QueryParam(QP_LONG))
	if err != nil {
		return nil, err
	}
	if setting != "" || location1 != "" || location2 != "" || location3 != "" || lat != "" || long != "" {
		// add query module Location
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterLocationsStart)
		// add location filters
		if setting != "" {
			query.AddFilter("s.setting", setting, opSetting, junctor)
			junctor = sql.OpAnd // after first filter is added with "WHERE", change to "AND" for following filters
		}
		if location1 != "" {
			query.AddFilter("toplevelloc.locationname", location1, opLoc1, junctor)
			junctor = sql.OpAnd
		}
		if location2 != "" {
			query.AddFilter("secondlevelloc.locationname", location2, opLoc2, junctor)
			junctor = sql.OpAnd
		}
		if location3 != "" {
			query.AddFilter("thirdlevelloc.locationname", location3, opLoc3, junctor)
			junctor = sql.OpAnd
		}
		if lat != "" {
			query.AddFilter("s.latitude", lat, opLat, junctor)
			junctor = sql.OpAnd
		}
		if long != "" {
			query.AddFilter("s.longitude", long, opLong, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterLocationsEnd)
	}

	// taxonomic classifiers
	junctor = sql.OpWhere // reset junctor for new subquery
	rockType, opRType, err := parseParam(c.QueryParam(QP_ROCKTYPE))
	if err != nil {
		return nil, err
	}
	rockClass, opRClass, err := parseParam(c.QueryParam(QP_ROCKCLASS))
	if err != nil {
		return nil, err
	}
	mineral, opMin, err := parseParam(c.QueryParam(QP_MINERAL))
	if err != nil {
		return nil, err
	}
	if rockType != "" || rockClass != "" || mineral != "" {
		// add query module taxonomic classifiers
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersStart)
		// add taxonomic classifiers filters
		if rockType != "" {
			query.AddFilter("tax_type.taxonomicclassifiername", rockType, opRType, junctor)
			junctor = sql.OpAnd
		}
		if rockClass != "" {
			query.AddFilter("tax_class.taxonomicclassifiername", rockClass, opRClass, junctor)
			junctor = sql.OpAnd
		}
		if mineral != "" {
			query.AddFilter("tax_min.taxonomicclassifiercommonname", mineral, opMin, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersEnd)
	}

	// annotations
	junctor = sql.OpWhere // reset junctor for new subquery
	material, opMat, err := parseParam(c.QueryParam(QP_MATERIAL))
	if err != nil {
		return nil, err
	}
	incType, opIncType, err := parseParam(c.QueryParam(QP_INCTYPE))
	if err != nil {
		return nil, err
	}
	sampTech, opSampTech, err := parseParam(c.QueryParam(QP_SAMPTECH))
	if err != nil {
		return nil, err
	}
	if material != "" || incType != "" || sampTech != "" {
		// add query module annotations
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterAnnotationsStart)
		// add annotaion filters
		if material != "" {
			query.AddFilter("ann_mat.annotationtext", material, opMat, junctor)
			junctor = sql.OpAnd
		}
		if incType != "" {
			query.AddFilter("ann_inc_type.annotationtext", incType, opIncType, junctor)
			junctor = sql.OpAnd
		}
		if sampTech != "" {
			query.AddFilter("ann_samp_tech.annotationtext", sampTech, opSampTech, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterAnnotationsEnd)
	}

	// results
	junctor = sql.OpWhere // reset junctor for new subquery
	elem, opElem, err := parseParam(c.QueryParam(QP_ELEM))
	if err != nil {
		return nil, err
	}
	elemType, opElemType, err := parseParam(c.QueryParam(QP_ELEMTYPE))
	if err != nil {
		return nil, err
	}
	value, opValue, err := parseParam(c.QueryParam(QP_VALUE))
	if err != nil {
		return nil, err
	}
	if elem != "" || elemType != "" || value != "" {
		// add query module results
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterResultsStart)
		if elem != "" {
			query.AddFilter("mv.variablecode", elem, opElem, junctor)
			junctor = sql.OpAnd
		}
		if elemType != "" {
			query.AddFilter("mv.variabletypecode", elemType, opElemType, junctor)
			junctor = sql.OpAnd
		}
		if value != "" {
			query.AddFilter("mv.datavalue", value, opValue, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterResultsEnd)
	}

	// citation
	junctor = sql.OpWhere // reset junctor for new subquery
	title, opTitle, err := parseParam(c.QueryParam(QP_TITLE))
	if err != nil {
		return nil, err
	}
	pubYear, opPubYear, err := parseParam(c.QueryParam(QP_PUBYEAR))
	if err != nil {
		return nil, err
	}
	doi, opDOI, err := parseParam(c.QueryParam(QP_DOI))
	if err != nil {
		return nil, err
	}
	authorFirst, opAuthorFirst, err := parseParam(c.QueryParam(QP_AUTHOR_FIRST))
	if err != nil {
		return nil, err
	}
	authorLast, opAuthorLast, err := parseParam(c.QueryParam(QP_AUTHOR_LAST))
	if err != nil {
		return nil, err
	}
	if title != "" || pubYear != "" || doi != "" || authorFirst != "" || authorLast != "" {
		// add query module citations
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterCitationsStart)
		if title != "" {
			query.AddFilter("c.title", title, opTitle, junctor)
			junctor = sql.OpAnd
		}
		if pubYear != "" {
			query.AddFilter("c.publicationyear", pubYear, opPubYear, junctor)
			junctor = sql.OpAnd
		}
		if doi != "" {
			query.AddFilter("cid.citationexternalidentifier", doi, opDOI, junctor)
			query.AddFilter("e.externalidentifiersystemname", "DOI", sql.OpEq, sql.OpAnd)
			junctor = sql.OpAnd
		}
		if authorFirst != "" {
			query.AddFilter("p.personfirstname", authorFirst, opAuthorFirst, junctor)
			junctor = sql.OpAnd
		}
		if authorLast != "" {
			query.AddFilter("p.personlastname", authorLast, opAuthorLast, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterCitationsEnd)
	}

	// Ages
	junctor = sql.OpWhere // reset junctor for new subquery
	ageMin, opAgeMin, err := parseParam(c.QueryParam(QP_AGE_MIN))
	if err != nil {
		return nil, err
	}
	ageMax, opAgeMax, err := parseParam(c.QueryParam(QP_AGE_MAX))
	if err != nil {
		return nil, err
	}
	geoAge, opGeoAge, err := parseParam(c.QueryParam(QP_GEO_AGE))
	if err != nil {
		return nil, err
	}
	geoPrefix, opGeoPrefix, err := parseParam(c.QueryParam(QP_GEO_AGE_PREFIX))
	if err != nil {
		return nil, err
	}
	if ageMin != "" || ageMax != "" || geoAge != "" || geoPrefix != "" {
		// add query module age
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterAgesStart)
		if ageMin != "" {
			query.AddFilter("sa.specimenagemin", ageMin, opAgeMin, junctor)
			junctor = sql.OpAnd
		}
		if ageMax != "" {
			query.AddFilter("sa.specimenagemax", ageMax, opAgeMax, junctor)
			junctor = sql.OpAnd
		}
		if geoAge != "" {
			query.AddFilter("sa.specimengeolage", geoAge, opGeoAge, junctor)
			junctor = sql.OpAnd
		}
		if geoPrefix != "" {
			query.AddFilter("sa.specimengeolageprefix", geoPrefix, opGeoPrefix, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterAgesEnd)
	}

	// Organizations
	junctor = sql.OpWhere // reset junctor for new subquery
	labName, opLabName, err := parseParam(c.QueryParam(QP_LAB))
	if err != nil {
		return nil, err
	}
	if labName != "" {
		// add query module organizations
		query.AddSQLBlock(sql.GestSamplingfeatureIdsByFilterOrganizationsStart)
		if labName != "" {
			query.AddFilter("o.organizationname", labName, opLabName, junctor)
		}
		query.AddSQLBlock(sql.GestSamplingfeatureIdsByFilterOrganizationsEnd)
	}

	// Geometries
	junctor = sql.OpWhere // reset junctor for new subquery
	polygonString, _, err := parseParam(c.QueryParam(QP_POLY))
	if err != nil {
		return nil, err
	}
	if polygonString != "" {
		// add query module geometry
		query.AddSQLBlock(sql.GestSamplingfeatureIdsByFilterGeometryStart)
		if polygonString != "" {
			// format polygon string for postGIS/SQL syntax
			polygonFormatted, err := formatPolygonArray(polygonString)
			if err != nil {
				return nil, err
			}
			query.AddFilter("sg.geometry", polygonFormatted, sql.OpInPolygon, junctor)
		}
		query.AddSQLBlock(sql.GestSamplingfeatureIdsByFilterGeometryEnd)
	}
	return query, nil
}

// formatPolygonArray formats a given input polygon for usage in postGIS/SQL syntax
// Input is json formatted: [[val1,val2],[val1,val2],...]
// Output is postGIS geometry syntax: (val1 val2, val1 val2, ...)
func formatPolygonArray(polygonString string) (string, error) {
	polygon := [][]float64{}
	err := json.Unmarshal([]byte(polygonString), &polygon)
	if err != nil {
		return "", err
	}
	output := "("
	for i, point := range polygon {
		if i > 0 {
			// add separator before adding next point
			output += ","
		}
		for _, coordinate := range point {
			output += fmt.Sprintf(" %f", coordinate)
		}
	}
	output += ")"
	return output, nil
}

// parseClusterToGeoJSON takes an array of model.ClusteredSample and parses it into GeoJSON
func parseClusterToGeoJSON(clusterData []model.ClusteredSample) ([]model.GeoJSONCluster, error) {
	geoJSON := make([]model.GeoJSONCluster, 0, len(clusterData))
	for _, cluster := range clusterData {
		centroid := model.GeoJSONFeature{
			Type:     model.GEOJSONTYPE_FEATURE,
			Geometry: cluster.Centroid,
		}
		convexHull := model.GeoJSONFeature{
			Type:     model.GEOJSONTYPE_FEATURE,
			Geometry: cluster.ConvexHull,
		}
		points := model.GeoJSONFeature{
			Type:     model.GEOJSONTYPE_FEATURE,
			Geometry: cluster.Points,
		}
		geoJSONCluster := model.GeoJSONCluster{
			ClusterID:  cluster.ClusterID,
			Centroid:   centroid,
			ConvexHull: convexHull,
			Points:     points,
			Samples:    cluster.Samples,
		}
		geoJSON = append(geoJSON, geoJSONCluster)
	}
	return geoJSON, nil
}
