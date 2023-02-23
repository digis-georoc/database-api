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

// GetSamplesFiltered godoc
// @Summary     Retrieve all samplingfeatureIDs filtered by a variety of fields
// @Description Get all samplingfeatureIDs matching the current filters
// @Description Filter DSL syntax:
// @Description ?<field>=<op>:<value>
// @Description where <field> is one of the accepted query params; <op> is one of "lt", "gt", "eq", "in" and <value> is an unquoted string, integer or decimal
// @Description Multiple values for an "in"-filter must be comma-separated and will be interpreted as a discunctive filter.
// @Description The filters are evaluated conjunctively.
// @Description Note that applying more filters can slow down the query as more tables have to be considered in the evaluation.
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit         query    int    false "limit"
// @Param       offset        query    int    false "offset"
// @Param       setting       query    string false "tectonic setting"
// @Param       location1     query    string false "location level 1"
// @Param       location2     query    string false "location level 2"
// @Param       location3     query    string false "location level 3"
// @Param       rocktype      query    string false "rock type"
// @Param       rockclass     query    string false "taxonomic classifier name"
// @Param       mineral       query    string false "mineral"
// @Param       material      query    string false "material"
// @Param       inclusiontype query    string false "inclusion type"
// @Param       sampletech    query    string false "sampling technique"
// @Param       element       query    string false "chemical element"
// @Param       elementtype   query    string false "element type"
// @Param       value         query    string false "measured value"
// @Success     200           {array}  model.Specimen
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
	specimen := []model.Specimen{}
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

// parseParam parses a given query parameter and validates the contents
func parseParam(queryParam string) (string, string, error) {
	if queryParam == "" {
		return "", "", nil
	}
	operator, value, found := strings.Cut(queryParam, ":")
	if !found {
		return "", "", fmt.Errorf("Invalid param format")
	}
	// validate operator
	operator, opIsValid := sql.OperatorMap[operator]
	if !opIsValid {
		return "", "", fmt.Errorf("Invalid operator")
	}
	return value, operator, nil
}

func (h *Handler) TestQuery(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	query := `select distinct spec.samplingfeatureid
	from odm2.specimens spec
	join (
		-- taxonomic classifiers
		select stc.samplingfeatureid
		from odm2.specimentaxonomicclassifiers stc
		left join odm2.taxonomicclassifiers tax_type on tax_type.taxonomicclassifierid = stc.taxonomicclassifierid and tax_type.taxonomicclassifiertypecv = 'Rock'
		left join odm2.taxonomicclassifiers tax_class on tax_class.taxonomicclassifierid = stc.taxonomicclassifierid and tax_class.taxonomicclassifiertypecv = 'Lithology'
		left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifierdescription  = 'Mineral Classification from GEOROC'
	WHERE tax_type.taxonomicclassifiername = $1
	) tax on tax.samplingfeatureid = spec.samplingfeatureid`
	arg := c.QueryParam("arg")
	specimen := []model.Specimen{}
	err := h.db.Query(query, &specimen, fmt.Sprintf("'%s'", arg))
	if err != nil {
		logger.Errorf("Can not TEST: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve TEST data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(specimen), specimen}
	return c.JSON(http.StatusOK, response)
}
