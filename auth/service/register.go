package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/simt/auth/repo"
	"github.com/simt/dtacc/model"
	"github.com/simt/pkg/cipherx"
)

type RegisterService struct {
	UserRepo  repo.UserRepository
	validator *validator.Validate
}

func NewRegisterService(userRepo repo.UserRepository) *RegisterService {
	return &RegisterService{
		UserRepo:  userRepo,
		validator: validator.New(),
	}
}

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type RegisterResponse struct {
	UserID    string
	CreatedAt time.Time
}

var ErrEmailAlreadyExists = errors.New("email already registered")

func (s *RegisterService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	exists, err := s.UserRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := cipherx.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
	}

	if err := s.UserRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		UserID:    user.ID.String(),
		CreatedAt: user.CreatedAt,
	}, nil
}
