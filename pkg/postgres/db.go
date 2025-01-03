// Package postgres implements postgres connection.
package postgres

import (
	"database/sql"
	"log"
	"time"
)

const (
	defaultMaxOpenConns    = 20
	defaultMaxIdleConns    = 20
	defaultConnMaxLifetime = time.Minute * 5
	defaultConnAttempts    = 10
	defaultConnTimeout     = time.Second
)

// DB -.
type DB struct {
	maxOpenConns    int
	maxIdleConns    int
	maxConnLifetime time.Duration
	connAttempts    int
	connTimeout     time.Duration

	Pool *sql.DB
}

// New -.
func New(url string, opts ...Option) (*DB, error) {
	cfg := &DB{
		maxOpenConns:    defaultMaxOpenConns,
		maxIdleConns:    defaultMaxIdleConns,
		maxConnLifetime: defaultConnMaxLifetime,
		connAttempts:    defaultConnAttempts,
		connTimeout:     defaultConnTimeout,
	}

	for cfg.connAttempts > 0 {
		db, err := sql.Open("postgres", url)
		if err == nil {
			cfg.Pool = db
			break
		}

		log.Printf("[Postgres] trying to connect, attempts left: %d", cfg.connAttempts)

		time.Sleep(cfg.connTimeout)

		cfg.connAttempts--
	}

	// Custom options
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, nil
}

// Close -.
func (p *DB) Close() {
	if p.Pool != nil {
		_ = p.Pool.Close()
	}
}
