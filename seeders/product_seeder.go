package seeders

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5"
)

var productCategories = []string{
	"Electronics",
	"Clothing",
	"Home & Garden",
	"Sports",
	"Books",
	"Toys",
	"Beauty",
	"Food & Beverage",
}

func seedProducts(ctx context.Context, tx pgx.Tx) error {
	const numProducts = 20  // Reduced from 100

	for i := 0; i < numProducts; i++ {
		name := gofakeit.ProductName()
		description := gofakeit.Sentence(10)
		price := gofakeit.Price(10, 500)
		stock := gofakeit.Number(10, 1000)
		category := productCategories[gofakeit.Number(0, len(productCategories)-1)]
		// Make SKU unique by adding timestamp
		sku := fmt.Sprintf("SKU-%s-%d-%d", category[:3], i, time.Now().UnixNano())

		_, err := tx.Exec(ctx,
			`INSERT INTO products (name, description, price, stock_quantity, category, sku)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			name, description, price, stock, category, sku)

		if err != nil {
			return err
		}
	}

	log.Printf("[Seeder] Inserted %d products", numProducts)
	return nil
}
