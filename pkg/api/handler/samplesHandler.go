package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

// GetSamplesByGeoSetting godoc
// @Summary     Retrieve all samples filtered by a variety of fields
// @Description Get all samples matching the current filters
// @Description Multiple values in a single filter must be comma separated
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit      query    int    false "limit"
// @Param       offset     query    int    false "offset"
// @Param       setting    query    string false "tectonic setting"
// @Param       location1  query    string false "location level 1"
// @Param       location2  query    string false "location level 2"
// @Param       location3  query    string false "location level 3"
// @Param       samplename query    string false "samplingfeature name"
// @Param       sampletech query    string false "sampling technique"
// @Param       landorsea  query    string false "land or sea"
// @Param       rockclass  query    string false "taxonomic classifier name"
// @Param       rocktype   query    string false "rock type"
// @Param       material   query    string false "material"
// @Param       majorelem  query    string false "chemical element"
// @Success     200        {array}  model.Sample
// @Failure     401        {object} string
// @Failure     404        {object} string
// @Failure     422        {object} string
// @Failure     500        {object} string
// @Router      /queries/samples [get]
func (h *Handler) GetSamplesByGeoSetting(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	samples := []model.Sample{}
	query := sql.NewQuery(sql.GetSamplesByGeoSettingQuery)

	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)

	// add optional search filters
	setting := c.QueryParam("setting")
	if setting != "" {
		query.AddInFilterQuoted("s.setting", setting)
	}
	location1 := c.QueryParam("location1")
	if location1 != "" {
		query.AddInFilterQuoted("toplevelloc.locationname", location1)
	}
	location2 := c.QueryParam("location2")
	if location2 != "" {
		query.AddInFilterQuoted("secondlevelloc.locationname", location2)
	}
	location3 := c.QueryParam("location3")
	if location3 != "" {
		query.AddInFilterQuoted("thirdlevelloc.locationname", location3)
	}
	samplename := c.QueryParam("samplename")
	if samplename != "" {
		query.AddInFilterQuoted("sf.samplingfeaturename", samplename)
	}
	sampletech := c.QueryParam("sampletech")
	if sampletech != "" {
		query.AddInFilterQuoted("ann_samp_tech.annotationtext", sampletech)
	}
	landorsea := c.QueryParam("landorsea")
	if landorsea != "" {
		query.AddInFilterQuoted("s.sitedescription", landorsea)
	}
	rockclass := c.QueryParam("rockclass")
	if rockclass != "" {
		query.AddInFilterQuoted("tax_class.taxonomicclassifiername", rockclass)
	}
	rocktype := c.QueryParam("rocktype")
	if rocktype != "" {
		query.AddInFilterQuoted("tax_type.taxonomicclassifiername", rocktype)
	}
	material := c.QueryParam("material")
	if material != "" {
		query.AddInFilterQuoted("ann_mat.annotationtext", material)
	}
	majorelem := c.QueryParam("majorelem")
	if majorelem != "" {
		query.AddInFilterQuoted("v.variablecode", majorelem)
	}

	err = h.db.Query(query.String(), &samples)
	if err != nil {
		logger.Errorf("Can not GetSamplesByGeoSetting: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(samples), samples}
	return c.JSON(http.StatusOK, response)
}
