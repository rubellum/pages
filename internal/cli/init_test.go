package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesContentDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pages-init-test")
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
		t.Fatalf("runInit failed: %v", err)
	}

	// Must create content/ (not src/)
	contentDir := filepath.Join(tmpDir, "content")
	if st, err := os.Stat(contentDir); err != nil || !st.IsDir() {
		t.Errorf("expected content/ directory to exist, err=%v", err)
	}
	indexMD := filepath.Join(contentDir, "index.md")
	if _, err := os.Stat(indexMD); err != nil {
		t.Errorf("expected content/index.md to exist, err=%v", err)
	}
	templatesBase := filepath.Join(tmpDir, "templates", "base.html")
	if _, err := os.Stat(templatesBase); err != nil {
		t.Errorf("expected templates/base.html to exist, err=%v", err)
	}
	publicCSS := filepath.Join(tmpDir, "public", "css", "style.css")
	if _, err := os.Stat(publicCSS); err != nil {
		t.Errorf("expected public/css/style.css to exist, err=%v", err)
	}

	// Must NOT create src/
	srcDir := filepath.Join(tmpDir, "src")
	if st, err := os.Stat(srcDir); err == nil && st.IsDir() {
		t.Error("expected src/ not to be created (default is content/)")
	}
}
