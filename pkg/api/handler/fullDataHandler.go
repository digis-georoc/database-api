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
	QP_IDENTIFIER = "identifier"
)

// GetFullDataByID godoc
// @Summary     Retrieve full dataset by samplingfeatureid
// @Description get full dataset by samplingfeatureid
// @Security    ApiKeyAuth
// @Tags        fulldata
// @Accept      json
// @Produce     json
// @Param       samplingfeatureid path     string true "Samplingfeature identifier"
// @Success     200               {array}  model.FullData
// @Failure     401               {object} string
// @Failure     404               {object} string
// @Failure     500               {object} string
// @Router      /queries/fulldata/{samplingfeatureid} [get]
func (h *Handler) GetFullDataByID(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	fullData := []model.FullData{}
	query := sql.NewQuery(sql.FullDataByIdQuery)
	err := h.db.Query(query.String(), &fullData, c.Param(QP_IDENTIFIER))
	if err != nil {
		logger.Errorf("Can not retrieve FullDataById: %v", err)
		return c.String(http.StatusInternalServerError, "Can not retrieve full data")
	}
	response := struct {
		NumItems int
		Data     interface{}
	}{1, fullData}
	return c.JSON(http.StatusOK, response)
}
