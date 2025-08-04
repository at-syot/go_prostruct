package repo

import (
	"context"
	"database/sql"

	"github.com/simt/dtacc/model"
	"github.com/uptrace/bun"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// BunUserRepository is a Bun-based implementation of UserRepository.
type BunUserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &BunUserRepository{db: db}
}

func (r *BunUserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *BunUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *BunUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *BunUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	user, err := r.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
