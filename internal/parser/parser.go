package parser

import (
	"bytes"
	"fmt"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

var (
	frontmatterDelimiter = []byte("---")
	md                    = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
)

// Frontmatter represents the YAML frontmatter of a Markdown file
type Frontmatter struct {
	Title string    `yaml:"title"`
	Date  time.Time `yaml:"date"`
}

// ParseResult contains the parsed frontmatter and content
type ParseResult struct {
	Frontmatter *Frontmatter
	Content     []byte
}

// Parse extracts frontmatter and content from a Markdown file
func Parse(data []byte) (*ParseResult, error) {
	data = bytes.TrimSpace(data)
	if !bytes.HasPrefix(data, frontmatterDelimiter) {
		return nil, fmt.Errorf("invalid frontmatter: missing opening delimiter")
	}
	rest := data[len(frontmatterDelimiter):]
	idx := bytes.Index(rest, frontmatterDelimiter)
	if idx == -1 {
		return nil, fmt.Errorf("invalid frontmatter: missing closing delimiter")
	}
	fmData := bytes.TrimSpace(rest[:idx])
	content := bytes.TrimSpace(rest[idx+len(frontmatterDelimiter):])

	var fm Frontmatter
	if err := yaml.Unmarshal(fmData, &fm); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}
	return &ParseResult{Frontmatter: &fm, Content: content}, nil
}

// ConvertMarkdown converts Markdown content to HTML
func ConvertMarkdown(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
