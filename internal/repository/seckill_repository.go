package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

// SeckillRepository 秒杀活动仓储。
type SeckillRepository interface {
	Create(activity *model.SeckillActivity) error
	FindByID(id uint) (*model.SeckillActivity, error)
	List(page, pageSize int) ([]model.SeckillActivity, int64, error)
	Update(activity *model.SeckillActivity) error
	DecrementStock(id uint) error
	DecrementStockTx(tx *gorm.DB, id uint) error
}

type seckillRepository struct {
	db *gorm.DB
}

// NewSeckillRepository 构造秒杀活动仓储。
func NewSeckillRepository(db *gorm.DB) SeckillRepository {
	return &seckillRepository{db: db}
}

func (r *seckillRepository) Create(activity *model.SeckillActivity) error {
	if err := r.db.Create(activity).Error; err != nil {
		return fmt.Errorf("create seckill: %w", err)
	}
	return nil
}

func (r *seckillRepository) FindByID(id uint) (*model.SeckillActivity, error) {
	var a model.SeckillActivity
	err := r.db.Preload("Product").First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find seckill by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find seckill by id: %w", err)
	}
	return &a, nil
}

func (r *seckillRepository) List(page, pageSize int) ([]model.SeckillActivity, int64, error) {
	var list []model.SeckillActivity
	var total int64
	q := r.db.Model(&model.SeckillActivity{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count seckills: %w", err)
	}
	if err := q.Preload("Product").Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list seckills: %w", err)
	}
	return list, total, nil
}

func (r *seckillRepository) Update(activity *model.SeckillActivity) error {
	if err := r.db.Save(activity).Error; err != nil {
		return fmt.Errorf("update seckill: %w", err)
	}
	return nil
}

func (r *seckillRepository) DecrementStock(id uint) error {
	return r.DecrementStockTx(nil, id)
}

func (r *seckillRepository) DecrementStockTx(tx *gorm.DB, id uint) error {
	res := dbOrTx(r.db, tx).Model(&model.SeckillActivity{}).Where("id = ? AND seckill_stock > 0", id).
		UpdateColumn("seckill_stock", gorm.Expr("seckill_stock - 1"))
	if res.Error != nil {
		return fmt.Errorf("decrement seckill stock: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("decrement seckill stock: %w", util.ErrConflict)
	}
	return nil
}
