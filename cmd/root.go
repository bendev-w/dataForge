package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dataforge",
	Short: "DataForge - Database seeding tool",
	Long:  "DataForge is a powerful CLI tool for seeding PostgreSQL databases with realistic data at specified intervals",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
}
