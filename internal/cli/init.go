package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rubellum/pages/internal/template"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new pages project",
	Long:  `Initialize a new pages project with default configuration, template, and styles.`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Create directory structure
	dirs := []string{
		"src",
		"public/css",
		"templates",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create base.html template
	if err := writeFileIfNotExists(filepath.Join("templates", "base.html"), []byte(template.DefaultBaseHTML)); err != nil {
		return err
	}

	// Create style.css
	if err := writeFileIfNotExists(filepath.Join("public", "css", "style.css"), []byte(template.DefaultStyleCSS)); err != nil {
		return err
	}

	// Create TOP page (its title becomes the site title)
	if err := writeFileIfNotExists(filepath.Join("src", "index.md"), []byte(template.DefaultIndexMD)); err != nil {
		return err
	}

	// Create .gitkeep in src (for when user deletes index.md and wants to keep dir)
	if err := writeFileIfNotExists(filepath.Join("src", ".gitkeep"), []byte{}); err != nil {
		return err
	}

	fmt.Println("Initialized pages project")
	fmt.Println("")
	fmt.Println("Created:")
	fmt.Println("  templates/base.html - HTML template")
	fmt.Println("  public/css/style.css - Default styles")
	fmt.Println("  src/index.md       - TOP page (its title is the site title)")
	fmt.Println("  src/               - Markdown source directory")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("  1. Run 'pages build' to generate HTML")
	fmt.Println("  (Optional) Add pages.yaml to override site title")

	return nil
}

func writeFileIfNotExists(path string, content []byte) error {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Skipping %s (already exists)\n", path)
		return nil
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}

	return nil
}
