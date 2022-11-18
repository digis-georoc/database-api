package handler

import (
	"net/http"

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
// @Success     200 {array}  string
// @Failure     404 {object} string
// @Router      /ping [get]
func (h *Handler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "Pong")
}
