package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}
	if cfg.Title != "" {
		t.Errorf("expected default Title \"\", got %q", cfg.Title)
	}
}

func TestLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pages-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name      string
		content   string
		wantErr   bool
		errMsg    string
		wantTitle string
	}{
		{
			name:      "valid config",
			content:   `title: "My Site"` + "\n",
			wantErr:   false,
			wantTitle: "My Site",
		},
		{
			name:      "empty config (defaults applied)",
			content:   ``,
			wantErr:   false,
			wantTitle: "",
		},
		{
			name:    "invalid yaml",
			content: `title: "unclosed string`,
			wantErr: true,
			errMsg:  "failed to parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write config file
			configPath := filepath.Join(tmpDir, "pages.yaml")
			if err := os.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			cfg, err := Load(configPath)

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

			if cfg.Title != tt.wantTitle {
				t.Errorf("expected title %q, got %q", tt.wantTitle, cfg.Title)
			}
		})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	cfg, err := Load("/nonexistent/pages.yaml")
	if err != nil {
		t.Errorf("expected no error for nonexistent file (return default), got %v", err)
		return
	}
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}
	if cfg.Title != "" {
		t.Errorf("expected default config (empty title), got title=%q", cfg.Title)
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := &Config{Title: "Custom"}
	cfg.ApplyDefaults()
	if cfg.Title != "Custom" {
		t.Errorf("expected Title unchanged \"Custom\", got %q", cfg.Title)
	}
}

