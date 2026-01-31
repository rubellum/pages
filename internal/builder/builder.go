package builder

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/rubellum/pages/internal/config"
	"github.com/rubellum/pages/internal/parser"
	tmpl "github.com/rubellum/pages/internal/template"
)

// Options contains build configuration
type Options struct {
	InputDir    string
	OutputDir   string
	TemplateDir string
}

// Builder handles the site building process
type Builder struct {
	config  *config.Config
	options *Options
}

// New creates a new Builder
func New(cfg *config.Config, opts *Options) *Builder {
	return &Builder{
		config:  cfg,
		options: opts,
	}
}

// Build executes the build process. It returns the number of pages built.
func (b *Builder) Build() (int, error) {
	renderer, err := tmpl.NewRenderer(filepath.Join(b.options.TemplateDir, "base.html"))
	if err != nil {
		return 0, fmt.Errorf("failed to load template: %w", err)
	}

	if err := b.cleanHTMLFiles(); err != nil {
		return 0, fmt.Errorf("failed to clean output directory: %w", err)
	}

	mdFiles, err := b.findMarkdownFiles()
	if err != nil {
		return 0, fmt.Errorf("failed to find markdown files: %w", err)
	}

	builtCount := 0
	for _, mdFile := range mdFiles {
		built, err := b.processFile(mdFile, renderer)
		if err != nil {
			return 0, fmt.Errorf("failed to process %s: %w", mdFile, err)
		}
		if built {
			builtCount++
		}
	}

	return builtCount, nil
}

// cleanHTMLFiles removes all .html files from the output directory
func (b *Builder) cleanHTMLFiles() error {
	return filepath.Walk(b.options.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

// findMarkdownFiles recursively finds all .md files in the input directory
func (b *Builder) findMarkdownFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(b.options.InputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("input directory not found: %s", b.options.InputDir)
			}
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// processFile converts a single markdown file to HTML
func (b *Builder) processFile(mdPath string, renderer *tmpl.Renderer) (bool, error) {
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return false, err
	}

	result, err := parser.Parse(data)
	if err != nil {
		return false, fmt.Errorf("parse %s: %w", mdPath, err)
	}

	htmlContent, err := parser.ConvertMarkdown(result.Content)
	if err != nil {
		return false, fmt.Errorf("failed to parse markdown in %s: %w", mdPath, err)
	}

	outputPath, err := b.getOutputPath(mdPath)
	if err != nil {
		return false, err
	}
	relPrefix := b.relPrefixFromOutputPath(outputPath)

	hasTitle := result.Title != ""
	pageData := &tmpl.PageData{
		Title:     result.Title,
		HasTitle:  hasTitle,
		HasDate:   false,
		Content:   template.HTML(htmlContent),
		Site:      tmpl.SiteData{Title: b.config.Title},
		RelPrefix: relPrefix,
	}

	output, err := renderer.Render(pageData)
	if err != nil {
		return false, err
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return false, err
	}

	// Write output file
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return false, err
	}

	return true, nil
}

// getOutputPath converts an input markdown path to an output HTML path.
// It returns an error if mdPath is not under InputDir (path traversal).
func (b *Builder) getOutputPath(mdPath string) (string, error) {
	relPath, err := filepath.Rel(b.options.InputDir, mdPath)
	if err != nil {
		return "", fmt.Errorf("invalid path %s: %w", mdPath, err)
	}
	if strings.Contains(filepath.ToSlash(relPath), "..") {
		return "", fmt.Errorf("path traversal not allowed: %s", mdPath)
	}
	relPath = strings.TrimSuffix(relPath, ".md") + ".html"
	return filepath.Join(b.options.OutputDir, relPath), nil
}

// relPrefixFromOutputPath returns the relative path from the output file to site root (e.g. "" or "../", "../../").
func (b *Builder) relPrefixFromOutputPath(outputPath string) string {
	rel, err := filepath.Rel(b.options.OutputDir, outputPath)
	if err != nil {
		return ""
	}
	dir := filepath.Dir(rel)
	if dir == "." {
		return ""
	}
	depth := len(strings.Split(filepath.ToSlash(dir), "/"))
	return strings.Repeat("../", depth)
}
