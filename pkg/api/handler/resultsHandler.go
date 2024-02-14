// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_ELEMENTTYPE = "type"
)

// GetElements godoc
// @Summary     Retrieve chemical elements
// @Description get chemical elements
// @Security    ApiKeyAuth
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       type   query    string false "Element type"
// @Param       limit  query    int    false "limit"
// @Param       offset query    int    false "offset"
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

	elementType, opElementType, err := parseParam(c.QueryParam(QP_ELEMENTTYPE))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	query := sql.NewQuery(sql.ElementsQuery)
	if elementType != "" {
		// add filter for element type
		query.AddFilter("v.variabletypecode", elementType, opElementType, sql.OpWhere)
	}
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	elements, err := repository.Query[model.Element](c.Request().Context(), h.db, query.GetQueryString(), query.GetFilterValues()...)
	if err != nil {
		logger.Errorf("Can not GetElements: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve chemical element data")
	}
	// TODO: fill standard-units
	for _, e := range elements {
		e.Unit = "tbd"
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

	query := sql.NewQuery(sql.ElementTypesQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	elementTypes, err := repository.Query[model.ElementType](c.Request().Context(), h.db, query.GetQueryString())
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
