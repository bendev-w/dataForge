package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/dataforge/db"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all data from database tables",
	Long:  "Deletes all records from all tables in the database (respecting foreign key constraints)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClear()
	},
}

func init() {
	clearCmd.Flags().StringVar(&dbURL, "db", "", "Database URL (postgres://... or mysql://user:password@host:3306/database)")
	clearCmd.MarkFlagRequired("db")
}

func runClear() error {
	// Connect to database
	conn, err := db.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()

	log.Printf("Connected to database: %s", dbURL)
	log.Println("Clearing all data from database tables...")

	// Start transaction
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

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("✅ Successfully cleared all database records")
	return nil
}
