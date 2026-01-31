package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rubellum/pages/internal/config"
)

func setupTestProject(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "pages-builder-test")
	if err != nil {
		t.Fatal(err)
	}

	// Create directories (default input is content/)
	dirs := []string{"content", "public", "templates"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			os.RemoveAll(tmpDir)
			t.Fatal(err)
		}
	}

	// Create base template
	templateContent := `<!DOCTYPE html>
<html>
<head><title>{{if .HasTitle}}{{.Title}} | {{end}}{{.Site.Title}}</title></head>
<body>
{{- if .HasTitle}}
<h1>{{.Title}}</h1>
{{- end}}
{{- if .HasDate}}
<time>{{.DateFormatted}}</time>
{{- end}}
<div>{{.Content}}</div>
</body>
</html>`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates", "base.html"), []byte(templateContent), 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestBuilder_Build(t *testing.T) {
	tmpDir, cleanup := setupTestProject(t)
	defer cleanup()

	mdContent := `# Test Post

Hello **world**!
`
	if err := os.WriteFile(filepath.Join(tmpDir, "content", "test.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Title: "Test Site"}
	opts := &Options{
		InputDir:    filepath.Join(tmpDir, "content"),
		OutputDir:   filepath.Join(tmpDir, "public"),
		TemplateDir: filepath.Join(tmpDir, "templates"),
	}

	b := New(cfg, opts)

	if _, err := b.Build(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Check output file exists
	outputPath := filepath.Join(tmpDir, "public", "test.html")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	html := string(content)

	checks := []string{
		"<title>Test Post | Test Site</title>",
		"<h1>Test Post</h1>",
		"<strong>world</strong>",
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("Expected HTML to contain %q, got:\n%s", check, html)
		}
	}
}

func TestBuilder_NestedDirectory(t *testing.T) {
	tmpDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "content", "blog", "tech")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}

	mdContent := `# Nested Post

Nested content.
`
	if err := os.WriteFile(filepath.Join(nestedDir, "post.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Title: "Test Site"}
	opts := &Options{
		InputDir:    filepath.Join(tmpDir, "content"),
		OutputDir:   filepath.Join(tmpDir, "public"),
		TemplateDir: filepath.Join(tmpDir, "templates"),
	}

	b := New(cfg, opts)

	if _, err := b.Build(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Check nested output file exists
	outputPath := filepath.Join(tmpDir, "public", "blog", "tech", "post.html")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Nested file should be generated with correct directory structure")
	}
}

func TestBuilder_CleanHTMLFiles(t *testing.T) {
	tmpDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create existing HTML file that should be cleaned
	existingHTML := filepath.Join(tmpDir, "public", "old.html")
	if err := os.WriteFile(existingHTML, []byte("<html>old</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create CSS file that should NOT be cleaned
	cssDir := filepath.Join(tmpDir, "public", "css")
	if err := os.MkdirAll(cssDir, 0755); err != nil {
		t.Fatal(err)
	}
	cssFile := filepath.Join(cssDir, "style.css")
	if err := os.WriteFile(cssFile, []byte("body {}"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Title: "Test Site"}
	opts := &Options{
		InputDir:    filepath.Join(tmpDir, "content"),
		OutputDir:   filepath.Join(tmpDir, "public"),
		TemplateDir: filepath.Join(tmpDir, "templates"),
	}

	b := New(cfg, opts)

	if _, err := b.Build(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Old HTML should be removed
	if _, err := os.Stat(existingHTML); !os.IsNotExist(err) {
		t.Error("Old HTML file should be cleaned up")
	}

	// CSS file should remain
	if _, err := os.Stat(cssFile); os.IsNotExist(err) {
		t.Error("CSS file should not be removed during cleanup")
	}
}

func TestBuilder_TemplateNotFound(t *testing.T) {
	tmpDir, cleanup := setupTestProject(t)
	defer cleanup()

	mdContent := `# Test

Content.
`
	if err := os.WriteFile(filepath.Join(tmpDir, "content", "test.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Remove template
	os.Remove(filepath.Join(tmpDir, "templates", "base.html"))

	cfg := &config.Config{Title: "Test Site"}
	opts := &Options{
		InputDir:    filepath.Join(tmpDir, "content"),
		OutputDir:   filepath.Join(tmpDir, "public"),
		TemplateDir: filepath.Join(tmpDir, "templates"),
	}

	b := New(cfg, opts)

	_, err := b.Build()
	if err == nil {
		t.Error("Expected error for missing template")
		return
	}

	if !strings.Contains(err.Error(), "template not found") && !strings.Contains(err.Error(), "failed to load template") {
		t.Errorf("Expected template-related error, got: %v", err)
	}
}

func TestGetOutputPath(t *testing.T) {
	b := &Builder{
		options: &Options{
			InputDir:  "content",
			OutputDir: "public",
		},
	}

	tests := []struct {
		input string
		want  string
	}{
		{"content/test.md", "public/test.html"},
		{"content/blog/post.md", "public/blog/post.html"},
		{"content/blog/tech/go.md", "public/blog/tech/go.html"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := b.getOutputPath(tt.input)
			if err != nil {
				t.Errorf("getOutputPath(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("getOutputPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuilder_NoDate(t *testing.T) {
	tmpDir, cleanup := setupTestProject(t)
	defer cleanup()

	mdContent := `# No Date Post

Content without date.
`
	if err := os.WriteFile(filepath.Join(tmpDir, "content", "nodate.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Title: "Test Site"}
	opts := &Options{
		InputDir:    filepath.Join(tmpDir, "content"),
		OutputDir:   filepath.Join(tmpDir, "public"),
		TemplateDir: filepath.Join(tmpDir, "templates"),
	}

	b := New(cfg, opts)

	if _, err := b.Build(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Check output file exists
	outputPath := filepath.Join(tmpDir, "public", "nodate.html")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	html := string(content)

	if !strings.Contains(html, "<h1>No Date Post</h1>") {
		t.Error("Expected title to be present")
	}
}

func TestBuilder_NoTitle(t *testing.T) {
	tmpDir, cleanup := setupTestProject(t)
	defer cleanup()

	mdContent := `Content without title.
`
	if err := os.WriteFile(filepath.Join(tmpDir, "content", "notitle.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Title: "Test Site"}
	opts := &Options{
		InputDir:    filepath.Join(tmpDir, "content"),
		OutputDir:   filepath.Join(tmpDir, "public"),
		TemplateDir: filepath.Join(tmpDir, "templates"),
	}

	b := New(cfg, opts)

	if _, err := b.Build(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Check output file exists
	outputPath := filepath.Join(tmpDir, "public", "notitle.html")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	html := string(content)

	if !strings.Contains(html, "<title>Test Site</title>") {
		t.Errorf("Expected <title>Test Site</title>, got: %s", html)
	}
	if strings.Contains(html, "<h1>") {
		t.Error("Expected no h1 when title is not specified")
	}
}

