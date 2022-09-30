package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/src/middleware"
)

// GET /authors/:lastName
func (h *Handler) GetAuthors(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	authors, err := h.db.GetAuthorByName(c.Param("lastName"))
	if err != nil {
		logger.Errorf("Can not GetAuthorsByName: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	logger.Infof("Retrieved author data: %v", authors)
	response := struct {
		NumItems int
		Data     interface{}
	}{len(authors), authors}
	return c.JSON(http.StatusOK, response)
}
