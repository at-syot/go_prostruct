package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/simt/auth/repo"
	"github.com/simt/dtacc/model"
	"github.com/simt/pkg/testingx"
)

func TestMain(m *testing.M) {
	if err := testingx.ApplyProjectRootDir(); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestRegisterService(t *testing.T) {
	// stub
	oat := &model.User{Email: "oat@example.com"}
	art := &model.User{Email: "art@example.com"}
	yok := &model.User{Email: "yok@example.com"}

	usersByEmail := make(map[string]*model.User)
	usersByEmail["oat@example.com"] = oat
	usersByEmail["art@example.com"] = art
	usersByEmail["yok@example.com"] = yok

	// sut
	ctx := context.Background()
	sutRepo := newMockUserRepo()
	sutRepo.usersByEmail = usersByEmail

	type expected struct {
		hasResp bool
		errFn   func(error) bool
	}
	type testcase struct {
		name      string
		setupRepo func() repo.UserRepository
		req       RegisterRequest
		expected  expected
	}

	cases := []testcase{
		{
			name: "should validate req fails",
			setupRepo: func() repo.UserRepository {
				return sutRepo
			},
			req: RegisterRequest{},
			expected: expected{
				hasResp: false,
				errFn: func(err error) bool {
					var valErrs validator.ValidationErrors
					return errors.As(err, &valErrs)
				},
			},
		},
		{
			name: "should got email is exists, when register with duplicate email",
			setupRepo: func() repo.UserRepository {
				sutRepo.usersByEmail["aot@example.com"] = &model.User{Email: "aot@example.com"}
				return sutRepo
			},
			req: RegisterRequest{Email: "aot@example.com", Password: "12345678"},
			expected: expected{
				hasResp: false,
				errFn: func(err error) bool {
					return errors.Is(err, ErrEmailAlreadyExists)
				},
			},
		},
		{
			name: "should register successfully",
			setupRepo: func() repo.UserRepository {
				return sutRepo
			},
			req: RegisterRequest{Email: "new@example.com", Password: "12345678"},
			expected: expected{
				hasResp: true,
				errFn: func(err error) bool {
					return err == nil
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sutService := NewRegisterService(c.setupRepo())
			gotResp, gotErr := sutService.Register(ctx, c.req)
			if !c.expected.errFn(gotErr) {
				t.Error("error is not match")
			}

			if c.expected.hasResp {
				if gotResp == nil {
					t.Fatal("expected resp")
				}
			}
		})
	}
}

type mockUserRepo struct {
	usersByID    map[string]*model.User
	usersByEmail map[string]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByID:    make(map[string]*model.User),
		usersByEmail: make(map[string]*model.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if user == nil || user.ID.String() == "" || user.Email == "" {
		return nil
	}
	m.usersByID[user.ID.String()] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, ok := m.usersByID[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	_, ok := m.usersByEmail[email]
	return ok, nil
}
