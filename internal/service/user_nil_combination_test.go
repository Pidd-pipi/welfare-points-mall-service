package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

type nilUserRepo struct{}

func (m *nilUserRepo) Create(u *model.User) error { return nil }
func (m *nilUserRepo) FindByUsername(username string) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) FindByID(id uint) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) ListEmployees() ([]model.User, error) { return nil, nil }
func (m *nilUserRepo) List(page, pageSize int) ([]model.User, int64, error) { return nil, 0, nil }
func (m *nilUserRepo) Update(u *model.User) error { return nil }

func newNilUserService() UserService {
	return NewUserService(&nilUserRepo{}, newSlogLogger(), "test-secret", 72)
}

func TestGetByIDNilUserReturnsError(t *testing.T) {
	svc := newNilUserService()
	user, err := svc.GetByID(1)
	if err == nil {
		t.Fatalf("expected error for nil user, got user=%+v", user)
	}
	if user != nil {
		t.Fatalf("user should be nil, got %+v", user)
	}
}

func TestLoginNilUserUnauthorized(t *testing.T) {
	svc := newNilUserService()
	_, _, err := svc.Login("alice", "123456")
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	if !errors.Is(err, util.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestRegisterStillWorks(t *testing.T) {
	svc := newTestUserService()
	u, err := svc.Register("bob", "123456", "测试", "", "", "", constants.RoleEmployee)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if u == nil {
		t.Fatal("register returned nil user")
	}
}

func newSlogLogger() *slog.Logger { return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})) }
