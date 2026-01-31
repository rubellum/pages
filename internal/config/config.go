package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the site configuration from pages.yaml (optional)
type Config struct {
	Title string `yaml:"title"`
}

// DefaultSiteTitle is the fallback when neither pages.yaml nor index.md provides a title
const DefaultSiteTitle = "Site"

// Default returns the default configuration (site title resolved later from index.md or DefaultSiteTitle)
func Default() *Config {
	return &Config{Title: ""}
}

// Load reads and parses the configuration file. If the file does not exist, returns Default().
// When the file exists, unspecified fields are filled with default values.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse pages.yaml: %w", err)
	}

	cfg.ApplyDefaults()
	return &cfg, nil
}

// ApplyDefaults fills empty fields with default values
func (c *Config) ApplyDefaults() {
	if c.Title == "" {
		c.Title = Default().Title
	}
}

