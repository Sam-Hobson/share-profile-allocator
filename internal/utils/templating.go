package utils

import (
	"bytes"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// dict function to create a map in templates
func dict(values ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key, _ := values[i].(string)
		if i+1 < len(values) {
			m[key] = values[i+1]
		}
	}
	return m
}

var funcMap = template.FuncMap{
	"dict": dict,
	"add1": func(i int) int { return i + 1 },
}

var Template = &templateRenderer{
	templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("public/views/*.html")),
}

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func RenderTemplateToString(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := Template.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
