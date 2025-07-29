package dtacc

import (
	"context"
	"database/sql"

	"github.com/simt/dtacc/model"
	_ "github.com/simt/dtacc/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDB() (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := tempMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func tempMigrate(db *bun.DB) error {
	ctx := context.Background()
	if _, err := db.NewCreateTable().Model((*model.User)(nil)).Exec(ctx); err != nil {
		return err
	}

	return nil
}
