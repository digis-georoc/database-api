package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/geometry"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
)

// GetSampleIDStreamed_v2 godoc
//
//	@Summary		Retrieve all samplingfeatureIDs filtered by a variety of fields, streamed as pages of results
//	@Description	Get all samplingfeatureIDs matching the current filters
//	@Description	Filter DSL syntax:
//	@Description	FIELD=OPERATOR:VALUE
//	@Description	where FIELD is one of the accepted query params; OPERATOR is one of "lt" (<), "gt" (>), "eq" (=), "in" (IN), "lk" (LIKE), "btw" (BETWEEN)
//	@Description	and VALUE is an unquoted string, integer or decimal
//	@Description	Multiple VALUEs for an "in"-filter must be comma-separated and will be interpreted as a discunctive filter.
//	@Description	The OPERATORs "lt", "gt" and "btw" are only applicable to numerical values.
//	@Description	The OPERATOR "lk" is only applicable to string values and supports wildcards `*`(0 or more chars) and `?`(one char).
//	@Description	The OPERATOR "btw" accepts two comma-separated values as the inclusive lower and upper bound. Missing values are assumed as 0 and 9999999 respectively.
//	@Description	If no OPERATOR is specified, "eq" is assumed as the default OPERATOR.
//	@Description	The filters are evaluated conjunctively.
//	@Description	Note that applying more filters can slow down the query as more tables have to be considered in the evaluation.
//	@Security		ApiKeyAuth
//	@Tags			samples
//	@Accept			json
//	@Produce		json
//	@Param			limit				query		int		false	"limit"
//	@Param			offset				query		int		false	"offset"
//	@Param			setting				query		string	false	"tectonic setting - see /queries/sites/settings (supports Filter DSL)"
//	@Param			location1			query		string	false	"location level 1 - see /queries/locations/l1 (supports Filter DSL)"
//	@Param			location2			query		string	false	"location level 2 - see /queries/locations/l2 (supports Filter DSL)"
//	@Param			location3			query		string	false	"location level 3 - see /queries/locations/l3 (supports Filter DSL)"
//	@Param			latitude			query		string	false	"latitude (supports Filter DSL)"
//	@Param			longitude			query		string	false	"longitude (supports Filter DSL)"
//	@Param			rocktype			query		string	false	"rock type - see /queries/samples/rocktypes (supports Filter DSL)"
//	@Param			rockclassID			query		int		false	"taxonomic classifier ID - see /queries/samples/rockclasses value (supports Filter DSL)"
//	@Param			mineral				query		string	false	"mineral - see /queries/samples/minerals (supports Filter DSL)"
//	@Param			material			query		string	false	"material - see /queries/samples/materials (supports Filter DSL)"
//	@Param			inclusiontype		query		string	false	"inclusion type - see /queries/samples/inclusiontypes (supports Filter DSL)"
//	@Param			hostmaterial		query		string	false	"host material - see /queries/samples/hostmaterials (supports Filter DSL)"
//	@Param			inclusionmaterial	query		string	false	"inclusion material - see /queries/samples/inclusionmaterials (supports Filter DSL)"
//	@Param			sampletech			query		string	false	"sampling technique - see /queries/samples/samplingtechniques (supports Filter DSL)"
//	@Param			rimorcore			query		string	false	"rim or core - R = Rim, C = Core, I = Intermediate (supports Filter DSL)"
//	@Param			chemistry			query		string	false	"chemical filter using the form `(TYPE,ELEMENT,MIN,MAX),...` where the filter tuples are evaluated conjunctively"
//	@Param			title				query		string	false	"title of publication (supports Filter DSL)"
//	@Param			publicationyear		query		string	false	"publication year (supports Filter DSL)"
//	@Param			doi					query		string	false	"DOI (supports Filter DSL)"
//	@Param			firstname			query		string	false	"Author first name (supports Filter DSL)"
//	@Param			lastname			query		string	false	"Author last name (supports Filter DSL)"
//	@Param			agemin				query		string	false	"Specimen age min (supports Filter DSL)"
//	@Param			agemax				query		string	false	"Specimen age max (supports Filter DSL)"
//	@Param			geoage				query		string	false	"Specimen geological age - see /queries/samples/geoages (supports Filter DSL)"
//	@Param			geoageprefix		query		string	false	"Specimen geological age prefix - see /queries/samples/geoageprefixes (supports Filter DSL)"
//	@Param			lab					query		string	false	"Laboratory name - see /queries/samples/organizationnames (supports Filter DSL)"
//	@Param			polygon				query		string	false	"Coordinate-Polygon formatted as 2-dimensional json array: [[LONG,LAT],[2.4,6.3]]"
//	@Param			addcoordinates		query		bool	false	"Add coordinates to each sample"
//	@Success		206					{object}	model.SearchIndexPage
//	@Success		200					{object}	JSON
//	@Failure		401					{object}	string
//	@Failure		404					{object}	string
//	@Failure		422					{object}	string
//	@Failure		500					{object}	string
//	@Router			/v2/queries/samples [get]
func (h *Handler) GetSampleIDStreamed_v2(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	// parse filters
	filters, err := parseFilters(c)
	if err != nil {
		logger.Errorf("can not parse filters: %s", err.Error())
		return c.String(http.StatusUnprocessableEntity, "Invalid filters")
	}

	// prepare result channel and start the query concurrently
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resultChan := make(chan model.SearchIndexPage, 2)
	go h.searchIndex.QuerySortSearchAfterStream(c.Request().Context(), repository.MINIMALFIELDS, filters, resultChan)

	// stream response
	c.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(c.Response())
	totalCount := 0
	for {
		page, open := <-resultChan
		if !open {
			// channel closed because of error or finished fetching data
			break
		}
		totalCount += len(page.Documents)
		err := enc.Encode(page)
		if err != nil {
			return err
		}
		err = http.NewResponseController(c.Response()).Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

// parseFilters parses filter values from the incoming request
var skip []string = []string{"zoomlevel"}

func parseFilters(c echo.Context) (map[string]string, error) {
	filters := map[string]string{}
	for k, v := range c.QueryParams() {
		if slices.Contains(skip, k) {
			continue
		}
		filters[k] = strings.Join(v, ",")
		fmt.Printf("%s:%s\n", k, filters[k])
	}
	return filters, nil
}

// GetSamplesClustered_v2 godoc
//
//	@Summary		Retrieve all samplingfeatureIDs filtered by a variety of fields and clustered
//	@Description	Get all samplingfeatureIDs matching the current filters clustered
//	@Description	Filter DSL syntax:
//	@Description	FIELD=OPERATOR:VALUE
//	@Description	where FIELD is one of the accepted query params; OPERATOR is one of "lt" (<), "gt" (>), "eq" (=), "in" (IN), "lk" (LIKE), "btw" (BETWEEN)
//	@Description	and VALUE is an unquoted string, integer or decimal
//	@Description	Multiple VALUEs for an "in"-filter must be comma-separated and will be interpreted as a discunctive filter.
//	@Description	The OPERATORs "lt", "gt" and "btw" are only applicable to numerical values.
//	@Description	The OPERATOR "lk" is only applicable to string values and supports wildcards `*`(0 or more chars) and `?`(one char).
//	@Description	The OPERATOR "btw" accepts two comma-separated values as the inclusive lower and upper bound. Missing values are assumed as 0 and 9999999 respectively.
//	@Description	If no OPERATOR is specified, "eq" is assumed as the default OPERATOR.
//	@Description	The filters are evaluated conjunctively.
//	@Description	Note that applying more filters can slow down the query as more tables have to be considered in the evaluation.
//	@Security		ApiKeyAuth
//	@Tags			geodata
//	@Accept			json
//	@Produce		json
//	@Param			limit				query		int		false	"limit"
//	@Param			offset				query		int		false	"offset"
//	@Param			setting				query		string	false	"tectonic setting - see /queries/sites/settings (supports Filter DSL)"
//	@Param			location1			query		string	false	"location level 1 - see /queries/locations/l1 (supports Filter DSL)"
//	@Param			location2			query		string	false	"location level 2 - see /queries/locations/l2 (supports Filter DSL)"
//	@Param			location3			query		string	false	"location level 3 - see /queries/locations/l3 (supports Filter DSL)"
//	@Param			latitude			query		string	false	"latitude (supports Filter DSL)"
//	@Param			longitude			query		string	false	"longitude (supports Filter DSL)"
//	@Param			rocktype			query		string	false	"rock type - see /queries/samples/rocktypes (supports 'eq', 'in')"
//	@Param			rockclassID			query		int		false	"taxonomic classifier ID - see /queries/samples/rockclasses value (supports 'eq', 'in')"
//	@Param			mineral				query		string	false	"mineral - see /queries/samples/minerals (supports 'eq', 'in')"
//	@Param			material			query		string	false	"material - see /queries/samples/materials (supports Filter DSL)"
//	@Param			inclusiontype		query		string	false	"inclusion type - see /queries/samples/inclusiontypes (supports Filter DSL)"
//	@Param			hostmineral			query		string	false	"host mineral - see /queries/samples/hostmaterials (supports 'eq', 'in')"
//	@Param			inclusionmineral	query		string	false	"inclusion mineral - see /queries/samples/inclusionmaterials (supports 'eq', 'in')"
//	@Param			sampletech			query		string	false	"sampling technique - see /queries/samples/samplingtechniques (supports Filter DSL)"
//	@Param			rimorcore			query		string	false	"rim or core - R = Rim, C = Core, I = Intermediate (supports Filter DSL)"
//	@Param			chemistry			query		string	false	"chemical filter using the form `(TYPE,ELEMENT,MIN,MAX),...` where the filter tuples are evaluated conjunctively"
//	@Param			title				query		string	false	"title of publication (supports Filter DSL)"
//	@Param			publicationyear		query		string	false	"publication year (supports Filter DSL)"
//	@Param			doi					query		string	false	"DOI (supports Filter DSL)"
//	@Param			firstname			query		string	false	"Author first name (supports 'eq', 'in')"
//	@Param			lastname			query		string	false	"Author last name (supports 'eq', 'in')"
//	@Param			agemin				query		string	false	"Specimen age min (supports Filter DSL)"
//	@Param			agemax				query		string	false	"Specimen age max (supports Filter DSL)"
//	@Param			geoage				query		string	false	"Specimen geological age - see /queries/samples/geoages (supports Filter DSL)"
//	@Param			geoageprefix		query		string	false	"Specimen geological age prefix - see /queries/samples/geoageprefixes (supports Filter DSL)"
//	@Param			lab					query		string	false	"Laboratory name - see /queries/samples/organizationnames (supports Filter DSL)"
//	@Param			polygon				query		string	false	"Coordinate-Polygon formatted as 2-dimensional json array: [[LONG,LAT],[2.4,6.3]]"
//	@Param			bbox				query		string	true	"BoundingBox formatted as 2-dimensional json array: [[SW_Long,SW_Lat],[SE_Long,SE_Lat],[NE_Long,NE_Lat],[NW_Long,NW_Lat]]"
//	@Param			zoomlevel			query		int		false	"Zoom level of the map. Must be at least 1"
//	@Success		200					{object}	model.ClusterResponse
//	@Failure		401					{object}	string
//	@Failure		404					{object}	string
//	@Failure		422					{object}	string
//	@Failure		500					{object}	string
//	@Router			/v2/geodata/samplesclustered [get]
func (h *Handler) GetSamplesClustered_v2(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	// get zoomLevel
	zoomLevelS := c.QueryParam(QP_ZOOMLEVEL)
	if zoomLevelS == "" {
		zoomLevelS = "1"
	}
	// parse zoomLevel to int
	zoomLevel, err := strconv.Atoi(zoomLevelS)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid zoom level - must be an integer")
	}

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
	// scale bbox
	if !geometry.IsZoom0(bbox) {
		// add frame around bbox to avoid reloading on small panning
		bbox = geometry.ScaleBBox(bbox)
	}
	// truncate bbox after scaling so it contains at most one whole world
	bbox = geometry.TruncateBBox(bbox)
	// add first point again to make closed polygon shape
	bbox = append(bbox, bbox[0])

	filters, err := parseFilters(c)
	if err != nil {
		logger.Errorf("can not parse filters: %s", err.Error())
		return c.String(http.StatusUnprocessableEntity, "Invalid filters")
	}

	// start the query
	response, err := h.searchIndex.QueryClustered(repository.SEARCH_FIELDS, filters, zoomLevel)
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
	return c.JSON(http.StatusOK, response)
}
