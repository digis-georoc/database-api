package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
)

const (
	KEY_USER     = "username"
	KEY_PASSWORD = "password"
)

type loginRequest struct {
	Username string
	Password string
}

// POST /login
// requires keys "username" and "password" in body
func (h *Handler) Login(c echo.Context) error {
	logger, ok := c.Get(middleware.LOGGER_KEY).(middleware.APILogger)
	if !ok {
		panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(middleware.LOGGER_KEY), middleware.APILogger{}))
	}
	req := loginRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Can not bind request body. Expected %+v", req))
	}
	if req.Username == "" {
		return c.JSON(http.StatusUnauthorized, "Username is required")
	}
	if req.Password == "" {
		return c.JSON(http.StatusUnauthorized, "Password is required")
	}
	config := h.config
	client := gocloak.NewClient(config.Host)
	ctx := context.Background()
	token, err := client.Login(ctx, config.ClientID, config.ClientSecret, config.Realm, req.Username, req.Password)
	if err != nil {
		msg := fmt.Sprintf("Login failed: %v", err)
		logger.Error(msg)
		return c.JSON(http.StatusUnauthorized, msg)
	}
	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, config.ClientID, config.ClientSecret, config.Realm)
	if err != nil {
		msg := fmt.Sprintf("Token inspection failed: %v", err)
		logger.Error(msg)
		return c.JSON(http.StatusUnauthorized, msg)
	}

	if !*rptResult.Active {
		msg := "Token is not active"
		logger.Error(msg)
		return c.JSON(http.StatusUnauthorized, msg)
	}
	return c.JSON(http.StatusOK, token)
}
