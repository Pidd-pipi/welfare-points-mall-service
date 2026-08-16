package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

// LotteryRepository 抽奖活动仓储。
type LotteryRepository interface {
	Create(activity *model.LotteryActivity) error
	FindByID(id uint) (*model.LotteryActivity, error)
	List(page, pageSize int) ([]model.LotteryActivity, int64, error)
	Update(activity *model.LotteryActivity) error
	CreateRecord(record *model.LotteryRecord) error
	CreateRecordTx(tx *gorm.DB, record *model.LotteryRecord) error
	ListRecords(activityID uint, page, pageSize int) ([]model.LotteryRecord, int64, error)
}

type lotteryRepository struct {
	db *gorm.DB
}

// NewLotteryRepository 构造抽奖活动仓储。
func NewLotteryRepository(db *gorm.DB) LotteryRepository {
	return &lotteryRepository{db: db}
}

func (r *lotteryRepository) Create(activity *model.LotteryActivity) error {
	if err := r.db.Create(activity).Error; err != nil {
		return fmt.Errorf("create lottery: %w", err)
	}
	return nil
}

func (r *lotteryRepository) FindByID(id uint) (*model.LotteryActivity, error) {
	var a model.LotteryActivity
	err := r.db.First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find lottery by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find lottery by id: %w", err)
	}
	return &a, nil
}

func (r *lotteryRepository) List(page, pageSize int) ([]model.LotteryActivity, int64, error) {
	var list []model.LotteryActivity
	var total int64
	q := r.db.Model(&model.LotteryActivity{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count lotteries: %w", err)
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list lotteries: %w", err)
	}
	return list, total, nil
}

func (r *lotteryRepository) Update(activity *model.LotteryActivity) error {
	if err := r.db.Save(activity).Error; err != nil {
		return fmt.Errorf("update lottery: %w", err)
	}
	return nil
}

func (r *lotteryRepository) CreateRecord(record *model.LotteryRecord) error {
	return r.CreateRecordTx(nil, record)
}

func (r *lotteryRepository) CreateRecordTx(tx *gorm.DB, record *model.LotteryRecord) error {
	if err := dbOrTx(r.db, tx).Create(record).Error; err != nil {
		return fmt.Errorf("create lottery record: %w", err)
	}
	return nil
}

func (r *lotteryRepository) ListRecords(activityID uint, page, pageSize int) ([]model.LotteryRecord, int64, error) {
	var list []model.LotteryRecord
	var total int64
	q := r.db.Model(&model.LotteryRecord{}).Where("activity_id = ?", activityID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count lottery records: %w", err)
	}
	if err := q.Preload("User").Preload("Activity").Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list lottery records: %w", err)
	}
	return list, total, nil
}
