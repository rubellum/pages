package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rubellum/pages/internal/builder"
	"github.com/rubellum/pages/internal/config"
	"github.com/rubellum/pages/internal/parser"
	"github.com/spf13/cobra"
)

var (
	inputDir  string
	outputDir string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	Long:  `Convert Markdown files to HTML and generate the static site.`,
	RunE:  runBuild,
}

func init() {
	buildCmd.Flags().StringVarP(&inputDir, "input", "i", "content", "Input directory containing Markdown files")
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "public", "Output directory for generated HTML")
	rootCmd.AddCommand(buildCmd)
}

// ResolveSiteTitleFromIndex reads inputDir/index.md and returns the frontmatter title, or "" if not found/unset.
// Exported for testing.
func ResolveSiteTitleFromIndex(inputDir string) string {
	indexPath := filepath.Join(inputDir, "index.md")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return ""
	}
	result, err := parser.Parse(data)
	if err != nil {
		return ""
	}
	return result.Title
}

func runBuild(cmd *cobra.Command, args []string) error {
	inputDir, _ := cmd.Flags().GetString("input")
	outputDir, _ := cmd.Flags().GetString("output")
	count, err := doBuild(inputDir, outputDir)
	if err != nil {
		return err
	}
	fmt.Printf("Built %d pages\n", count)
	return nil
}

// doBuild loads config, resolves site title, and runs the builder. Used by runBuild and tests.
func doBuild(inputDir, outputDir string) (int, error) {
	cfg, err := config.Load("pages.yaml")
	if err != nil {
		return 0, fmt.Errorf("error: %w", err)
	}

	// Resolve site title: pages.yaml title → index.md title → "Site"
	if cfg.Title == "" {
		cfg.Title = ResolveSiteTitleFromIndex(inputDir)
	}
	if cfg.Title == "" {
		cfg.Title = config.DefaultSiteTitle
	}

	opts := &builder.Options{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		TemplateDir: "templates",
	}

	b := builder.New(cfg, opts)
	return b.Build()
}
