package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

// GetCitations godoc
// @Summary     Retrieve citations
// @Description get citations
// @Security    ApiKeyAuth
// @Tags        citations
// @Accept      json
// @Produce     json
// @Success     200      {array}  model.Citation
// @Failure     401      {object} string
// @Failure     404      {object} string
// @Failure     500      {object} string
// @Router      /queries/citations [get]
func (h *Handler) GetCitations(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	citations := []model.Citation{}
	err := h.db.Query(sql.CitationsQuery, &citations)
	if err != nil {
		logger.Errorf("Can not GetCitations: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve citation data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(citations), citations}
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
// @Success     200      {array}  model.Citation
// @Failure     401      {object} string
// @Failure     404      {object} string
// @Failure     500      {object} string
// @Router      /queries/citations/{citationID} [get]
func (h *Handler) GetCitationByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	citations := []model.Citation{}
	err := h.db.Query(sql.CitationByIDQuery, &citations, c.Param("citationID"))
	if err != nil {
		logger.Errorf("Can not GetCitationByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve citation data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(citations), citations}
	return c.JSON(http.StatusOK, response)
}
