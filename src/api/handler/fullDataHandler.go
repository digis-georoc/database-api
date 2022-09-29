package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/src/middleware"
)

// GET /fullData/:identifier
func (h *Handler) GetFullData(c echo.Context) error {
	fullData, err := h.db.GetFullDataById(c.Param("identifier"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Can not retrieve full data: %v", err))
	}
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	logger.Infof("Retrieved full data: %v", fullData)
	response := struct {
		NumItems int
		Data     interface{}
	}{1, fullData}
	return c.JSON(http.StatusOK, response)
}
