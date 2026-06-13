package seeders

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var paymentMethods = []string{
	"credit_card",
	"debit_card",
	"paypal",
	"bank_transfer",
	"apple_pay",
	"google_pay",
}

var paymentStatuses = []string{
	"pending",
	"completed",
	"failed",
	"refunded",
}

func seedPayments(ctx context.Context, tx pgx.Tx) error {
	// Get all orders
	log.Println("[Payments] Querying orders...")
	rows, err := tx.Query(ctx, "SELECT id, total_amount, order_date FROM orders")
	if err != nil {
		return err
	}

	type orderData struct {
		id        int
		amount    float64
		orderDate time.Time
	}

	var orders []orderData
	for rows.Next() {
		var o orderData
		if err := rows.Scan(&o.id, &o.amount, &o.orderDate); err != nil {
			rows.Close()
			return err
		}
		orders = append(orders, o)
	}
	rows.Close()  // Explicitly close before doing Exec
	
	if err := rows.Err(); err != nil {
		return err
	}

	log.Printf("[Payments] Got %d orders to create payments for", len(orders))

	insertCount := 0
	for _, order := range orders {
		method := paymentMethods[gofakeit.Number(0, len(paymentMethods)-1)]
		status := paymentStatuses[gofakeit.Number(0, len(paymentStatuses)-1)]
		transactionID := fmt.Sprintf("TXN-%d-%d", order.id, gofakeit.Number(100000, 999999))

		// Payment date usually shortly after order date
		paymentDate := order.orderDate.Add(time.Duration(gofakeit.Number(1, 48)) * time.Hour)

		var dt pgtype.Timestamp
		dt.Scan(paymentDate)

		_, err := tx.Exec(ctx,
			`INSERT INTO payments (order_id, amount, payment_method, status, transaction_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			order.id, order.amount, method, status, transactionID, dt)

		if err != nil {
			log.Printf("[Payments] Error inserting payment for order %d: %v", order.id, err)
			return err
		}
		insertCount++
	}

	log.Printf("[Payments] Inserted %d payments", insertCount)
	return nil
}
