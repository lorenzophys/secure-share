package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sethvargo/go-password/password"
)

func (app *Application) NewRouter(templatesGlob string) *echo.Echo {
	e := echo.New()
	e.Debug = app.Config.DebugMode

	e.Renderer = NewTemplateRenderer(e, templatesGlob)

	e.Static("/dist", "web/dist")

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				app.logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				app.logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())
	e.Use(app.CommonSecurityHeadersMiddleware)

	e.GET("/", func(c echo.Context) error {
		data := TemplateData{
			ProjectTitle:    app.Config.ProjectTitle,
			ProjectSubtitle: app.Config.ProjectSubtitle,
		}
		return c.Render(http.StatusOK, "base.tmpl.html", data)
	})

	e.GET("/generatepwd", func(c echo.Context) error {
		// Generate a password that is 20 characters long with 4 digits, 4 symbols,
		// allowing upper and lower case letters, allowing repeat characters.
		generatedPwd, err := password.Generate(25, 4, 4, false, true)
		if err != nil {
			app.logger.Error("failed to generate password.", "err", err)
		}
		return c.String(http.StatusOK, generatedPwd)
	})

	e.POST("/", func(c echo.Context) error {
		data := TemplateData{
			ProjectTitle:    app.Config.ProjectTitle,
			ProjectSubtitle: app.Config.ProjectSubtitle,
			BaseUrl:         fmt.Sprintf("http://%s", app.Config.BaseUrl),
		}
		if app.Config.TLS.Enabled {
			data.BaseUrl = fmt.Sprintf("https://%s", app.Config.BaseUrl)
		}

		urlKey, err := app.HandlePostSecret(c)
		if err != nil {
			return err
		}

		data.UrlKey = urlKey

		return c.Render(http.StatusOK, "link.tmpl.html", data)
	})

	e.GET("/:key", func(c echo.Context) error {
		data := TemplateData{
			ProjectTitle:    app.Config.ProjectTitle,
			ProjectSubtitle: app.Config.ProjectSubtitle,
			RenderSecret:    false,
		}

		secret, err := app.HandleGetSecret(c)
		if err != nil {
			data.Error = true
			return c.Render(http.StatusNotFound, "base.tmpl.html", data)
		}

		data.RenderSecret = true
		data.Secret = secret

		return c.Render(http.StatusOK, "base.tmpl.html", data)
	})

	app.logger.Info("setup new router", "debug_mode", e.Debug)

	return e
}

func (app *Application) HandlePostSecret(c echo.Context) (string, error) {
	textareaContent := c.FormValue("textareaContent")
	menuTimeSelection := c.FormValue("menuSelection")

	ttl, err := time.ParseDuration(menuTimeSelection)
	if err != nil {
		return "", err
	}

	return app.Store.Set(textareaContent, ttl), nil
}

func (app *Application) HandleGetSecret(c echo.Context) (string, error) {
	key := c.Param("key")
	password, ok := app.Store.Get(key)
	if !ok {
		return "", echo.NewHTTPError(http.StatusNotFound, "secret not found")
	}
	return password, nil
}
