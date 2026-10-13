package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dataforge/db"
	"github.com/dataforge/seeders"
	"github.com/spf13/cobra"
)

var (
	dbURL    string
	interval string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start seeding the database at specified intervals",
	Long:  "Connects to the database and seeds it with realistic data at the specified interval",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStart()
	},
}

func init() {
	startCmd.Flags().StringVar(&dbURL, "db", "", "Database URL (postgres://... or mysql://user:password@host:3306/database)")
	startCmd.Flags().StringVar(&interval, "interval", "5s", "Interval between seeding operations (e.g., 5s, 1m, 10m)")

	startCmd.MarkFlagRequired("db")
}

func runStart() error {
	// Parse interval
	dur, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	// Connect to database
	conn, err := db.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()

	// Create schema
	if err := db.InitializeSchema(context.Background(), conn); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Printf("Connected to database: %s", dbURL)
	log.Printf("Seeding interval: %v", dur)
	log.Println("Starting data seeding in APPEND MODE (data accumulates over time)...")

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	// Seed immediately on start
	seedDatabase(conn)

	// Seed at intervals (synchronously, one at a time)
	for range ticker.C {
		seedDatabase(conn)
	}

	return nil
}

func seedDatabase(conn *db.DB) {
	startTime := time.Now()
	log.Println("Starting seeding cycle...")

	if err := seeders.SeedAll(context.Background(), conn); err != nil {
		log.Printf("❌ Error seeding database: %v", err)
	} else {
		duration := time.Since(startTime)
		log.Printf("✅ Successfully seeded database in %v", duration)
	}
}
