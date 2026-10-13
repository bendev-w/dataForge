package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
	driver string
}

type Tx struct {
	*sql.Tx
	driver string
}

func Connect(ctx context.Context, connString string) (*DB, error) {
	driver, dsn, err := connectionDetails(connString)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	return &DB{DB: conn, driver: driver}, nil
}

func (db *DB) Begin(ctx context.Context) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, driver: db.driver}, nil
}

func (db *DB) Driver() string {
	return db.driver
}

func (tx *Tx) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.ExecContext(ctx, rebind(query, tx.driver), args...)
}

func (tx *Tx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.QueryContext(ctx, rebind(query, tx.driver), args...)
}

func connectionDetails(connString string) (string, string, error) {
	if strings.HasPrefix(connString, "mysql://") || strings.HasPrefix(connString, "mysql+tcp://") {
		dsn, err := mysqlDSN(connString)
		return "mysql", dsn, err
	}
	return "pgx", connString, nil
}

func mysqlDSN(connString string) (string, error) {
	parsed, err := url.Parse(strings.TrimPrefix(connString, "mysql+tcp://"))
	if err != nil {
		return "", fmt.Errorf("invalid MySQL connection URL: %w", err)
	}
	if parsed.Host == "" || parsed.Path == "" {
		return "", fmt.Errorf("MySQL URL must include a host and database")
	}

	user := ""
	password := ""
	if parsed.User != nil {
		user = parsed.User.Username()
		password, _ = parsed.User.Password()
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true", user, password, parsed.Host, strings.TrimPrefix(parsed.Path, "/"))
	if parsed.RawQuery != "" {
		dsn += "&" + parsed.RawQuery
	}
	return dsn, nil
}

func rebind(query, driver string) string {
	if driver != "mysql" {
		return query
	}
	for i := 20; i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("$%d", i), "?")
	}
	return query
}
