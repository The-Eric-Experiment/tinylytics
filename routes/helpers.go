package routes

import (
	"html/template"
)

// LoadTemplates loads all HTML templates
func LoadTemplates() *template.Template {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return tmpl
}

