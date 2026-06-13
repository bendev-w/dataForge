package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
