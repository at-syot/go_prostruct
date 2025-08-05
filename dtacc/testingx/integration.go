package testingx

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const migrationsDir = "dtacc/migrations"

// SetupTestDB starts a Postgres Testcontainer, applies migrations, and returns a bun.DB instance.
// It is reusable for integration tests in other modules.
func SetupTestDB() (*bun.DB, error) {
	wd, _ := os.Getwd()
	println("current testing dir", wd)

	ctx := context.Background()
	pgContainer, err := tcpostgres.Run(ctx, "postgres:16")
	if err != nil {
		return nil, err
	}
	defer func() { _ = pgContainer.Terminate(ctx) }()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	println("got dsn", dsn)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := WaitDB(ctx, db); err != nil {
		return nil, err
	}

	// Run Goose migrations
	if err := goose.Up(sqldb, migrationsDir); err != nil {
		return nil, err
	}

	return db, nil
}

func WaitDB(ctx context.Context, db *bun.DB) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := db.Ping(); err == nil {
				return nil
			}
			time.Sleep(time.Second)
		}
	}

}
