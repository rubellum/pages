package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTitle string
		wantDate  string
		wantBody  string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid frontmatter",
			input: `---
title: "Hello World"
date: 2025-01-15
---

This is the content.`,
			wantTitle: "Hello World",
			wantDate:  "2025-01-15",
			wantBody:  "This is the content.",
			wantErr:   false,
		},
		{
			name:    "missing opening delimiter",
			input:   `title: "Test"`,
			wantErr: true,
			errMsg:  "missing opening delimiter",
		},
		{
			name: "missing closing delimiter",
			input: `---
title: "Test"
date: 2025-01-15`,
			wantErr: true,
			errMsg:  "missing closing delimiter",
		},
		{
			name: "optional title (no title)",
			input: `---
date: 2025-01-15
---

Content without title.`,
			wantTitle: "",
			wantDate:  "2025-01-15",
			wantBody:  "Content without title.",
			wantErr:   false,
		},
		{
			name: "optional date (no date)",
			input: `---
title: "No Date Post"
---

Content without date.`,
			wantTitle: "No Date Post",
			wantDate:  "",
			wantBody:  "Content without date.",
			wantErr:   false,
		},
		{
			name: "invalid yaml",
			input: `---
title: "unclosed
date: 2025-01-15
---

Content.`,
			wantErr: true,
			errMsg:  "failed to parse frontmatter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

		if result.Frontmatter.Title != tt.wantTitle {
			t.Errorf("title: got %q, want %q", result.Frontmatter.Title, tt.wantTitle)
		}

		if tt.wantDate == "" {
			if !result.Frontmatter.Date.IsZero() {
				t.Errorf("date: expected zero value, got %v", result.Frontmatter.Date)
			}
		} else {
			expectedDate, _ := time.Parse("2006-01-02", tt.wantDate)
			if !result.Frontmatter.Date.Equal(expectedDate) {
				t.Errorf("date: got %v, want %v", result.Frontmatter.Date, expectedDate)
			}
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

