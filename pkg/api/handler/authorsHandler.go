package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /authors/:lastName
func (h *Handler) GetAuthors(c echo.Context) error {
	authors, err := h.db.GetAuthorByName(c.Param("lastName"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not retrieve author data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{len(authors), authors}
	return c.JSON(http.StatusOK, response)
}
