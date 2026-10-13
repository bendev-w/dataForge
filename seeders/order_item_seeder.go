package seeders

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dataforge/db"
)

func seedOrderItems(ctx context.Context, tx *db.Tx) error {
	// Get all product IDs
	log.Println("[OrderItems] Querying products...")
	productRows, err := tx.Query(ctx, "SELECT id, price FROM products ORDER BY id")
	if err != nil {
		return err
	}

	type product struct {
		id    int
		price float64
	}

	var products []product
	for productRows.Next() {
		var p product
		if err := productRows.Scan(&p.id, &p.price); err != nil {
			productRows.Close()
			return err
		}
		products = append(products, p)
	}
	productRows.Close() // Explicitly close before next query

	if err := productRows.Err(); err != nil {
		return err
	}

	log.Printf("[OrderItems] Got %d products", len(products))

	if len(products) == 0 {
		return nil // No products to add to orders
	}

	// Get all orders
	log.Println("[OrderItems] Querying orders...")
	orderRows, err := tx.Query(ctx, "SELECT id FROM orders")
	if err != nil {
		return err
	}

	var orderIDs []int
	for orderRows.Next() {
		var id int
		if err := orderRows.Scan(&id); err != nil {
			orderRows.Close()
			return err
		}
		orderIDs = append(orderIDs, id)
	}
	orderRows.Close() // Explicitly close before doing Exec

	if err := orderRows.Err(); err != nil {
		return err
	}

	log.Printf("[OrderItems] Got %d orders", len(orderIDs))

	// For each order, add 1-2 items (reduced from 2-5)
	insertCount := 0
	for _, orderID := range orderIDs {
		numItems := gofakeit.Number(1, 2)

		for i := 0; i < numItems; i++ {
			// Pick random product
			randomProduct := products[gofakeit.Number(0, len(products)-1)]
			productID := randomProduct.id
			unitPrice := randomProduct.price

			quantity := gofakeit.Number(1, 5)
			subtotal := float64(quantity) * unitPrice

			_, err = tx.Exec(ctx,
				`INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
				VALUES ($1, $2, $3, $4, $5)`,
				orderID, productID, quantity, unitPrice, subtotal)

			if err != nil {
				return err
			}
			insertCount++
		}
	}

	log.Printf("[OrderItems] Inserted %d order items", insertCount)
	return nil
}
