package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/sql"
)

const (
	QP_PERSONID = "personID"
)

// GetAuthors godoc
// @Summary     Retrieve authors
// @Description get authors
// @Security    ApiKeyAuth
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       limit  query    int false "limit"
// @Param       offset query    int false "offset"
// @Success     200    {object} model.PeopleResponse
// @Failure     401    {object} string
// @Failure     404    {object} string
// @Failure     422    {object} string
// @Failure     500    {object} string
// @Router      /queries/authors [get]
func (h *Handler) GetAuthors(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	authors := []model.Person{}
	query := sql.NewQuery(sql.AuthorsQuery)
	limit, offset, err := handlePaginationParams(c)
	if err != nil {
		logger.Errorf("Invalid pagination params: %v", err)
		return c.String(http.StatusUnprocessableEntity, "Invalid pagination parameters")
	}
	query.AddLimit(limit)
	query.AddOffset(offset)
	err = h.db.Query(query.GetQueryString(), &authors)
	if err != nil {
		logger.Errorf("Can not GetAuthors: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	response := model.PeopleResponse{
		NumItems: len(authors),
		Data:     authors,
	}
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
// @Success     200      {object} model.PeopleResponse
// @Failure     401      {object} string
// @Failure     404      {object} string
// @Failure     500      {object} string
// @Router      /queries/authors/{personID} [get]
func (h *Handler) GetAuthorByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}

	authors := []model.Person{}
	err := h.db.Query(sql.AuthorByIDQuery, &authors, c.Param(QP_PERSONID))
	if err != nil {
		logger.Errorf("Can not GetAuthorByID: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	num := len(authors)
	if num == 0 {
		return c.String(http.StatusNotFound, "No data found")
	}
	response := model.PeopleResponse{
		NumItems: num,
		Data:     authors,
	}
	return c.JSON(http.StatusOK, response)
}
