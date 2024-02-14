// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"
)

const (
	HEADER_ACCESS_KEY   = "DIGIS-API-ACCESSKEY"
	ACCESSKEY_DELIMITER = ":"
)

// GetAccessKeyMiddleware returns the middleware to validate authentication via a predefined access key
// param secStore: the secretstore.Secretstore to load the permitted access keys from vault injected secret
func GetAccessKeyMiddleware(secStore secretstore.SecretStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessKey := c.Request().Header.Get(HEADER_ACCESS_KEY)
			if accessKey == "" {
				return c.JSON(http.StatusUnauthorized, "No access key provided")
			}
			logger, ok := c.Get(LOGGER_KEY).(APILogger)
			if !ok {
				panic(fmt.Sprintf("Can not get context.logger of type %T as type %T", c.Get(LOGGER_KEY), APILogger{}))
			}
			err := secStore.LoadSecretsFromFile("/vault/secrets/accesskeys.txt")
			if err != nil {
				logger.Errorf("Can not load secrets from file: %v", err)
				return c.JSON(http.StatusInternalServerError, "Can not verify allowed access keys")
			}
			allowedKeys, err := secStore.GetMap()
			if err != nil {
				logger.Errorf("No allowed access keys found: %v", err)
				return c.JSON(http.StatusInternalServerError, "Can not verify allowed access keys: none configured")
			}
			for _, v := range allowedKeys {
				if v == accessKey {
					return next(c)
				}
			}
			return c.JSON(http.StatusUnauthorized, "Invalid access key")
		}
	}
}
