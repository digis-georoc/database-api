package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

// GetElements godoc
// @Summary     Retrieve chemical elements
// @Description get chemical elements
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.ElementResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/results/elements [get]
func (h *Handler) GetElements(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	elements := []model.Element{}
	query := sql.NewQuery(sql.ElementsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &elements)
	if err != nil {
		logger.Errorf("Can not GetElements: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve chemical element data")
	}
	response := model.ElementResponse{
		NumItems: len(elements),
		Data:     elements,
	}
	return c.JSON(http.StatusOK, response)
}

// GetElementTypes godoc
// @Summary     Retrieve chemical element types
// @Description get chemical element types
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.ElementTypeResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/results/elementtypes [get]
func (h *Handler) GetElementTypes(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	elementTypes := []model.ElementType{}
	query := sql.NewQuery(sql.ElementTypesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &elementTypes)
	if err != nil {
		logger.Errorf("Can not GetElementTypes: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve chemical element type data")
	}
	response := model.ElementTypeResponse{
		NumItems: len(elementTypes),
		Data:     elementTypes,
	}
	return c.JSON(http.StatusOK, response)
}
