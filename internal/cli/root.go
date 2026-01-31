package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pages",
	Short: "A minimal static site generator",
	Long:  `pages is a minimal static site generator that converts Markdown files to HTML.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

