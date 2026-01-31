package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDefaultFlags(t *testing.T) {
	inputDef := buildCmd.Flags().Lookup("input").DefValue
	if inputDef != "content" {
		t.Errorf("build -i default = %q, want content", inputDef)
	}
	outputDef := buildCmd.Flags().Lookup("output").DefValue
	if outputDef != "public" {
		t.Errorf("build -o default = %q, want public", outputDef)
	}
}

func TestE2E_InitThenBuild(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pages-e2e")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := runInit(nil, nil); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	inputDir = "content"
	outputDir = "public"
	defer func() {
		inputDir = "content"
		outputDir = "public"
	}()

	if err := runBuild(nil, nil); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	indexHTML := filepath.Join(tmpDir, "public", "index.html")
	body, err := os.ReadFile(indexHTML)
	if err != nil {
		t.Fatalf("public/index.html not found or unreadable: %v", err)
	}
	html := string(body)
	if !strings.Contains(html, "My Site") {
		t.Errorf("expected HTML to contain site title 'My Site', got:\n%s", html)
	}
	if !strings.Contains(html, "Welcome") {
		t.Errorf("expected HTML to contain body 'Welcome', got:\n%s", html)
	}
}

func TestResolveSiteTitleFromIndex(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pages-build-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("index.md with title", func(t *testing.T) {
		content := `# My Blog

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

	t.Run("index.md without h1 title", func(t *testing.T) {
		subDir := filepath.Join(tmpDir, "sub")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}
		content := `No title.
`
		if err := os.WriteFile(filepath.Join(subDir, "index.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ResolveSiteTitleFromIndex(subDir)
		if got != "" {
			t.Errorf("expected empty title when no # line, got %q", got)
		}
	})
}
