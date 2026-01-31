package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"
)

// PageData contains all data available to templates
type PageData struct {
	Title         string
	HasTitle      bool
	Date          string
	DateFormatted string
	HasDate       bool
	Content       template.HTML
	Site          SiteData
	// RelPrefix is the relative path from this page to site root (e.g. "" or "../", "../../"). Use for href/src.
	RelPrefix string
}

// SiteData contains site-wide configuration
type SiteData struct {
	Title string
}

// Renderer handles template rendering
type Renderer struct {
	tmpl *template.Template
}

// NewRenderer creates a new template renderer
func NewRenderer(templatePath string) (*Renderer, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template not found: %s", templatePath)
		}
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Renderer{tmpl: tmpl}, nil
}

// Render renders a page with the given data
func (r *Renderer) Render(data *PageData) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}

// FormatDate formats a time.Time to ISO 8601 format
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateJapanese formats a time.Time to Japanese date format
func FormatDateJapanese(t time.Time) string {
	return fmt.Sprintf("%d年%d月%d日", t.Year(), t.Month(), t.Day())
}
