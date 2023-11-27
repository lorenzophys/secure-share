package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lorenzophys/secure_share/internal/store"
	memoryStore "github.com/lorenzophys/secure_share/internal/store/in-memory"
)

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

func main() {
	var store store.SecretStore

	storeType := "in-memory"

	switch storeType {
	case "in-memory":
		store = memoryStore.NewMemoryStore()
	default:
		log.Fatal("Invalid storage type")
	}

	e := echo.New()
	e.Debug = true

	templatesRenderer := NewTemplateRenderer(e, "web/templates/*.tmpl.html")
	e.Renderer = templatesRenderer

	e.Static("/dist", "web/dist")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		data := map[string]string{
			"title": "Home",
		}
		return c.Render(http.StatusOK, "base.tmpl.html", data)
	})

	e.POST("/", func(c echo.Context) error {
		textareaContent := c.FormValue("textareaContent")
		//menuSelection := c.FormValue("menuSelection")
		data := map[string]string{
			"urlKey": store.Set(textareaContent),
		}
		return c.Render(http.StatusOK, "link.tmpl.html", data)
	})

	e.GET("/:key", func(c echo.Context) error {
		key := c.Param("key")
		password, ok := store.Get(key)
		if !ok {
			return c.String(http.StatusNotFound, "Password not found")
		}
		return c.String(http.StatusOK, password)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
