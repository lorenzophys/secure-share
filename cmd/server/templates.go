package main

import (
	"html/template"
	"io"

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
	tmpl, err := template.ParseGlob(globPattern)
	if err != nil {
		e.Logger.Fatal(err)
	}

	return &Template{Templates: tmpl}
}
