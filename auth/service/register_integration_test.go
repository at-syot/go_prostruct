//go:build integration

package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/simt/auth/repo"
	"github.com/simt/dtacc/model"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func setupTestDB(t *testing.T) *bun.DB {
	ctx := context.Background()
	pgContainer, err := tcpostgres.RunContainer(ctx,
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Wait for DB to be ready
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	// Create users table (or run migrations here)
	if _, err := db.NewCreateTable().Model((*model.User)(nil)).IfNotExists().Exec(ctx); err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	return db
}

func TestRegisterService_Integration(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repo.NewUserRepository(db)
	service := NewRegisterService(userRepo)

	ctx := context.Background()

	t.Run("register new user", func(t *testing.T) {
		req := RegisterRequest{Email: "integration@example.com", Password: "password123"}
		resp, err := service.Register(ctx, req)
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
		if resp.UserID == "" {
			t.Error("expected user ID to be set")
		}
	})

	t.Run("register duplicate email", func(t *testing.T) {
		req := RegisterRequest{Email: "integration@example.com", Password: "password123"}
		_, err := service.Register(ctx, req)
		if err == nil {
			t.Fatal("expected error for duplicate email, got nil")
		}
	})
}
