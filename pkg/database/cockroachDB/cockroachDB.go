package cockroachdb

import (
	"context"

	"embed"
	"fmt"
	"log"
	"time"

	config "proposal-template/pkg/utils/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// region: ======= CockroachDB Configuration =======

var DefaultConfig = CockroachDBConfig{
	URI:                   "postgresql://root@localhost:26257/defaultdb?sslmode=disable",
	MaxOpenConns:          25,
	MaxIdleConns:          25,
	ConnMaxLifetimeInSecs: 300,
}

var _ config.IConfig = (*CockroachDBConfig)(nil)

func (c *CockroachDBConfig) Load() error {
	log.Printf("Loading CockroachDBConfig")
	return config.ParseConfig(c)
}

func NewCockroachDB(opts ...Option) (*gorm.DB, error) {
	cfg := DefaultConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure GORM database connection
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(cfg.URI), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CockroachDB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeInSecs) * time.Second)

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// region: ======= CockroachDB Migration =======

func CockroachDBGooseMigrate(db *gorm.DB, baseFS embed.FS, migrationFolder string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from GORM: %v", err)
	}
	
	log.Printf("Running CockroachDB migrations using goose\n")

	// Set the base filesystem for migrations
	goose.SetBaseFS(baseFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("CockroachDB migrations: could not set dialect: %v", err)
	}

	// Run migrations
	if err := goose.Up(sqlDB, migrationFolder); err != nil {
		return fmt.Errorf("cockroachDB migrations: could not apply migrations: %v", err)
	}
	
	log.Printf("CockroachDB migrations applied successfully")
	return nil
}