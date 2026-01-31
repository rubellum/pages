package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTitle string
		wantBody  string
	}{
		{
			name:      "first line is h1",
			input:     "# Hello World\n\nThis is the content.",
			wantTitle: "Hello World",
			wantBody:  "This is the content.",
		},
		{
			name:      "no h1, full content",
			input:     "This is the content.",
			wantTitle: "",
			wantBody:  "This is the content.",
		},
		{
			name:      "h1 with extra spaces",
			input:     "#  No Date Post  \n\nContent without date.",
			wantTitle: "No Date Post",
			wantBody:  "Content without date.",
		},
		{
			name:      "empty file",
			input:     "",
			wantTitle: "",
			wantBody:  "",
		},
		{
			name:      "only h1",
			input:     "# Only Title",
			wantTitle: "Only Title",
			wantBody:  "",
		},
		{
			name:      "h2 first, no title",
			input:     "## Not H1\n\nBody.",
			wantTitle: "",
			wantBody:  "## Not H1\n\nBody.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse([]byte(tt.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result.Title != tt.wantTitle {
				t.Errorf("title: got %q, want %q", result.Title, tt.wantTitle)
			}
			content := strings.TrimSpace(string(result.Content))
			if content != tt.wantBody {
				t.Errorf("content: got %q, want %q", content, tt.wantBody)
			}
		})
	}
}

func TestConvertMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "paragraph",
			input:    "Hello, world!",
			contains: []string{"<p>Hello, world!</p>"},
		},
		{
			name:     "heading",
			input:    "# Heading 1\n\n## Heading 2",
			contains: []string{"<h1>Heading 1</h1>", "<h2>Heading 2</h2>"},
		},
		{
			name:     "bold and italic",
			input:    "**bold** and *italic*",
			contains: []string{"<strong>bold</strong>", "<em>italic</em>"},
		},
		{
			name:     "link",
			input:    "[link](https://example.com)",
			contains: []string{`<a href="https://example.com">link</a>`},
		},
		{
			name:     "code block",
			input:    "```go\nfmt.Println(\"Hello\")\n```",
			contains: []string{"<pre>", "<code", "fmt.Println"},
		},
		{
			name:     "inline code",
			input:    "Use `code` here",
			contains: []string{"<code>code</code>"},
		},
		{
			name:     "unordered list",
			input:    "- Item 1\n- Item 2",
			contains: []string{"<ul>", "<li>Item 1</li>", "<li>Item 2</li>"},
		},
		{
			name:     "ordered list",
			input:    "1. First\n2. Second",
			contains: []string{"<ol>", "<li>First</li>", "<li>Second</li>"},
		},
		{
			name:     "blockquote",
			input:    "> Quote here",
			contains: []string{"<blockquote>", "Quote here"},
		},
		{
			name:     "horizontal rule",
			input:    "Above\n\n---\n\nBelow",
			contains: []string{"<hr"},
		},
		{
			name:     "table (GFM)",
			input:    "| A | B |\n|---|---|\n| 1 | 2 |",
			contains: []string{"<table>", "<th>A</th>", "<td>1</td>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertMarkdown([]byte(tt.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			html := string(result)
			for _, want := range tt.contains {
				if !strings.Contains(html, want) {
					t.Errorf("expected HTML to contain %q, got:\n%s", want, html)
				}
			}
		})
	}
}
