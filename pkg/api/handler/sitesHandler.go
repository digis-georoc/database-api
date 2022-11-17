package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
)

// GET /sites
func (h *Handler) GetSites(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	sites, err := h.db.GetSites()
	if err != nil {
		logger.Errorf("Can not GetSites: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve site data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(sites), sites}
	return c.JSON(http.StatusOK, response)
}
