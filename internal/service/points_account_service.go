package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

// PointsAccountService 积分账户业务逻辑。
type PointsAccountService interface {
	GetOrCreate(userID uint) (*model.PointsAccount, error)
	GetByUserID(userID uint) (*model.PointsAccount, error)
	Grant(userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error)
	GrantTx(tx *gorm.DB, userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error)
	Deduct(userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error)
	DeductTx(tx *gorm.DB, userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error)
	Refund(userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error)
	RefundTx(tx *gorm.DB, userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error)
	ListTransactions(userID uint, page, pageSize int) ([]model.PointsTransaction, int64, error)
	MonthlyStats(userID uint, month string) (*MonthlyPointsStats, error)
}

// MonthlyPointsStats 月度积分收支。
type MonthlyPointsStats struct {
	Earned int `json:"earned"`
	Spent  int `json:"spent"`
}

type pointsAccountService struct {
	accountRepo repository.PointsAccountRepository
	db          *gorm.DB
	logger      *slog.Logger
}

// NewPointsAccountService 构造积分账户服务。
func NewPointsAccountService(accountRepo repository.PointsAccountRepository, db *gorm.DB, logger *slog.Logger) PointsAccountService {
	return &pointsAccountService{accountRepo: accountRepo, db: db, logger: logger}
}

func (s *pointsAccountService) GetOrCreate(userID uint) (*model.PointsAccount, error) {
	return s.getOrCreateTx(nil, userID)
}

func (s *pointsAccountService) getOrCreateTx(tx *gorm.DB, userID uint) (*model.PointsAccount, error) {
	account, err := s.accountRepo.FindByUserIDTx(tx, userID)
	if err == nil {
		return account, nil
	}
	if !errors.Is(err, util.ErrNotFound) {
		return nil, fmt.Errorf("get or create account user[%d]: %w", userID, err)
	}
	account = &model.PointsAccount{UserID: userID, Balance: 0, TotalEarned: 0, TotalSpent: 0}
	if err := s.accountRepo.CreateTx(tx, account); err != nil {
		return nil, fmt.Errorf("get or create account user[%d]: %w", userID, err)
	}
	return account, nil
}

func (s *pointsAccountService) GetByUserID(userID uint) (*model.PointsAccount, error) {
	return s.GetOrCreate(userID)
}

func (s *pointsAccountService) Grant(userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error) {
	var account *model.PointsAccount
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		account, err = s.GrantTx(tx, userID, amount, description, relatedOrderID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *pointsAccountService) GrantTx(tx *gorm.DB, userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("grant points amount[%d]: %w", amount, util.ErrValidation)
	}
	account, err := s.getOrCreateTx(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("grant points user[%d]: %w", userID, err)
	}
	account.Balance = util.DeductPoints(account.Balance, amount)
	account.TotalEarned += amount
	if err := s.accountRepo.UpdateTx(tx, account); err != nil {
		return nil, fmt.Errorf("grant points user[%d]: %w", userID, err)
	}
	transaction := &model.PointsTransaction{
		AccountID: account.ID, UserID: userID, ChangeType: "earn", Amount: amount,
		BalanceAfter: account.Balance, Description: description, RelatedOrderID: relatedOrderID,
	}
	if err := s.accountRepo.CreateTransactionTx(tx, transaction); err != nil {
		return nil, fmt.Errorf("grant points transaction user[%d]: %w", userID, err)
	}
	s.logger.Info(constants.LogPointsGranted, "user_id", userID, "amount", amount)
	return account, nil
}

func (s *pointsAccountService) Deduct(userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error) {
	var account *model.PointsAccount
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		account, err = s.DeductTx(tx, userID, amount, description, relatedOrderID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *pointsAccountService) DeductTx(tx *gorm.DB, userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("deduct points amount[%d]: %w", amount, util.ErrValidation)
	}
	account, err := s.getOrCreateTx(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("deduct points user[%d]: %w", userID, err)
	}
	if !util.CanAfford(account.Balance, amount) {
		return nil, fmt.Errorf("deduct points user[%d] balance[%d] need[%d]: %w", userID, account.Balance, amount, util.ErrPointsNotEnough)
	}
	account.Balance = util.GrantPoints(account.Balance, amount)
	account.TotalSpent += amount
	if err := s.accountRepo.UpdateTx(tx, account); err != nil {
		return nil, fmt.Errorf("deduct points user[%d]: %w", userID, err)
	}
	transaction := &model.PointsTransaction{
		AccountID: account.ID, UserID: userID, ChangeType: "spend", Amount: amount,
		BalanceAfter: account.Balance, Description: description, RelatedOrderID: relatedOrderID,
	}
	if err := s.accountRepo.CreateTransactionTx(tx, transaction); err != nil {
		return nil, fmt.Errorf("deduct points transaction user[%d]: %w", userID, err)
	}
	s.logger.Info(constants.LogPointsDeducted, "user_id", userID, "amount", amount)
	return account, nil
}

func (s *pointsAccountService) Refund(userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error) {
	var account *model.PointsAccount
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		account, err = s.RefundTx(tx, userID, amount, description, relatedOrderID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *pointsAccountService) RefundTx(tx *gorm.DB, userID uint, amount int, description string, relatedOrderID *uint) (*model.PointsAccount, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("refund points amount[%d]: %w", amount, util.ErrValidation)
	}
	account, err := s.getOrCreateTx(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("refund points user[%d]: %w", userID, err)
	}
	account.Balance = util.DeductPoints(account.Balance, amount)
	if err := s.accountRepo.UpdateTx(tx, account); err != nil {
		return nil, fmt.Errorf("refund points user[%d]: %w", userID, err)
	}
	transaction := &model.PointsTransaction{
		AccountID: account.ID, UserID: userID, ChangeType: "refund", Amount: amount,
		BalanceAfter: account.Balance, Description: description, RelatedOrderID: relatedOrderID,
	}
	if err := s.accountRepo.CreateTransactionTx(tx, transaction); err != nil {
		return nil, fmt.Errorf("refund points transaction user[%d]: %w", userID, err)
	}
	s.logger.Info(constants.LogPointsRefunded, "user_id", userID, "amount", amount)
	return account, nil
}

func (s *pointsAccountService) ListTransactions(userID uint, page, pageSize int) ([]model.PointsTransaction, int64, error) {
	list, total, err := s.accountRepo.ListTransactions(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list points transactions: %w", err)
	}
	s.logger.Info(constants.LogAccountFetched, "user_id", userID, "total", total)
	return list, total, nil
}

func (s *pointsAccountService) MonthlyStats(userID uint, month string) (*MonthlyPointsStats, error) {
	earned, spent, err := s.accountRepo.MonthlyEarnSpend(userID, month)
	if err != nil {
		return nil, fmt.Errorf("monthly points stats user[%d]: %w", userID, err)
	}
	s.logger.Info(constants.LogAccountStatsFetched, "user_id", userID, "earned", earned, "spent", spent)
	return &MonthlyPointsStats{Earned: earned, Spent: spent}, nil
}
