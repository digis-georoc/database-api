package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

// GetAuthors godoc
// @Summary     Retrieve authors
// @Description get authors
// @Security    ApiKeyAuth
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       limit  query    int true "limit"
// @Param       offset query    int true "offset"
// @Success     200    {array}  model.People
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     500    {object} string
// @Router      /queries/authors [get]
func (h *Handler) GetAuthors(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	authors := []model.People{}
	query := sql.NewQuery(sql.AuthorsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.String(), &authors)
	if err != nil {
		logger.Errorf("Can not GetAuthors: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(authors), authors}
	return c.JSON(http.StatusOK, response)
}

// GetAuthorsByID godoc
// @Summary     Retrieve authors by personID
// @Description get authors by personID
// @Security    ApiKeyAuth
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       personID path     string true "Person ID"
// @Success     200      {array}  model.People
// @Failure     401      {object} string
// @Failure     404      {object} string
// @Failure     500      {object} string
// @Router      /queries/authors/{personID} [get]
func (h *Handler) GetAuthorByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	authors := []model.People{}
	query := sql.NewQuery(sql.AuthorByIDQuery)
	err := h.db.Query(query.String(), &authors, c.Param("personID"))
	if err != nil {
		logger.Errorf("Can not GetAuthorByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(authors), authors}
	return c.JSON(http.StatusOK, response)
}
