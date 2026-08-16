package service

import (
	"testing"

	"github.com/ld/welfaremall/internal/model"
)

func TestPointsFlowGrantAndDeduct(t *testing.T) {
	svc, repo := newTestPointsAccountService()
	if _, err := svc.GrantTx(nil, 1, 500, "初始积分", nil); err != nil {
		t.Fatalf("grant failed: %v", err)
	}
	account, err := svc.GetByUserID(1)
	if err != nil {
		t.Fatalf("get account failed: %v", err)
	}
	if account.Balance != 500 {
		t.Fatalf("balance after grant = %d, want 500", account.Balance)
	}
	// manually seed a clean account for deduct test
	repo.accounts[2] = &model.PointsAccount{ID: 2, UserID: 2, Balance: 100, TotalEarned: 100}
	if _, err := svc.DeductTx(nil, 2, 30, "兑换", nil); err != nil {
		t.Fatalf("deduct failed: %v", err)
	}
	account, _ = svc.GetByUserID(2)
	if account.Balance != 70 {
		t.Fatalf("balance after deduct = %d, want 70", account.Balance)
	}
	if account.TotalSpent != 30 {
		t.Fatalf("total spent after deduct = %d, want 30", account.TotalSpent)
	}
}

func TestPointsFlowRefund(t *testing.T) {
	svc, repo := newTestPointsAccountService()
	repo.accounts[3] = &model.PointsAccount{ID: 3, UserID: 3, Balance: 70, TotalSpent: 30}
	if _, err := svc.RefundTx(nil, 3, 30, "取消订单", nil); err != nil {
		t.Fatalf("refund failed: %v", err)
	}
	account, _ := svc.GetByUserID(3)
	if account.Balance != 100 {
		t.Fatalf("balance after refund = %d, want 100", account.Balance)
	}
}
