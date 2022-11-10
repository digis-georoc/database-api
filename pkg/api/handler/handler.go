package handler

import (
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
)

// Handler is the core strunct holding all dependencies to handle api requests
type Handler struct {
	db     repository.PostgresConnector
	config middleware.KeycloakConfig
}

// NewHandler returns a pointer to a new Handler instance
func NewHandler(db repository.PostgresConnector, config middleware.KeycloakConfig) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}
