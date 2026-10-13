package seeders

import (
	"context"
	"log"

	"github.com/dataforge/db"
)

func SeedAll(ctx context.Context, conn *db.DB) error {
	// Start transaction
	log.Println("[Seeder] Starting transaction...")
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Seed in order due to foreign key dependencies (append mode - no deletes)
	log.Println("[Seeder] Seeding customers...")
	if err := seedCustomers(ctx, tx); err != nil {
		return err
	}

	log.Println("[Seeder] Seeding products...")
	if err := seedProducts(ctx, tx); err != nil {
		return err
	}

	log.Println("[Seeder] Seeding orders...")
	if err := seedOrders(ctx, tx); err != nil {
		return err
	}

	log.Println("[Seeder] Seeding order items...")
	if err := seedOrderItems(ctx, tx); err != nil {
		return err
	}

	log.Println("[Seeder] Seeding payments...")
	if err := seedPayments(ctx, tx); err != nil {
		log.Printf("[Seeder] Error in seedPayments: %v", err)
		return err
	}

	// Commit transaction
	log.Println("[Seeder] About to commit transaction...")
	commitErr := tx.Commit()
	log.Printf("[Seeder] Commit result: err=%v", commitErr)
	return commitErr
}
