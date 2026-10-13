package seeders

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dataforge/db"
)

var orderStatuses = []string{
	"pending",
	"processing",
	"shipped",
	"delivered",
	"cancelled",
}

func seedOrders(ctx context.Context, tx *db.Tx) error {
	// Get all customer IDs from the database
	log.Println("[Orders] Querying customers...")
	rows, err := tx.Query(ctx, "SELECT id FROM customers ORDER BY id")
	if err != nil {
		return err
	}

	var customerIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		customerIDs = append(customerIDs, id)
	}
	rows.Close() // Explicitly close before doing Exec

	if err := rows.Err(); err != nil {
		return err
	}

	log.Printf("[Orders] Got %d customers", len(customerIDs))

	if len(customerIDs) == 0 {
		return nil // No customers to create orders for
	}

	const ordersPerCustomer = 2 // Reduced from 5

	insertCount := 0
	for _, customerID := range customerIDs {
		for j := 0; j < ordersPerCustomer; j++ {
			totalAmount := gofakeit.Price(50, 2000)
			status := orderStatuses[gofakeit.Number(0, len(orderStatuses)-1)]

			// Random order date within last 90 days
			daysAgo := gofakeit.Number(0, 90)
			orderDate := time.Now().AddDate(0, 0, -daysAgo)

			address := gofakeit.Address().Address

			_, err := tx.Exec(ctx,
				`INSERT INTO orders (customer_id, order_date, total_amount, status, shipping_address)
				VALUES ($1, $2, $3, $4, $5)`,
				customerID, orderDate, totalAmount, status, address)

			if err != nil {
				return err
			}
			insertCount++
		}
	}

	log.Printf("[Orders] Inserted %d orders", insertCount)
	return nil
}
