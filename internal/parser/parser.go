package parser

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

// ParseResult contains the extracted title and body from a Markdown file
type ParseResult struct {
	Title   string
	Content []byte
}

// Parse reads markdown and extracts title from the first # line; the rest is content.
// If there is no # line, Title is empty and Content is the whole file.
func Parse(data []byte) (*ParseResult, error) {
	data = bytes.TrimSpace(data)
	lines := strings.Split(string(data), "\n")

	var title string
	var contentStart int

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			title = strings.TrimSpace(trimmed[2:])
			contentStart = i + 1
			break
		}
		if trimmed != "" {
			// First non-empty line is not a heading; no title
			contentStart = 0
			break
		}
	}

	var content []byte
	if contentStart > 0 {
		content = []byte(strings.TrimSpace(strings.Join(lines[contentStart:], "\n")))
	} else {
		content = data
	}

	return &ParseResult{Title: title, Content: content}, nil
}

// ConvertMarkdown converts Markdown content to HTML
func ConvertMarkdown(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
