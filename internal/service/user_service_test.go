package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

type mockUserRepo struct {
	users map[string]*model.User
	seq   uint
}

func newMockUserRepo() *mockUserRepo { return &mockUserRepo{users: map[string]*model.User{}} }

func (m *mockUserRepo) Create(u *model.User) error {
	if _, ok := m.users[u.Username]; ok {
		return repository.ErrDuplicate
	}
	m.seq++
	u.ID = m.seq
	m.users[u.Username] = u
	return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, util.ErrNotFound
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, util.ErrNotFound
}

func (m *mockUserRepo) ListEmployees() ([]model.User, error) {
	var list []model.User
	for _, u := range m.users {
		if u.Role == constants.RoleEmployee {
			list = append(list, *u)
		}
	}
	return list, nil
}

func (m *mockUserRepo) List(page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) Update(u *model.User) error {
	m.users[u.Username] = u
	return nil
}

func newTestUserService() UserService {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewUserService(newMockUserRepo(), logger, "test-secret", 72)
}

func TestUserServiceRegister(t *testing.T) {
	svc := newTestUserService()
	tests := []struct {
		name     string
		username string
		password string
		role     constants.UserRole
		wantErr  bool
	}{
		{"valid employee", "alice", "123456", constants.RoleEmployee, false},
		{"valid hr", "bob", "123456", constants.RoleHR, false},
		{"empty username", "", "123456", constants.RoleAdmin, true},
		{"short password", "carol", "123", constants.RoleAdmin, true},
		{"duplicate", "alice", "123456", constants.RoleAdmin, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := svc.Register(tt.username, tt.password, "测试", "", "", "", tt.role)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", u)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if u.PasswordHash == "" || u.PasswordHash == tt.password {
				t.Error("password not hashed")
			}
		})
	}
}

func TestUserServiceLogin(t *testing.T) {
	svc := newTestUserService()
	if _, err := svc.Register("lisa", "123456", "丽萨", "", "", "", constants.RoleEmployee); err != nil {
		t.Fatalf("seed register: %v", err)
	}
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"correct", "123456", false},
		{"wrong", "bad-pass", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, token, err := svc.Login("lisa", tt.password)
			if tt.wantErr {
				if !errors.Is(err, util.ErrUnauthorized) {
					t.Fatalf("expected ErrUnauthorized, got %v", err)
				}
				return
			}
			if err != nil || token == "" {
				t.Fatalf("login failed: %v", err)
			}
		})
	}
}
