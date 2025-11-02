package main

import (
	"embed"
	"net/http"
	"text/template"
)

type Template struct {
	templates *template.Template
}

//go:embed templates
var fs embed.FS

func NewTemplate() (*Template, error) {
	templates, err := template.ParseFS(fs, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Template{templates}, nil
}

func (t *Template) Render(w http.ResponseWriter, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name+".html", data)
}
