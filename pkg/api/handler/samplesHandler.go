package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_SAMPLINGFEATUREID = "samplingfeatureID"
	QP_SETTING           = "setting"
	QP_LOC1              = "location1"
	QP_LOC2              = "location2"
	QP_LOC3              = "location3"
	QP_SAMPNAME          = "samplename"
	QP_SAMPTECH          = "sampletech"
	QP_LORS              = "landorsea"
	QP_ROCKCLASS         = "rockclass"
	QP_ROCKTYPE          = "rocktype"
	QP_MATERIAL          = "material"
	QP_MAJORELEM         = "majorelem"
)

// GetSamples godoc
// @Summary     Retrieve all samples
// @Description get all sample data
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {array}  model.Sample
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/samples [get]
func (h *Handler) GetSamples(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	samples := []model.Sample{}
	query := sql.NewQuery(sql.GetSamplesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.String(), &samples)
	if err != nil {
		logger.Errorf("Can not GetSamples: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve sample data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(samples), samples}
	return c.JSON(http.StatusOK, response)
}

// GetSampleByID godoc
// @Summary     Retrieve sample by samplingfeatureid
// @Description get sample by samplingfeatureid
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       samplingfeatureID path     string true "Sample ID"
// @Success     200        {array}  model.Sample
// @Failure     401        {object} string
// @Failure     404        {object} string
// @Failure     500        {object} string
// @Router      /queries/samples/{samplingfeatureID} [get]
func (h *Handler) GetSampleByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	samples := []model.Sample{}
	query := sql.NewQuery(sql.GetSampleByIDQuery)
	err := h.db.Query(query.String(), &samples, c.Param(QP_SAMPLINGFEATUREID))
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
// @Param       setting    query    string false "tectonic setting - choose from /queries/sites/settings"
// @Param       location1  query    string false "location level 1 - choose from /queries/locations/l1"
// @Param       location2  query    string false "location level 2 - choose from /queries/locations/l2"
// @Param       location3  query    string false "location level 3 - choose from /queries/locations/l3"
// @Param       samplename query    string false "samplingfeature name"
// @Param       sampletech query    string false "sampling technique - choose from /queries/samples/samplingtechniques"
// @Param       landorsea  query    string false "land or sea - choose from /queries/sites/landorsea"
// @Param       rockclass  query    string false "rock class - choose from /queries/samples/rockclasses"
// @Param       rocktype   query    string false "rock type - choose from /queries/samples/rocktypes"
// @Param       mineral    query    string false "mineral - choose from /queries/samples/minerals"
// @Param       majorelem  query    string false "chemical element"
// @Success     200        {array}  model.SampleByGeoSettingResponse
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
	samples := []model.SampleByGeoSettingResponse{}
	query := sql.NewQuery(sql.GetSamplesByGeoSettingQuery)

	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)

	// add optional search filters
	setting := c.QueryParam(QP_SETTING)
	if setting != "" {
		query.AddInFilterQuoted("s.setting", setting)
	}
	location1 := c.QueryParam(QP_LOC1)
	if location1 != "" {
		query.AddInFilterQuoted("toplevelloc.locationname", location1)
	}
	location2 := c.QueryParam(QP_LOC2)
	if location2 != "" {
		query.AddInFilterQuoted("secondlevelloc.locationname", location2)
	}
	location3 := c.QueryParam(QP_LOC3)
	if location3 != "" {
		query.AddInFilterQuoted("thirdlevelloc.locationname", location3)
	}
	samplename := c.QueryParam(QP_SAMPNAME)
	if samplename != "" {
		query.AddInFilterQuoted("sf.samplingfeaturename", samplename)
	}
	sampletech := c.QueryParam(QP_SAMPTECH)
	if sampletech != "" {
		query.AddInFilterQuoted("ann_samp_tech.annotationtext", sampletech)
	}
	landorsea := c.QueryParam(QP_LORS)
	if landorsea != "" {
		query.AddInFilterQuoted("s.sitedescription", landorsea)
	}
	rockclass := c.QueryParam(QP_ROCKCLASS)
	if rockclass != "" {
		query.AddInFilterQuoted("tax_class.taxonomicclassifiername", rockclass)
	}
	rocktype := c.QueryParam(QP_ROCKTYPE)
	if rocktype != "" {
		query.AddInFilterQuoted("tax_type.taxonomicclassifiername", rocktype)
	}
	material := c.QueryParam(QP_MATERIAL)
	if material != "" {
		query.AddInFilterQuoted("ann_mat.annotationtext", material)
	}
	majorelem := c.QueryParam(QP_MAJORELEM)
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
	err = h.db.Query(query.String(), &specimentypes)
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
	err = h.db.Query(query.String(), &samplingtechniques)
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
	err = h.db.Query(query.String(), &rockclasses)
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
	err = h.db.Query(query.String(), &rocktypes)
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
	err = h.db.Query(query.String(), &minerals)
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
