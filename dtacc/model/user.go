package model

import (
	"database/sql"

	"github.com/uptrace/bun"
)

type DTModel struct {
	Cat sql.NullTime   `bun:"created_at"`
	Cby sql.NullString `bun:"created_by"`
	Uat sql.NullTime   `bun:"updated_at"`
	Uby sql.NullString `bun:"updated_by"`
	Dat sql.NullTime   `bun:"deleted_at"`
	Dby sql.NullString `bun:"deleted_by"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string `bun:",pk,autoincrement"`
	Uname         string `bun:",unique,notnull"`
	Hash          string `bun:",notnull"`
	DTModel
}
