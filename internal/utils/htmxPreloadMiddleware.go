package utils

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func HtmxCacheMiddleware(cacheTimeSeconds uint) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("HX-Request") == "true" {
				c.Response().Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d", cacheTimeSeconds))
				c.Response().Header().Set("Vary", "HX-Request, Accept")
			}

			return next(c)
		}
	}
}
