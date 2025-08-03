package dtacc

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/simt/dtacc/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDB() (*bun.DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("could not load .env file: %w", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run temporary migration
	if err := tempMigrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}

func tempMigrate(db *bun.DB) error {
	ctx := context.Background()
	if _, err := db.NewCreateTable().Model((*model.User)(nil)).IfNotExists().Exec(ctx); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}
