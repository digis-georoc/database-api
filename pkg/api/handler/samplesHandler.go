package handler

import (
	"fmt"
	"net/http"
	"strings"

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

	QP_ROCKTYPE  = "rocktype"
	QP_ROCKCLASS = "rockclass"
	QP_MINERAL   = "mineral"

	QP_MATERIAL = "material"
	QP_INCTYPE  = "inclusiontype"
	QP_SAMPTECH = "sampletech"

	QP_ELEM     = "element"
	QP_ELEMTYPE = "elementtype"
	QP_VALUE    = "value"
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
// @Description where FIELD is one of the accepted query params; OPERATOR is one of "lt", "gt", "eq", "in" and VALUE is an unquoted string, integer or decimal
// @Description Multiple VALUEs for an "in"-filter must be comma-separated and will be interpreted as a discunctive filter.
// @Description The OPERATORs "lt" and "gt" are only applicable to numerical values.
// @Description If no OPERATOR is specified, "eq" is assumed as the default OPERATOR
// @Description The filters are evaluated conjunctively.
// @Description Note that applying more filters can slow down the query as more tables have to be considered in the evaluation.
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit         query    int    false "limit"
// @Param       offset        query    int    false "offset"
// @Param       setting       query    string false "tectonic setting - see /queries/sites/settings"
// @Param       location1     query    string false "location level 1 - see /queries/locations/l1"
// @Param       location2     query    string false "location level 2 - see /queries/locations/l2"
// @Param       location3     query    string false "location level 3 - see /queries/locations/l3"
// @Param       rocktype      query    string false "rock type - see /queries/samples/rocktypes"
// @Param       rockclass     query    string false "taxonomic classifier name - see /queries/samples/rockclasses"
// @Param       mineral       query    string false "mineral - see /queries/samples/minerals"
// @Param       material      query    string false "material - see /queries/samples/materials"
// @Param       inclusiontype query    string false "inclusion type - see /queries/samples/inclusiontypes"
// @Param       sampletech    query    string false "sampling technique - see /queries/samples/samplingtechniques"
// @Param       element       query    string false "chemical element - see /queries/samples/elements"
// @Param       elementtype   query    string false "element type - see /queries/samples/elementtypes"
// @Param       value         query    number false "measured value"
// @Success     200           {array}  model.SampleByFiltersResponse
// @Failure     401           {object} string
// @Failure     404           {object} string
// @Failure     422           {object} string
// @Failure     500           {object} string
// @Router      /queries/samples [get]
func (h *Handler) GetSamplesFiltered(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	specimen := []model.SampleByFiltersResponse{}
	query := sql.NewQuery(sql.GetSamplingfeatureIdsByFilterBaseQuery)

	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)

	// add optional search filters
	junctor := sql.OpWhere // junctor to connect a new filter clause to the query: can be "WHERE" or "AND/OR"
	// location filters
	setting, opSetting, err := parseParam(c.QueryParam(QP_SETTING))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	location1, opLoc1, err := parseParam(c.QueryParam(QP_LOC1))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	location2, opLoc2, err := parseParam(c.QueryParam(QP_LOC2))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	location3, opLoc3, err := parseParam(c.QueryParam(QP_LOC3))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	if setting != "" || location1 != "" || location2 != "" || location3 != "" {
		// add query module Location
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterLocationsStart)
		// add location filters
		if setting != "" {
			query.AddFilter("s.setting", setting, opSetting, junctor)
			junctor = sql.OpAnd
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
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterLocationsEnd)
	}

	// taxonomic classifiers
	junctor = sql.OpWhere
	rockType, opRType, err := parseParam(c.QueryParam(QP_ROCKTYPE))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	rockClass, opRClass, err := parseParam(c.QueryParam(QP_ROCKCLASS))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	mineral, opMin, err := parseParam(c.QueryParam(QP_MINERAL))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
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
	junctor = sql.OpWhere
	material, opMat, err := parseParam(c.QueryParam(QP_MATERIAL))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	incType, opIncType, err := parseParam(c.QueryParam(QP_INCTYPE))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	sampTech, opSampTech, err := parseParam(c.QueryParam(QP_SAMPTECH))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
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
	junctor = sql.OpWhere
	elem, opElem, err := parseParam(c.QueryParam(QP_ELEM))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	elemType, opElemType, err := parseParam(c.QueryParam(QP_ELEMTYPE))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	value, opValue, err := parseParam(c.QueryParam(QP_VALUE))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
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

	err = h.db.Query(query.GetQueryString(), &specimen, query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetSamplesFiltered: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(specimen), specimen}
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

// parseParam parses a given query parameter and validates the contents
func parseParam(queryParam string) (string, string, error) {
	if queryParam == "" {
		return "", "", nil
	}
	operator, value, found := strings.Cut(queryParam, ":")
	if !found {
		// if no operator is specified, "eq" is assumed as default
		return queryParam, sql.OpEq, nil
	}
	// validate operator
	operator, opIsValid := sql.OperatorMap[operator]
	if !opIsValid {
		return "", "", fmt.Errorf("Invalid operator")
	}
	return value, operator, nil
}
