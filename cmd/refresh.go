package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/dataforge/db"
	"github.com/dataforge/seeders"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Clear database and seed with new data",
	Long:  "Deletes all records from the database and seeds it with fresh realistic data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRefresh()
	},
}

func init() {
	refreshCmd.Flags().StringVar(&dbURL, "db", "", "Database connection URL (e.g., postgres://localhost:5432/shop)")
	refreshCmd.MarkFlagRequired("db")
}

func runRefresh() error {
	// Connect to database
	conn, err := db.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()

	log.Printf("Connected to database: %s", dbURL)
	log.Println("Refreshing database: clearing all data and seeding with fresh data...")

	// Start transaction for clearing
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete in reverse order of foreign key dependencies
	tables := []string{"payments", "order_items", "orders", "products", "customers"}

	for _, table := range tables {
		log.Printf("Clearing table: %s...", table)
		if _, err := tx.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s;", table)); err != nil {
			return fmt.Errorf("failed to clear %s: %w", table, err)
		}
	}

	// Reset auto-increment sequences
	if conn.Driver() == "mysql" {
		for _, table := range tables {
			log.Printf("Resetting auto-increment: %s...", table)
			if _, err := tx.Exec(context.Background(), fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1;", table)); err != nil {
				log.Printf("Warning: Could not reset auto-increment for %s: %v", table, err)
			}
		}
	} else {
		sequences := []string{"customers_id_seq", "products_id_seq", "orders_id_seq", "order_items_id_seq", "payments_id_seq"}
		for _, seq := range sequences {
			log.Printf("Resetting sequence: %s...", seq)
			if _, err := tx.Exec(context.Background(), fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1;", seq)); err != nil {
				log.Printf("Warning: Could not reset sequence %s (might not exist)", seq)
			}
		}
	}

	// Commit the clear transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit clear transaction: %w", err)
	}

	log.Println("✅ Database cleared successfully")

	// Now seed with new data
	log.Println("Seeding database with fresh data...")
	if err := seeders.SeedAll(context.Background(), conn); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	log.Println("✅ Successfully refreshed database with new data")
	return nil
}
