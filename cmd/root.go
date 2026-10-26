package cmd

import (
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "dataforge",
	Version: version,
	Short: "DataForge - Database seeding tool",
	Long:  "DataForge is a powerful CLI tool for seeding PostgreSQL and MySQL databases with realistic data at specified intervals",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(refreshCmd)
}
