package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/geometry"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_SETTING = "setting"
	QP_LOC1    = "location1"
	QP_LOC2    = "location2"
	QP_LOC3    = "location3"
	QP_LAT     = "latitude"
	QP_LONG    = "longitude"

	QP_ROCKTYPE     = "rocktype"
	QP_ROCKCLASS    = "rockclass"
	QP_MINERAL      = "mineral"
	QP_INCLUSIONMAT = "inclusionmaterial"
	QP_HOSTMAT      = "hostmaterial"

	QP_ROCKCLASS_QUERY = "q"

	QP_MATERIAL = "material"
	QP_INCTYPE  = "inclusiontype"
	QP_SAMPTECH = "sampletech"

	QP_CHEMISTRY = "chemistry"

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

	QP_BBOX = "bbox"

	QP_NUM_CLUSTERS      = "numClusters"
	QP_MAX_DISTANCE      = "maxDistance"
	DEFAULT_NUM_CLUSTERS = "1"  // 1 produces any number of clusters that satisfy the max_distance but prevents error where fewer samples that NUM_CLUSTERS exist
	DEFAULT_MAX_DISTANCE = "50" // not in use currently

	CLUSTERING_THRESHOLD = 15 // clusters with less points are returned as individual points instead of a cluster

	KEY_BBOX                    = "key_bbox"
	KEY_TRANSLATION_FACTOR      = "key_translation_factor"
	KEY_BOUNDARY                = "key_boundary"
	KEY_POLYGON                 = "key_polygon"
	KEY_TRANSLATION_FACTOR_POLY = "key_translation_factor_poly"
	KEY_BOUNDARY_POLY           = "key_boundary_poly"

	KEY_ROCKTYPE  = "key_rocktype"
	KEY_ROCKCLASS = "key_rockclass"
)

// GetSampleByID godoc
// @Summary     Retrieve sample by samplingfeatureid
// @Description get sample by samplingfeatureid
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       samplingfeatureID path     string true "Sample ID"
// @Success     200               {object} model.SampleResponse
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /queries/samples/{samplingfeatureID} [get]
func (h *Handler) GetSampleByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	query := sql.NewQuery(sql.GetSampleByIDQuery)
	samples, err := repository.Query[model.Sample](c.Request().Context(), h.db, query.GetQueryString(), c.Param(QP_SAMPFEATUREID))
	if err != nil {
		logger.Errorf("Can not GetSampleByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	num := len(samples)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}
	response := model.SampleResponse{
		NumItems: num,
		Data:     samples,
	}
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
// @Param       limit             query    int    false "limit"
// @Param       offset            query    int    false "offset"
// @Param       setting           query    string false "tectonic setting - see /queries/sites/settings (supports Filter DSL)"
// @Param       location1         query    string false "location level 1 - see /queries/locations/l1 (supports Filter DSL)"
// @Param       location2         query    string false "location level 2 - see /queries/locations/l2 (supports Filter DSL)"
// @Param       location3         query    string false "location level 3 - see /queries/locations/l3 (supports Filter DSL)"
// @Param       latitude          query    string false "latitude (supports Filter DSL)"
// @Param       longitude         query    string false "longitude (supports Filter DSL)"
// @Param       rocktype          query    string false "rock type - see /queries/samples/rocktypes (supports Filter DSL)"
// @Param       rockclass         query    string false "taxonomic classifier name - see /queries/samples/rockclasses (supports Filter DSL)"
// @Param       mineral           query    string false "mineral - see /queries/samples/minerals (supports Filter DSL)"
// @Param       material          query    string false "material - see /queries/samples/materials (supports Filter DSL)"
// @Param       inclusiontype     query    string false "inclusion type - see /queries/samples/inclusiontypes (supports Filter DSL)"
// @Param       hostmaterial      query    string false "host material - see /queries/samples/hostmaterials (supports Filter DSL)"
// @Param       inclusionmaterial query    string false "inclusion material - see /queries/samples/inclusionmaterials (supports Filter DSL)"
// @Param       sampletech        query    string false "sampling technique - see /queries/samples/samplingtechniques (supports Filter DSL)"
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
// @Success     200               {object} model.SampleByFilterResponse
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     422               {object} string
// @Failure     500               {object} string
// @Router      /queries/samples [get]
func (h *Handler) GetSamplesFiltered(c echo.Context) error {
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
	kwargs := map[string]interface{}{
		KEY_ROCKCLASS: true,
		KEY_ROCKTYPE:  true,
	}
	query, err := buildSampleFilterQuery(c, coordData, kwargs)
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

	result, err := repository.Query[model.SampleByFilters](c.Request().Context(), h.db, query.GetQueryString(), query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetSamplesFiltered: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	// copy into model without totalCount on each sample
	responseData := []model.SampleByFiltersData{}
	totalCount := 0
	for _, sample := range result {
		totalCount = sample.TotalCount
		data := model.SampleByFiltersData{
			SampleID:   sample.SampleID,
			SampleName: sample.SampleName,
			Latitude:   sample.Latitude,
			Longitude:  sample.Longitude,
			RockType:   sample.RockType,
			RockClass:  sample.RockClass,
		}
		responseData = append(responseData, data)
	}
	response := model.SampleByFilterResponse{NumItems: len(responseData), TotalCount: totalCount, Data: responseData}
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
// @Param       limit             query    int    false "limit"
// @Param       offset            query    int    false "offset"
// @Param       setting           query    string false "tectonic setting - see /queries/sites/settings (supports Filter DSL)"
// @Param       location1         query    string false "location level 1 - see /queries/locations/l1 (supports Filter DSL)"
// @Param       location2         query    string false "location level 2 - see /queries/locations/l2 (supports Filter DSL)"
// @Param       location3         query    string false "location level 3 - see /queries/locations/l3 (supports Filter DSL)"
// @Param       latitude          query    string false "latitude (supports Filter DSL)"
// @Param       longitude         query    string false "longitude (supports Filter DSL)"
// @Param       rocktype          query    string false "rock type - see /queries/samples/rocktypes (supports Filter DSL)"
// @Param       rockclass         query    string false "taxonomic classifier name - see /queries/samples/rockclasses (supports Filter DSL)"
// @Param       mineral           query    string false "mineral - see /queries/samples/minerals (supports Filter DSL)"
// @Param       material          query    string false "material - see /queries/samples/materials (supports Filter DSL)"
// @Param       inclusiontype     query    string false "inclusion type - see /queries/samples/inclusiontypes (supports Filter DSL)"
// @Param       hostmaterial      query    string false "host material - see /queries/samples/hostmaterials (supports Filter DSL)"
// @Param       inclusionmaterial query    string false "inclusion material - see /queries/samples/inclusionmaterials (supports Filter DSL)"
// @Param       sampletech        query    string false "sampling technique - see /queries/samples/samplingtechniques (supports Filter DSL)"
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
// @Param       bbox              query    string true  "BoundingBox formatted as 2-dimensional json array: [[SW_Long,SW_Lat],[SE_Long,SE_Lat],[NE_Long,NE_Lat],[NW_Long,NW_Lat]]"
// @Param       numClusters       query    int    false "Number of clusters for k-means clustering. Default is 7. Can be more depending on maxDistance"
// @Param       maxDistance       query    int    false "Max size of cluster. Recommended values per zoom-level: Z0: 50, Z1: 50, Z2: 25, Z4: 12 -> Zi = 50/i"
// @Success     200               {object} model.ClusterResponse
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     422               {object} string
// @Failure     500               {object} string
// @Router      /geodata/samplesclustered [get]
func (h *Handler) GetSamplesFilteredClustered(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	// response object
	response := model.ClusterResponse{}

	coordData := map[string]interface{}{}
	// get the bbox
	bboxString, _, err := parseParam(c.QueryParam(QP_BBOX))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not parse bbox")
	}
	if bboxString == "" {
		return c.String(http.StatusInternalServerError, "No bbox provided")
	}
	bbox, err := geometry.ParsePointArray(bboxString)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not parse bbox")
	}
	// calc clustering param relative to original (visible) bbox size (max 1 world truncated)
	bbox = geometry.TruncateBBox(bbox)
	width := bbox[1][0] - bbox[0][0]
	kmeansMaxDistance := width / 12
	// scale bbox
	if !geometry.IsZoom0(bbox) {
		// add frame around bbox to avoid reloading on small panning
		bbox = geometry.ScaleBBox(bbox)
	}
	// truncate bbox after scaling so it contains at most one whole world
	bbox = geometry.TruncateBBox(bbox)
	// add first point again to make closed polygon shape
	bbox = append(bbox, bbox[0])
	boundary, translationFactor, err := geometry.CalcTranslation(bbox)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not calculate bbox translation")
	}
	coordData[KEY_BBOX] = bbox
	coordData[KEY_TRANSLATION_FACTOR] = translationFactor
	coordData[KEY_BOUNDARY] = boundary

	// get polygon filter
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

	// build query string
	kwargs := map[string]interface{}{
		KEY_ROCKCLASS: true,
		KEY_ROCKTYPE:  true,
	}
	query, err := buildSampleFilterQuery(c, coordData, kwargs)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	numClusters := c.QueryParam(QP_NUM_CLUSTERS)
	if numClusters == "" {
		numClusters = DEFAULT_NUM_CLUSTERS
	}
	maxDistance := c.QueryParam(QP_MAX_DISTANCE)
	if maxDistance == "" {
		// set relative to bbox size
		maxDistance = fmt.Sprintf("%f", kmeansMaxDistance)
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

	clusterData, err := repository.Query[model.ClusteredSample](c.Request().Context(), h.db, query.GetQueryString(), query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetSamplesFilteredClustered: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}

	// wrap bbox in []interface{} for geoJSON polygon
	bboxIWrap := []interface{}{bbox}
	response.Bbox = model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type:        model.GEOJSON_GEOMETRY_POLYGON,
			Coordinates: bboxIWrap,
		},
	}
	geoJSONClusters, geoJSONPoints, err := parseClusterToGeoJSON(clusterData)
	if err != nil {
		logger.Errorf("Can not parse cluster data: %v", err)
		return c.String(http.StatusInternalServerError, "Can not parse cluster data")
	}
	response.Clusters = geoJSONClusters
	response.Points = geoJSONPoints
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
// @Success     200    {object} model.SpecimenTypeResponse
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

	query := sql.NewQuery(sql.GetSpecimenTypesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	specimentypes, err := repository.Query[model.SpecimenType](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetSpecimenTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve specimentype data")
	}
	response := model.SpecimenTypeResponse{
		NumItems: len(specimentypes),
		Data:     specimentypes,
	}
	return c.JSON(http.StatusOK, response)
}

// GetRockClasses godoc
// @Summary     Retrieve rock classes
// @Description get rock classes
// @Description Filter DSL syntax:
// @Description FIELD=OPERATOR:VALUE
// @Description where FIELD is one of the accepted query params; OPERATOR is either "in" (IN) for rocktype; or "lk" (LIKE) for the search query q
// @Description and VALUE is an unquoted string
// @Description Multiple VALUEs for an "in"-filter must be comma-separated and will be interpreted as a discunctive filter.
// @Description The OPERATOR "lk" is only applicable to string values and supports wildcards `*`(0 or more chars) and `?`(one char).
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit    query    int    false "limit"
// @Param       offset   query    int    false "offset"
// @Param       rocktype query    string false "One or more Rocktypes to filter corresponding Rockclasses as a comma-separated list. Use "in" as the operator"
// @Param       q        query    string false "Search string for rockclass values. Use "lk:" as the operator"
// @Success     200      {object} model.TaxonomicClassifierResponse
// @Failure     401      {object} string
// @Failure     404      {object} string
// @Failure     422      {object} string
// @Failure     500      {object} string
// @Router      /queries/samples/rockclasses [get]
func (h *Handler) GetRockClasses(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	query := sql.NewQuery(sql.RockClassQueryStart)

	rocktypes, _, err := parseParam(c.QueryParam(QP_ROCKTYPE))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Invalid rocktype parameter. Expected format: `in:VAL1,VAL2`")
	}
	if rocktypes != "" {
		// add query filter for rock types
		query.AddFilter("t2.taxonomicclassifiername", rocktypes, sql.OpIn, sql.OpAnd)
	}
	query.AddSQLBlock(sql.RockClassQueryMid)
	rockClassQuery, _, err := parseParam(c.QueryParam(QP_ROCKCLASS_QUERY))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Invalid string search parameter q. Expected format: `lk:*F_o`")
	}
	if rockClassQuery != "" {
		// add query filter for search query
		query.AddFilter("t.taxonomicclassifiername", rockClassQuery, sql.OpLike, sql.OpAnd)
	}
	query.AddSQLBlock(sql.RockClassQueryEnd)

	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	rockclasses, err := repository.Query[model.TaxonomicClassifier](c.Request().Context(), h.db, query.GetQueryString(), query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetRockClasses: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve rock class data")
	}
	response := model.TaxonomicClassifierResponse{
		NumItems: len(rockclasses),
		Data:     rockclasses,
	}
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
// @Success     200    {object} model.TaxonomicClassifierResponse
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

	query := sql.NewQuery(sql.RockTypeQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	rocktypes, err := repository.Query[model.TaxonomicClassifier](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetRockTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve rock type data")
	}
	response := model.TaxonomicClassifierResponse{
		NumItems: len(rocktypes),
		Data:     rocktypes,
	}
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
// @Success     200    {object} model.TaxonomicClassifierResponse
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

	query := sql.NewQuery(sql.MineralQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	minerals, err := repository.Query[model.TaxonomicClassifier](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetMinerals: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve mineral data")
	}
	response := model.TaxonomicClassifierResponse{
		NumItems: len(minerals),
		Data:     minerals,
	}
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
// @Success     200    {object} model.MaterialResponse
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

	query := sql.NewQuery(sql.MaterialsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	materials, err := repository.Query[model.Material](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetMaterials: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve material data")
	}
	response := model.MaterialResponse{
		NumItems: len(materials),
		Data:     materials,
	}
	return c.JSON(http.StatusOK, response)
}

// GetHostMaterials godoc
// @Summary     Retrieve host materials
// @Description get host materials
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.TaxonomicClassifierResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/hostmaterials [get]
func (h *Handler) GetHostMaterials(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	query := sql.NewQuery(sql.HostMatQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	hostMaterials, err := repository.Query[model.TaxonomicClassifier](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetHostMaterials: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve host material data")
	}
	response := model.TaxonomicClassifierResponse{
		NumItems: len(hostMaterials),
		Data:     hostMaterials,
	}
	return c.JSON(http.StatusOK, response)
}

// GetInclusionMaterials godoc
// @Summary     Retrieve inclusion materials
// @Description get inclusion materials
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.TaxonomicClassifierResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples/inclusionmaterials [get]
func (h *Handler) GetInclusionMaterials(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	query := sql.NewQuery(sql.IncMatQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	incMaterials, err := repository.Query[model.TaxonomicClassifier](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetInclusionMaterials: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve inclusion material data")
	}
	response := model.TaxonomicClassifierResponse{
		NumItems: len(incMaterials),
		Data:     incMaterials,
	}
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
// @Success     200    {object} model.InclusionTypeResponse
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

	query := sql.NewQuery(sql.InclusionTypesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	inclusionTypes, err := repository.Query[model.InclusionType](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetInclusionTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve inclusion type data")
	}
	response := model.InclusionTypeResponse{
		NumItems: len(inclusionTypes),
		Data:     inclusionTypes,
	}
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
// @Success     200    {object} model.SamplingTechniqueResponse
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

	query := sql.NewQuery(sql.SamplingTechniquesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	samplingtechniques, err := repository.Query[model.SamplingTechnique](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetSamplingTechniques: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sampling technique data")
	}
	response := model.SamplingTechniqueResponse{
		NumItems: len(samplingtechniques),
		Data:     samplingtechniques,
	}
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
// @Success     200   {object} model.SpecimenResponse
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

	query := sql.NewQuery(sql.GetRandomSpecimensQuery)
	limit := c.QueryParam(QP_LIMIT)
	randomSpecimen, err := repository.Query[model.Specimen](c.Request().Context(), h.db, query.GetQueryString(), limit)
	if err != nil {
		logger.Errorf("Can not GetRandomSamples: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve random data sample")
	}
	response := model.SpecimenResponse{
		NumItems: len(randomSpecimen),
		Data:     randomSpecimen,
	}
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
// @Success     200    {object} model.GeoAgeResponse
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

	query := sql.NewQuery(sql.GetGeoAgesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	geoAges, err := repository.Query[model.GeoAge](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetGeoAges: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geological age data")
	}
	response := model.GeoAgeResponse{
		NumItems: len(geoAges),
		Data:     geoAges,
	}
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
// @Success     200    {object} model.GeoAgePrefixResponse
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

	query := sql.NewQuery(sql.GetGeoAgePrefixesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	geoAgePrefixes, err := repository.Query[model.GeoAgePrefix](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetGeoAgePrefixes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve geological age prefix data")
	}
	response := model.GeoAgePrefixResponse{
		NumItems: len(geoAgePrefixes),
		Data:     geoAgePrefixes,
	}
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
// @Success     200    {object} model.OrganizationResponse
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

	query := sql.NewQuery(sql.GetOrganizationNamesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	organizations, err := repository.Query[model.Organization](c.Request().Context(), h.db, query.GetQueryString())
	if err != nil {
		logger.Errorf("Can not GetOrganizationNames: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve organization name data")
	}
	response := model.OrganizationResponse{
		NumItems: len(organizations),
		Data:     organizations,
	}
	return c.JSON(http.StatusOK, response)
}

// buildSampleFilterQuery constructs a query using filter params from the request
func buildSampleFilterQuery(c echo.Context, coordData map[string]interface{}, kwargs map[string]interface{}) (*sql.Query, error) {
	var returnRockClass, returnRockType bool
	for k, v := range kwargs {
		switch k {
		case KEY_ROCKCLASS:
			returnRockClass = v.(bool)
		case KEY_ROCKTYPE:
			returnRockType = v.(bool)
		}
	}

	query := sql.NewQuery(sql.GetSamplingfeatureIdsByFilterBaseQuery)
	bbox := coordData[KEY_BBOX]
	if bbox != nil {
		query = sql.NewQuery("")
		factor := coordData[KEY_TRANSLATION_FACTOR].(float64)
		params := map[string]interface{}{
			"translationFactor": -factor,
		}
		query.AddSQLBlockParametrized(sql.GetSamplingFeatureIdsByFilteBaseQueryTranslated, params)
	}

	// add optional search filters
	junctor := sql.OpWhere // junctor to connect a new filter clause to the query: can be "WHERE" or "AND/OR"

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
			query.AddFilter("mat.material", material, opMat, junctor)
			junctor = sql.OpAnd
		}
		if incType != "" {
			query.AddFilter("inctype.inclusion_type", incType, opIncType, junctor)
			junctor = sql.OpAnd
		}
		if sampTech != "" {
			query.AddFilter("stech.sampling_technique", sampTech, opSampTech, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterAnnotationsEnd)
	}

	// location filters
	junctor = sql.OpWhere // reset junctor for new subquery
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
			query.AddFilter("gs.settingname", setting, opSetting, junctor)
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
	hostMaterial, opHostMaterial, err := parseParam(c.QueryParam(QP_HOSTMAT))
	if err != nil {
		return nil, err
	}
	inclMaterial, opInclMaterial, err := parseParam(c.QueryParam(QP_INCLUSIONMAT))
	if err != nil {
		return nil, err
	}
	if returnRockType || returnRockClass || rockType != "" || rockClass != "" || mineral != "" || hostMaterial != "" || inclMaterial != "" {
		// add query module taxonomic classifiers
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersStart)
		// add filter for each subquery for significant speedup
		if returnRockType || rockType != "" {
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersRockTypeStart)
			if rockType != "" {
				query.AddFilter("tax_type.taxonomicclassifiername", rockType, opRType, sql.OpWhere)
			}
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersRockTypeEnd)
		}
		if returnRockClass || rockClass != "" {
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersRockClassStart)
			if rockClass != "" {
				query.AddFilter("tax_class.taxonomicclassifiername", rockClass, opRClass, sql.OpWhere)
			}
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersRockClassEnd)
		}
		if mineral != "" {
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersMineralStart)
			query.AddFilter("tax_min.taxonomicclassifiername", mineral, opMin, sql.OpWhere)
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersMineralEnd)
		}
		if hostMaterial != "" {
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersHostMatStart)
			query.AddFilter("tax_host.taxonomicclassifiername", hostMaterial, opHostMaterial, sql.OpAnd)
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersHostMatEnd)
		}
		if inclMaterial != "" {
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersIncMatStart)
			query.AddFilter("tax_inc.taxonomicclassifiername", inclMaterial, opInclMaterial, sql.OpAnd)
			query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersIncMatEnd)
		}
		// add taxonomic classifiers filters at the end
		junctor = sql.OpWhere
		if rockType != "" {
			query.AddFilter("rt.rock_type", rockType, opRType, junctor)
			junctor = sql.OpAnd
		}
		if rockClass != "" {
			query.AddFilter("rc.rock_class", rockClass, opRClass, junctor)
			junctor = sql.OpAnd
		}
		if mineral != "" {
			query.AddFilter("min.mineral", mineral, opMin, junctor)
			junctor = sql.OpAnd
		}
		if hostMaterial != "" {
			query.AddFilter("hostmat.host_material", hostMaterial, opHostMaterial, junctor)
			junctor = sql.OpAnd
		}
		if inclMaterial != "" {
			query.AddFilter("incmat.inclusion_material", inclMaterial, opInclMaterial, junctor)
		}
		query.AddSQLBlock(sql.GetSamplingfeatureIdsByFilterTaxonomicClassifiersEnd)
	}

	// results
	junctor = sql.OpWhere // reset junctor for new subquery
	qryString := c.QueryParam(QP_CHEMISTRY)
	if qryString != "" {
		chemQry, err := parseChemQuery(qryString)
		if err != nil {
			return nil, err
		}
		// add query module results
		mvIDList := ""
		for i := range chemQry.Expressions {
			if i > 0 {
				mvIDList += ","
			}
			mvIDList += fmt.Sprintf("m%d.samplingfeatureid", i+1)
		}
		query.AddSQLBlock(fmt.Sprintf("%s%s%s", sql.GetSamplingfeatureIdsByFilterResultsStartPre, mvIDList, sql.GetSamplingfeatureIdsByFilterResultsStartPost))
		// add ResultFilterExpression for each expression in the chemQry
		for i, expr := range chemQry.Expressions {
			// interpret missing minValue as 0 to enable "element exists"-queries
			if expr.MinValue == "" {
				expr.MinValue = "0"
			}
			junctor = sql.OpWhere // reset junctor for new expression
			exprJunctor, exists := model.CQSQLMap[expr.Junctor]
			if !exists {
				return nil, fmt.Errorf("Invalid junctor in chemical query")
			}
			query.AddSQLBlock(exprJunctor + sql.GetSamplingfeatureIdsByFilterResultsExpression)
			if expr.Type != "" {
				query.AddFilter("mv.variabletypecode", expr.Type, sql.OpEq, junctor)
				junctor = sql.OpAnd
			}
			if expr.Element != "" {
				query.AddFilter("mv.variablecode", expr.Element, sql.OpEq, junctor)
				junctor = sql.OpAnd
			}
			if expr.MinValue != "" {
				query.AddFilter("mv.datavalue", expr.MinValue, sql.OpGte, junctor)
				junctor = sql.OpAnd
			}
			if expr.MaxValue != "" {
				query.AddFilter("mv.datavalue", expr.MaxValue, sql.OpLte, junctor)
			}
			if i == 0 {
				query.AddSQLBlock(fmt.Sprintf(") m%d", i+1))
			} else {
				query.AddSQLBlock(fmt.Sprintf(") m%d on m%d.samplingfeatureid = m%d.samplingfeatureid", i+1, i+1, i))
			}
		}

		// add closing block for results
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
	polygon := coordData[KEY_POLYGON]
	if polygon != nil || bbox != nil {
		// add query module geometry
		if bbox != nil {
			bboxSlice := bbox.([][]float64)
			// format bbox string for postGIS/SQL syntax
			bboxFormatted, err := geometry.FormatPolygonArray(bboxSlice)
			if err != nil {
				return nil, err
			}
			params := map[string]interface{}{
				"bboxPolygon": fmt.Sprintf("POLYGON(%s)", bboxFormatted),
			}
			query.AddSQLBlockParametrized(sql.GestSamplingfeatureIdsByFilterGeometryBBOXStart, params)
			boundary := coordData[KEY_BOUNDARY].(float64)
			translationFactor := coordData[KEY_TRANSLATION_FACTOR].(float64)
			query.AddInTranslatedPolygonFilter("sg.geometry", bboxFormatted, junctor, boundary, translationFactor)
			junctor = sql.OpAnd
		} else {
			// if no bbox is supplied, use filter block without bbox check
			query.AddSQLBlock(sql.GestSamplingfeatureIdsByFilterGeometryStart)
		}
		if polygon != nil {
			polygonSlice := polygon.([][]float64)
			// format polygon for postGIS/SQL syntax
			polygonFormatted, err := geometry.FormatPolygonArray(polygonSlice)
			if err != nil {
				return nil, err
			}
			boundary := coordData[KEY_BOUNDARY_POLY].(float64)
			translationFactor := coordData[KEY_TRANSLATION_FACTOR_POLY].(float64)
			query.AddInTranslatedPolygonFilter("sg.geometry", polygonFormatted, junctor, boundary, translationFactor)
		}
		query.AddSQLBlock(sql.GestSamplingfeatureIdsByFilterGeometryEnd)
	}

	// coordinates
	// add query module coordinates
	query.AddSQLBlock(sql.GetGestSamplingfeatureIdsByFilterCoordinates)

	return query, nil
}

// parseClusterToGeoJSON takes an array of model.ClusteredSample and parses it into GeoJSON
func parseClusterToGeoJSON(clusterData []model.ClusteredSample) ([]model.GeoJSONCluster, []model.GeoJSONFeature, error) {
	clusters := make([]model.GeoJSONCluster, 0, len(clusterData))
	points := []model.GeoJSONFeature{}
	for _, cluster := range clusterData {
		if len(cluster.PointStrings) <= CLUSTERING_THRESHOLD {
			for i, p := range cluster.PointStrings {
				pointGeom, err := parseGeometryString(p)
				if err != nil {
					return nil, nil, err
				}
				point := model.GeoJSONFeature{
					Type:     model.GEOJSONTYPE_FEATURE,
					Geometry: *pointGeom,
					Properties: map[string]interface{}{
						"sampleID": cluster.Samples[i],
					},
				}
				points = append(points, point)
			}
			continue
		}
		centroidGeom, err := parseGeometryString(cluster.CentroidString)
		if err != nil {
			return nil, nil, err
		}
		centroid := model.GeoJSONFeature{
			Type:     model.GEOJSONTYPE_FEATURE,
			Geometry: *centroidGeom,
			Properties: map[string]interface{}{
				"clusterID":   cluster.ClusterID,
				"clusterSize": len(cluster.Samples),
			},
		}
		convexHullGeom, err := parseGeometryString(cluster.ConvexHullString)
		if err != nil {
			return nil, nil, err
		}
		convexHull := model.GeoJSONFeature{
			Type:     model.GEOJSONTYPE_FEATURE,
			Geometry: *convexHullGeom,
		}
		geoJSONCluster := model.GeoJSONCluster{
			ClusterID:  cluster.ClusterID,
			Centroid:   centroid,
			ConvexHull: convexHull,
		}
		clusters = append(clusters, geoJSONCluster)
	}
	return clusters, points, nil
}

// parseChemQuery takes a chemistry query DSL string and parses it into a ChemQuery structure
func parseChemQuery(query string) (model.ChemQuery, error) {
	chemQuery := model.ChemQuery{}
	expressionRegex := regexp.MustCompile(`\(([\w]+)?,([\w\d]+)?,([\d\.]+)?,([\d\.]+)?\)`)
	matches := expressionRegex.FindAllStringSubmatch(query, -1)
	if len(matches) == 0 {
		return chemQuery, fmt.Errorf("Can not parse chemical query")
	}
	for i, match := range matches {
		junctor := model.CQ_JUNCTOR_AND
		if i == 0 {
			// first expression gets no junctor
			junctor = model.CQ_JUNCTOR_NONE
		}
		expr := model.CQExpression{
			Junctor:  junctor,
			Type:     match[1],
			Element:  match[2],
			MinValue: match[3],
			MaxValue: match[4],
		}
		// omit expressions without type or element as they make no sense
		if expr.Type == "" && expr.Element == "" {
			continue
		}
		chemQuery.Expressions = append(chemQuery.Expressions, expr)
	}
	return chemQuery, nil
}

var geomTypeRegexp = regexp.MustCompile(`([A-z]+)`)
var coordRegexp = regexp.MustCompile(`((-?\d+(\.\d+)?) (-?\d+(\.\d+)?))`)

func parseGeometryString(geomString string) (*model.Geometry, error) {
	geometry := model.Geometry{}
	matches := geomTypeRegexp.FindAllString(geomString, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("Can not match geometry type: %s", geomString)
	}
	// set the type
	switch matches[0] {
	case "POINT":
		geometry.Type = model.GEOJSON_GEOMETRY_POINT
	case "POLYGON":
		geometry.Type = model.GEOJSON_GEOMETRY_POLYGON
	case "LINESTRING":
		geometry.Type = model.GEOJSON_GEOMETRY_LINESTRING
	default:
		return nil, fmt.Errorf("Unexpected GeoJSON type: found %s", matches[0])
	}
	// parse the coordinates
	matches = coordRegexp.FindAllString(geomString, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("Can not match coordinates: %s", geomString)
	}
	coordinates := []interface{}{}
	for _, match := range matches {
		split := strings.Split(match, " ")
		if len(split) != 2 {
			return nil, fmt.Errorf("Invalid coordinates: %s", match)
		}
		x, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid x coordinate: %s", split[0])
		}
		y, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid y coordinate: %s", split[1])
		}
		wrapper := make([]interface{}, 0, 2)
		wrapper = append(wrapper, x)
		wrapper = append(wrapper, y)
		coordinates = append(coordinates, wrapper)
	}
	if len(matches) > 1 {
		// multiple coordinates belong to a polygon and have to wrapped in two layers of array...
		coordinates = []interface{}{coordinates}
	} else {
		// ... while a single set of coordinates belongs to a point and has NO wrapping layer
		coordinates = coordinates[0].([]interface{})
	}
	geometry.Coordinates = coordinates
	return &geometry, nil
}
