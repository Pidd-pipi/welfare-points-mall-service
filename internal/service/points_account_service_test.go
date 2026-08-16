package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

type mockPointsAccountRepo struct {
	accounts     map[uint]*model.PointsAccount
	transactions []*model.PointsTransaction
	seq          uint
}

func newMockPointsAccountRepo() *mockPointsAccountRepo {
	return &mockPointsAccountRepo{accounts: make(map[uint]*model.PointsAccount)}
}

func (m *mockPointsAccountRepo) Create(account *model.PointsAccount) error {
	return m.CreateTx(nil, account)
}
func (m *mockPointsAccountRepo) CreateTx(tx *gorm.DB, account *model.PointsAccount) error {
	m.seq++
	account.ID = m.seq
	m.accounts[account.UserID] = account
	return nil
}
func (m *mockPointsAccountRepo) FindByUserID(userID uint) (*model.PointsAccount, error) {
	return m.FindByUserIDTx(nil, userID)
}
func (m *mockPointsAccountRepo) FindByUserIDTx(tx *gorm.DB, userID uint) (*model.PointsAccount, error) {
	if a, ok := m.accounts[userID]; ok {
		cp := *a
		return &cp, nil
	}
	return nil, util.ErrNotFound
}
func (m *mockPointsAccountRepo) FindByID(id uint) (*model.PointsAccount, error) {
	return nil, util.ErrNotFound
}
func (m *mockPointsAccountRepo) Update(account *model.PointsAccount) error {
	return m.UpdateTx(nil, account)
}
func (m *mockPointsAccountRepo) UpdateTx(tx *gorm.DB, account *model.PointsAccount) error {
	m.accounts[account.UserID] = account
	return nil
}
func (m *mockPointsAccountRepo) CreateTransaction(transaction *model.PointsTransaction) error {
	return m.CreateTransactionTx(nil, transaction)
}
func (m *mockPointsAccountRepo) CreateTransactionTx(tx *gorm.DB, transaction *model.PointsTransaction) error {
	m.transactions = append(m.transactions, transaction)
	return nil
}
func (m *mockPointsAccountRepo) ListTransactions(userID uint, page, pageSize int) ([]model.PointsTransaction, int64, error) {
	return nil, 0, nil
}
func (m *mockPointsAccountRepo) MonthlyEarnSpend(userID uint, month string) (int, int, error) {
	return 0, 0, nil
}

func newTestPointsAccountService() (PointsAccountService, *mockPointsAccountRepo) {
	repo := newMockPointsAccountRepo()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewPointsAccountService(repo, nil, logger), repo
}

func TestPointsAccountGrantAndDeduct(t *testing.T) {
	svc, _ := newTestPointsAccountService()
	if _, err := svc.GrantTx(nil, 1, 500, "初始积分", nil); err != nil {
		t.Fatalf("grant failed: %v", err)
	}
	if _, err := svc.DeductTx(nil, 1, 200, "兑换", nil); err != nil {
		t.Fatalf("deduct failed: %v", err)
	}
	account, err := svc.GetByUserID(1)
	if err != nil {
		t.Fatalf("get account failed: %v", err)
	}
	if account.Balance != 300 {
		t.Fatalf("balance=%d, want 300", account.Balance)
	}
	if account.TotalEarned != 500 || account.TotalSpent != 200 {
		t.Fatalf("earned=%d spent=%d, want 500/200", account.TotalEarned, account.TotalSpent)
	}
}

func TestPointsAccountDeductInsufficient(t *testing.T) {
	svc, _ := newTestPointsAccountService()
	if _, err := svc.GrantTx(nil, 2, 100, "初始积分", nil); err != nil {
		t.Fatalf("grant failed: %v", err)
	}
	_, err := svc.DeductTx(nil, 2, 101, "超额兑换", nil)
	if err == nil {
		t.Fatal("expected insufficient balance error")
	}
	if !errors.Is(err, util.ErrPointsNotEnough) {
		t.Fatalf("expected ErrPointsNotEnough, got %v", err)
	}
}

func TestPointsAccountRefund(t *testing.T) {
	svc, _ := newTestPointsAccountService()
	if _, err := svc.GrantTx(nil, 3, 100, "初始积分", nil); err != nil {
		t.Fatalf("grant failed: %v", err)
	}
	if _, err := svc.DeductTx(nil, 3, 80, "兑换", nil); err != nil {
		t.Fatalf("deduct failed: %v", err)
	}
	if _, err := svc.RefundTx(nil, 3, 80, "取消订单", nil); err != nil {
		t.Fatalf("refund failed: %v", err)
	}
	account, err := svc.GetByUserID(3)
	if err != nil {
		t.Fatalf("get account failed: %v", err)
	}
	if account.Balance != 100 {
		t.Fatalf("balance=%d, want 100", account.Balance)
	}
}
