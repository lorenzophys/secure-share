package main

import (
	"github.com/labstack/echo/v4"
)

func (app *Application) CommonSecurityHeadersMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "0")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("Referrer-Policy", "origin-when-cross-origin")

		if app.Config.TLS.Enabled {
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		return next(c)
	}
}
