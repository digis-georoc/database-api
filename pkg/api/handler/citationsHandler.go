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
	QP_CITATIONID = "citationID"
)

// GetCitations godoc
// @Summary     Retrieve citations
// @Description get citations
// @Security    ApiKeyAuth
// @Tags        citations
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.CitationResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/citations [get]
func (h *Handler) GetCitations(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	citations := []model.Citation{}
	query := sql.NewQuery(sql.CitationsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &citations)
	if err != nil {
		logger.Errorf("Can not GetCitations: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve citation data")
	}
	response := model.CitationResponse{
		NumItems: len(citations),
		Data:     citations,
	}
	return c.JSON(http.StatusOK, response)
}

// GetCitationByID godoc
// @Summary     Retrieve citations by citationID
// @Description get citations by citationID
// @Security    ApiKeyAuth
// @Tags        citations
// @Accept      json
// @Produce     json
// @Param       citationID path     string true "Citation ID"
// @Success     200        {object} model.CitationResponse
// @Failure     401        {object} string
// @Failure     404        {object} string
// @Failure     500        {object} string
// @Router      /queries/citations/{citationID} [get]
func (h *Handler) GetCitationByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	citations := []model.Citation{}
	err := h.db.Query(sql.CitationByIDQuery, &citations, c.Param(QP_CITATIONID))
	if err != nil {
		logger.Errorf("Can not GetCitationByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve citation data")
	}
	num := len(citations)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}
	response := model.CitationResponse{
		NumItems: num,
		Data:     citations,
	}
	return c.JSON(http.StatusOK, response)
}
