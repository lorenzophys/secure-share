package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *Application) NewRouter(templatesGlob string) *echo.Echo {
	e := echo.New()
	e.Debug = app.Config.DebugMode

	e.Renderer = NewTemplateRenderer(e, templatesGlob)

	e.Static("/dist", "web/dist")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		data := TemplateData{
			ProjectTitle:    app.Config.ProjectTitle,
			ProjectSubtitle: app.Config.ProjectSubtitle,
		}
		return c.Render(http.StatusOK, "base.tmpl.html", data)
	})

	e.POST("/", func(c echo.Context) error {
		data := TemplateData{
			ProjectTitle:    app.Config.ProjectTitle,
			ProjectSubtitle: app.Config.ProjectSubtitle,
			BaseUrl:         app.Config.BaseUrl,
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
