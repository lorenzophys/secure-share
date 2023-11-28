package main

import (
	"github.com/labstack/echo/v4"
)

func CommonSecurityHeadersMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "0")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")

		return next(c)
	}
}
