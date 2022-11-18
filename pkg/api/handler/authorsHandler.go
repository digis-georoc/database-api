package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
)

// GetAuthorsByLastname godoc
// @Summary     Retrieve authors by lastname
// @Description get authors by lastname
// @securityDefinitions.apikey ApiKeyAuth
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       lastName path     string true "Author lastname"
// @Success     200      {array}  model.People
// @Failure     401      {object} string
// @Failure     404      {object} string
// @Failure     500      {object} string
// @Router      /secured/authors/{lastName} [get]
func (h *Handler) GetAuthorsByLastName(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	authors, err := h.db.GetAuthorByLastName(c.Param("lastName"))
	if err != nil {
		logger.Errorf("Can not GetAuthorsByName: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(authors), authors}
	return c.JSON(http.StatusOK, response)
}
