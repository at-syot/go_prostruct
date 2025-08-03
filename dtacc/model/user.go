package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

// User represents a user in the auth system
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Email        string    `bun:"email,unique,notnull" json:"email"`
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:now()" json:"created_at"`
}

// NewUser creates a new user with a generated UUID
func NewUser(email, passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
}
