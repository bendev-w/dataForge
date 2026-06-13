package seeders

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5"
)

func seedCustomers(ctx context.Context, tx pgx.Tx) error {
	const numCustomers = 10  // Reduced from 50

	for i := 0; i < numCustomers; i++ {
		person := gofakeit.Person()
		// Make email unique by adding timestamp
		email := gofakeit.Email() + "-" + fmt.Sprintf("%d", time.Now().UnixNano())
		phone := gofakeit.Phone()
		address := gofakeit.Address().Address
		city := gofakeit.City()
		state := gofakeit.State()
		zipCode := gofakeit.Zip()
		country := gofakeit.Country()

		_, err := tx.Exec(ctx,
			`INSERT INTO customers (first_name, last_name, email, phone, address, city, state, zip_code, country)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			person.FirstName, person.LastName, email, phone, address, city, state, zipCode, country)

		if err != nil {
			return err
		}
	}

	log.Printf("[Seeder] Inserted %d customers", numCustomers)
	return nil
}
