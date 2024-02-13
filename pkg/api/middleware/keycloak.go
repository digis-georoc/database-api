// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/labstack/echo/v4"
)

type KeycloakConfig struct {
	Host         string
	ClientID     string
	ClientSecret string
	Realm        string
}

// GetAcademicCloudAuthMW returns a echo.MiddlewareFunc
// that returns the echo.HandlerFunc
// which executes the auth middleware method
// param: config The KeycloakConfig to connect to and interact with the Keycloak instance
func GetAcademicCloudAuthMW(config *KeycloakConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config == nil {
				return fmt.Errorf("No keycloak config provided")
			}
			logger, ok := c.Get(LOGGER_KEY).(APILogger)
			if !ok {
				panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(LOGGER_KEY), APILogger{}))
			}
			err := academiccloudAuth(c, *config)
			if err != nil {
				logger.Errorf("Auth error: %v", err)
				// set Unauthorized header
				return c.String(http.StatusUnauthorized, "Invalid or expired token")
			}
			return next(c)
		}
	}
}

// academiccloudAuth connects to the configured Keycloak instance
// and retrieves a token
func academiccloudAuth(c echo.Context, config KeycloakConfig) error {
	client := gocloak.NewClient(config.Host)
	ctx := context.Background()
	token := c.Request().Header.Get("Bearer")
	rptResult, err := client.RetrospectToken(ctx, token, config.ClientID, config.ClientSecret, config.Realm)
	if err != nil {
		return fmt.Errorf("Inspection failed: %v", err)
	}

	if !*rptResult.Active {
		return fmt.Errorf("Token is not active")
	}
	return nil
}
