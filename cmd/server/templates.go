package main

import (
	"html/template"
	"io"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

const templatesPath = "web/templates/*.tmpl.html"

type TemplateData struct {
	ProjectTitle    string
	ProjectSubtitle string
	BaseUrl         string
	UrlKey          string
	RenderSecret    bool
	Secret          string
	Error           bool
}

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func NewTemplateRenderer(e *echo.Echo, globPattern string) echo.Renderer {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tmpl, err := template.ParseGlob(globPattern)
	if err != nil {
		logger.Error("error parsing template directory glob.", "glob_pattern", globPattern, "error", err)
		os.Exit(1)
	}

	return &Template{Templates: tmpl}
}
