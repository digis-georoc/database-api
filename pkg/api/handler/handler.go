// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_LIMIT  = "limit"
	QP_OFFSET = "offset"
)

// Handler is the core strunct holding all dependencies to handle api requests
type Handler struct {
	db     repository.PostgresConnector
	config *middleware.KeycloakConfig
}

// NewHandler returns a pointer to a new Handler instance
func NewHandler(db repository.PostgresConnector, config *middleware.KeycloakConfig) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}

// Ping godoc
// @Summary     Health request to check db connection
// @Description Check connection to db
// @Tags        general
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Failure     424 {object} string
// @Router      /ping [get]
func (h *Handler) Ping(c echo.Context) error {
	err := h.db.Ping()
	if err != nil {
		return c.JSON(http.StatusFailedDependency, err.Error())
	}
	return c.JSON(http.StatusOK, "Pong")
}

// Alive godoc
// @Summary     Health request to check if api is responsive
// @Description Check connection to api
// @Tags        general
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Failure     404 {object} string
// @Router      /alive [get]
func (h *Handler) Alive(c echo.Context) error {
	return c.JSON(http.StatusOK, "I'm still alive...")
}

// Version godoc
// @Summary     Get api-version
// @Description Check current version of the api
// @Tags        general
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Failure     404 {object} string
// @Router      /version [get]
func (h *Handler) Version(c echo.Context) error {
	return c.JSON(http.StatusOK, "0.5.0")
}

// handlePaginationParams reads the pagination parameters from the request and returns them as integers
func handlePaginationParams(c echo.Context) (int, int, error) {
	var err error
	limit := c.QueryParam(QP_LIMIT)
	limVal := 0
	if limit != "" {
		limVal, err = strconv.Atoi(limit)
		if err != nil {
			return 0, 0, err
		}
	}

	offset := c.QueryParam(QP_OFFSET)
	offVal := 0
	if offset != "" {
		offVal, err = strconv.Atoi(offset)
		if err != nil {
			return 0, 0, err
		}
	}
	return limVal, offVal, nil
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
	// make operator lowercase
	operator = strings.ToLower(operator)
	// validate operator
	operator, opIsValid := sql.OperatorMap[operator]
	if !opIsValid {
		return "", "", fmt.Errorf("Invalid operator")
	}
	if operator == sql.OpLike {
		// LIKE is not supported for numeric values
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return "", "", fmt.Errorf("Operator LIKE cannot be applied to numeric value: %f", f)
		}
		// replace url-compatible wildcards with sql wildcards
		value = strings.ReplaceAll(value, "*", "%")
		value = strings.ReplaceAll(value, "?", "_")
	}
	return value, operator, nil
}
