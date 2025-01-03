package postgres

import "time"

// Option -.
type Option func(*DB)

// MaxOpenConns -.
func MaxOpenConns(size int) Option {
	return func(c *DB) {
		c.Pool.SetMaxOpenConns(size)
	}
}

// MaxIdleConns -.
func MaxIdleConns(size int) Option {
	return func(c *DB) {
		c.Pool.SetMaxIdleConns(size)
	}
}

// ConnMaxLifetime -.
func ConnMaxLifetime(duration string) Option {
	d, err := time.ParseDuration(duration)
	if err != nil {
		d = defaultConnMaxLifetime
	}

	return func(c *DB) {
		c.Pool.SetConnMaxLifetime(d)
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *DB) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *DB) {
		c.connTimeout = timeout
	}
}
