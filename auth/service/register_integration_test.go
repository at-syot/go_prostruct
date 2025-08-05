//##go:build integration

package service

import (
	"os"
	"testing"

	"context"

	"github.com/simt/auth/repo"
	"github.com/simt/dtacc/model"
	dtacctesting "github.com/simt/dtacc/testingx"
	pkgtesting "github.com/simt/pkg/testingx"
	"github.com/uptrace/bun"
)

var testDB *bun.DB

func TestMain(m *testing.M) {
	if err := pkgtesting.ApplyProjectRootDir(); err != nil {
		os.Exit(1)
	}

	db, err := dtacctesting.SetupTestDB()
	if err != nil {
		os.Exit(1)
	}
	testDB = db

	code := m.Run()
	os.Exit(code)
}

func TestRegisterService_Integration(t *testing.T) {
	// Clean up users table before each test run for isolation
	_, _ = testDB.NewTruncateTable().Table("users").Cascade().Exec(context.Background())

	sutRepo := repo.NewUserRepository(testDB)
	sutService := NewRegisterService(sutRepo)
	ctx := context.Background()

	type expected struct {
		hasResp bool
		errFn   func(error) bool
	}
	type testcase struct {
		name     string
		setup    func()
		req      RegisterRequest
		expected expected
	}

	cases := []testcase{
		{
			name:  "should validate req fails",
			setup: func() {},
			req:   RegisterRequest{},
			expected: expected{
				hasResp: false,
				errFn: func(err error) bool {
					// validator.ValidationErrors
					return err != nil
				},
			},
		},
		{
			name: "should got email is exists, when register with duplicate email",
			setup: func() {
				// Insert user with duplicate email
				user := &model.User{
					Email:        "aot@example.com",
					PasswordHash: "dummyhash",
				}
				_ = sutRepo.Create(ctx, user)
			},
			req: RegisterRequest{Email: "aot@example.com", Password: "12345678"},
			expected: expected{
				hasResp: false,
				errFn: func(err error) bool {
					return err != nil && err.Error() == ErrEmailAlreadyExists.Error()
				},
			},
		},
		{
			name:  "should register successfully",
			setup: func() {},
			req:   RegisterRequest{Email: "new@example.com", Password: "12345678"},
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
			// Clean up before each subtest for isolation
			_, _ = testDB.NewTruncateTable().Table("users").Cascade().Exec(ctx)
			c.setup()
			gotResp, gotErr := sutService.Register(ctx, c.req)
			if !c.expected.errFn(gotErr) {
				t.Errorf("error is not match, got: %v", gotErr)
			}

			if c.expected.hasResp {
				if gotResp == nil {
					t.Fatal("expected resp")
				}
			}
		})
	}
}
