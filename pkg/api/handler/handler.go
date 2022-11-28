package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
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
// @Summary     Sample request
// @Description Check connection to api
// @Tags        general
// @Accept      json
// @Produce     json
// @Success     200 {object}  string
// @Failure     404 {object} string
// @Router      /ping [get]
func (h *Handler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "Pong")
}

// handlePaginationParams reads the pagination parameters from the request and returns them as integers
func handlePaginationParams(c echo.Context) (int, int, error) {
	var err error
	limit := c.QueryParam("limit")
	limVal := 0
	if limit != "" {
		limVal, err = strconv.Atoi(limit)
		if err != nil {
			return 0, 0, err
		}
	}

	offset := c.QueryParam("offset")
	offVal := 0
	if offset != "" {
		offVal, err = strconv.Atoi(offset)
		if err != nil {
			return 0, 0, err
		}
	}
	return limVal, offVal, nil
}
