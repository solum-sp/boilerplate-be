package cockroachdb

import "proposal-template/pkg/logger"

// CockroachDBConfig holds database configuration options
type CockroachDBConfig struct {
	URI                   string
	MaxOpenConns          int
	MaxIdleConns          int
	ConnMaxLifetimeInSecs int
	Logger                logger.ILogger
}
// Option represents a functional option for CockroachDB configuration
type Option func(*CockroachDBConfig)

// WithURI sets the CockroachDB connection URI
func WithURI(uri string) Option {
	return func(c *CockroachDBConfig) {
		c.URI = uri
	}
}

// WithMaxOpenConns sets the max open connections
func WithMaxOpenConns(max int) Option {
	return func(c *CockroachDBConfig) {
		c.MaxOpenConns = max
	}
}

// WithMaxIdleConns sets the max idle connections
func WithMaxIdleConns(max int) Option {
	return func(c *CockroachDBConfig) {
		c.MaxIdleConns = max
	}
}

// WithConnMaxLifetime sets the max connection lifetime
func WithConnMaxLifetime(seconds int) Option {
	return func(c *CockroachDBConfig) {
		c.ConnMaxLifetimeInSecs = seconds
	}
}

func WithLogger(customLogger logger.ILogger) Option {
	return func(c *CockroachDBConfig) {
		c.Logger = customLogger
	}
}