package conn

import (
	"context"
	"database/sql"
)

type DB struct {
	// This is the database connection.
	conn *sql.DB
}

type DBInterface interface {
	Close() error
	Ping() error
	QueryWithContext(ctx context.Context, query string, args ...interface{}) error
	ExecWithContext(ctx context.Context, query string, args ...interface{}) error
	QueryRowWithContext(ctx context.Context, query string, args ...interface{}) error
	Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Close() error {
	// Close the database connection.
	return db.conn.Close()
}

func (db *DB) Ping() error {
	// Ping the database connection.
	return db.conn.Ping()
}

func (db *DB) QueryWithContext(ctx context.Context, query string, args ...interface{}) error {
	// QueryWithContext queries the database.
	return db.conn.QueryRowContext(ctx, query, args...).Scan()
}
