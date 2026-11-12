package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sync"
	texttemplate "text/template"
)

//go:embed templates/*.html templates/*.txt
var templateFS embed.FS

var (
	htmlOnce sync.Once
	htmlTmpl *template.Template
	htmlErr  error

	textOnce sync.Once
	textTmpl *texttemplate.Template
	textErr  error
)

func loadHTML() (*template.Template, error) {
	htmlOnce.Do(func() {
		htmlTmpl, htmlErr = template.ParseFS(templateFS, "templates/*.html")
	})
	return htmlTmpl, htmlErr
}

func loadText() (*texttemplate.Template, error) {
	textOnce.Do(func() {
		textTmpl, textErr = texttemplate.ParseFS(templateFS, "templates/*.txt")
	})
	return textTmpl, textErr
}

func renderHTML(name string, data any) (string, error) {
	t, err := loadHTML()
	if err != nil {
		return "", fmt.Errorf("parse html templates: %w", err)
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("execute html template %q: %w", name, err)
	}
	return buf.String(), nil
}

func renderText(name string, data any) (string, error) {
	t, err := loadText()
	if err != nil {
		return "", fmt.Errorf("parse text templates: %w", err)
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("execute text template %q: %w", name, err)
	}
	return buf.String(), nil
}
