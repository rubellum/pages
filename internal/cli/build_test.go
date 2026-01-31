package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSiteTitleFromIndex(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pages-build-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("index.md with title", func(t *testing.T) {
		content := `---
title: "My Blog"
date: 2025-01-01
---

Hello.
`
		if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ResolveSiteTitleFromIndex(tmpDir)
		if got != "My Blog" {
			t.Errorf("expected title \"My Blog\", got %q", got)
		}
	})

	t.Run("no index.md", func(t *testing.T) {
		emptyDir, err := os.MkdirTemp("", "pages-build-empty")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(emptyDir)
		got := ResolveSiteTitleFromIndex(emptyDir)
		if got != "" {
			t.Errorf("expected empty title when no index.md, got %q", got)
		}
	})

	t.Run("index.md without title in frontmatter", func(t *testing.T) {
		subDir := filepath.Join(tmpDir, "sub")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}
		content := `---
date: 2025-01-01
---

No title.
`
		if err := os.WriteFile(filepath.Join(subDir, "index.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ResolveSiteTitleFromIndex(subDir)
		if got != "" {
			t.Errorf("expected empty title when frontmatter has no title, got %q", got)
		}
	})
}
