// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package middleware

import (
	b64 "encoding/base64"
	"strings"

	"github.com/labstack/echo/v4"
)

const HEADER_USER_TRACKING = "DIGIS_API_UserTracking"

// GetUserTrackMiddleware returns the middleware to add a custom user header to the request to be used by the request logger
func GetUserTrackMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessKey := c.Request().Header.Get(HEADER_ACCESS_KEY)
			if accessKey != "" {
				// decode accesskey to get encoded accesskey name
				decoded, _ := b64.StdEncoding.DecodeString(accessKey)
				decoded_s := string(decoded)
				keyname := strings.Split(decoded_s, ACCESSKEY_DELIMITER)[0]
				// add keyname to request headers
				c.Request().Header.Set(HEADER_USER_TRACKING, keyname)
			}
			return next(c)
		}
	}
}
