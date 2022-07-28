package handler

import "gitlab.gwdg.de/fe/digis/database-api/src/repository"

// Handler is the core strunct holding all dependencies to handle api requests
type Handler struct {
	db repository.PostgresConnector
}

// NewHandler returns a pointer to a new Handler instance
func NewHandler(db repository.PostgresConnector) *Handler {
	return &Handler{
		db: db,
	}
}
