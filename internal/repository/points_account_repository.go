package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

// PointsAccountRepository 积分账户仓储。
type PointsAccountRepository interface {
	Create(account *model.PointsAccount) error
	CreateTx(tx *gorm.DB, account *model.PointsAccount) error
	FindByUserID(userID uint) (*model.PointsAccount, error)
	FindByUserIDTx(tx *gorm.DB, userID uint) (*model.PointsAccount, error)
	FindByID(id uint) (*model.PointsAccount, error)
	Update(account *model.PointsAccount) error
	UpdateTx(tx *gorm.DB, account *model.PointsAccount) error
	CreateTransaction(transaction *model.PointsTransaction) error
	CreateTransactionTx(tx *gorm.DB, transaction *model.PointsTransaction) error
	ListTransactions(userID uint, page, pageSize int) ([]model.PointsTransaction, int64, error)
	MonthlyEarnSpend(userID uint, month string) (earned, spent int, err error)
}

type pointsAccountRepository struct {
	db *gorm.DB
}

// NewPointsAccountRepository 构造积分账户仓储。
func NewPointsAccountRepository(db *gorm.DB) PointsAccountRepository {
	return &pointsAccountRepository{db: db}
}

func (r *pointsAccountRepository) Create(account *model.PointsAccount) error {
	return r.CreateTx(nil, account)
}

func (r *pointsAccountRepository) CreateTx(tx *gorm.DB, account *model.PointsAccount) error {
	if err := dbOrTx(r.db, tx).Create(account).Error; err != nil {
		return fmt.Errorf("create points account: %w", err)
	}
	return nil
}

func (r *pointsAccountRepository) FindByUserID(userID uint) (*model.PointsAccount, error) {
	return r.FindByUserIDTx(nil, userID)
}

func (r *pointsAccountRepository) FindByUserIDTx(tx *gorm.DB, userID uint) (*model.PointsAccount, error) {
	var account model.PointsAccount
	err := dbOrTx(r.db, tx).Preload("User").Where("user_id = ?", userID).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find points account by user: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find points account by user: %w", err)
	}
	return &account, nil
}

func (r *pointsAccountRepository) FindByID(id uint) (*model.PointsAccount, error) {
	var account model.PointsAccount
	err := r.db.Preload("User").First(&account, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find points account by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find points account by id: %w", err)
	}
	return &account, nil
}

func (r *pointsAccountRepository) Update(account *model.PointsAccount) error {
	return r.UpdateTx(nil, account)
}

func (r *pointsAccountRepository) UpdateTx(tx *gorm.DB, account *model.PointsAccount) error {
	if err := dbOrTx(r.db, tx).Save(account).Error; err != nil {
		return fmt.Errorf("update points account: %w", err)
	}
	return nil
}

func (r *pointsAccountRepository) CreateTransaction(transaction *model.PointsTransaction) error {
	return r.CreateTransactionTx(nil, transaction)
}

func (r *pointsAccountRepository) CreateTransactionTx(tx *gorm.DB, transaction *model.PointsTransaction) error {
	if err := dbOrTx(r.db, tx).Create(transaction).Error; err != nil {
		return fmt.Errorf("create points transaction: %w", err)
	}
	return nil
}

func (r *pointsAccountRepository) ListTransactions(userID uint, page, pageSize int) ([]model.PointsTransaction, int64, error) {
	var list []model.PointsTransaction
	var total int64
	q := r.db.Model(&model.PointsTransaction{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count points transactions: %w", err)
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list points transactions: %w", err)
	}
	return list, total, nil
}

func (r *pointsAccountRepository) MonthlyEarnSpend(userID uint, month string) (int, int, error) {
	start := month + "-01"
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return 0, 0, fmt.Errorf("monthly earn spend parse month: %w", err)
	}
	end := startTime.AddDate(0, 1, 0).Format("2006-01-02")
	var earned, spent int64
	q := r.db.Model(&model.PointsTransaction{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, start, end)
	if err := q.Where("change_type IN ?", []string{"earn", "refund"}).Select("COALESCE(SUM(amount),0)").Scan(&earned).Error; err != nil {
		return 0, 0, fmt.Errorf("monthly earned: %w", err)
	}
	if err := q.Where("change_type = ?", "spend").Select("COALESCE(SUM(amount),0)").Scan(&spent).Error; err != nil {
		return 0, 0, fmt.Errorf("monthly spent: %w", err)
	}
	return int(earned), int(spent), nil
}
