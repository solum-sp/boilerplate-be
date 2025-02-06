package dao

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	utils "proposal-template/pkg/utils"

	"github.com/pressly/goose/v3"
)

// region: ======= CockroachDB Configuration =======
type CockroachDBConfig struct {
	URI                   string `env:"COCKROACH_URI,required"`
	MaxOpenConns          int    `env:"COCKROACH_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns          int    `env:"COCKROACH_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetimeInSecs int    `env:"COCKROACH_CONN_MAX_LIFETIME_IN_SECS" envDefault:"300"`
}

var _ utils.Config = (*CockroachDBConfig)(nil)

func (c *CockroachDBConfig) Load() error {
	log.Printf("Loading CockroachDBConfig")
	return utils.ParseConfig(c)
}


func NewCockroachDB(cfg *CockroachDBConfig) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", cfg.URI)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeInSecs) * time.Second)

	return db, nil
}

// region: ======= CockroachDB Migration =======

func CockroachDBMigrate(db *sql.DB, baseFS embed.FS, migrationFolder string) error {
	log.Printf("Running CockroachDB migrations\n")

	// Set the base filesystem for migrations
	goose.SetBaseFS(baseFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("CockroachDB migrations: could not set dialect: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, migrationFolder); err != nil {
		return fmt.Errorf("cockroachDB migrations: could not apply migrations: %v", err)
	}

	log.Printf("CockroachDB migrations: applied successfully")
	return nil
}