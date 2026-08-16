package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

// OrderRepository 订单仓储。
type OrderRepository interface {
	Create(order *model.Order) error
	CreateTx(tx *gorm.DB, order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByOrderNo(orderNo string) (*model.Order, error)
	List(userID uint, page, pageSize int, status constants.OrderStatus) ([]model.Order, int64, error)
	ListAll(page, pageSize int) ([]model.Order, int64, error)
	Update(order *model.Order) error
	UpdateTx(tx *gorm.DB, order *model.Order) error
	TransitionStatusTx(tx *gorm.DB, id uint, from, to constants.OrderStatus) error
	CountByUserAndProduct(userID, productID uint) (int64, error)
	CountByUserAndProductTx(tx *gorm.DB, userID, productID uint) (int64, error)
	CountByUserAndSeckill(userID uint, seckillID uint) (int64, error)
	CountByUserAndSeckillTx(tx *gorm.DB, userID uint, seckillID uint) (int64, error)
	MonthlySummary(month string) (*BillSummary, error)
}

// BillSummary 月度福利账单统计。
type BillSummary struct {
	TotalOrders    int64 `json:"total_orders"`
	TotalPoints    int64 `json:"total_points"`
	ShippedCount   int64 `json:"shipped_count"`
	CompletedCount int64 `json:"completed_count"`
	CancelledCount int64 `json:"cancelled_count"`
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 构造订单仓储。
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.CreateTx(nil, order)
}

func (r *orderRepository) CreateTx(tx *gorm.DB, order *model.Order) error {
	if err := dbOrTx(r.db, tx).Create(order).Error; err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	return nil
}

func (r *orderRepository) FindByID(id uint) (*model.Order, error) {
	var o model.Order
	err := r.db.Preload("User").Preload("Product").First(&o, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find order by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find order by id: %w", err)
	}
	return &o, nil
}

func (r *orderRepository) FindByOrderNo(orderNo string) (*model.Order, error) {
	var o model.Order
	err := r.db.Where("order_no = ?", orderNo).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find order by order_no: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find order by order_no: %w", err)
	}
	return &o, nil
}

func (r *orderRepository) List(userID uint, page, pageSize int, status constants.OrderStatus) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	q := r.db.Model(&model.Order{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}
	if err := q.Preload("User").Preload("Product").Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	return orders, total, nil
}

func (r *orderRepository) ListAll(page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	q := r.db.Model(&model.Order{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count all orders: %w", err)
	}
	if err := q.Preload("User").Preload("Product").Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("list all orders: %w", err)
	}
	return orders, total, nil
}

func (r *orderRepository) Update(order *model.Order) error {
	return r.UpdateTx(nil, order)
}

func (r *orderRepository) UpdateTx(tx *gorm.DB, order *model.Order) error {
	if err := dbOrTx(r.db, tx).Save(order).Error; err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	return nil
}

func (r *orderRepository) TransitionStatusTx(tx *gorm.DB, id uint, from, to constants.OrderStatus) error {
	res := dbOrTx(r.db, tx).Model(&model.Order{}).
		Where("id = ? AND status = ?", id, from).
		Update("status", to)
	if res.Error != nil {
		return fmt.Errorf("transition order status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("transition order status: %w", util.ErrConflict)
	}
	return nil
}

func (r *orderRepository) CountByUserAndProduct(userID, productID uint) (int64, error) {
	return r.CountByUserAndProductTx(nil, userID, productID)
}

func (r *orderRepository) CountByUserAndProductTx(tx *gorm.DB, userID, productID uint) (int64, error) {
	var count int64
	if err := dbOrTx(r.db, tx).Model(&model.Order{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count orders by user product: %w", err)
	}
	return count, nil
}

func (r *orderRepository) CountByUserAndSeckill(userID uint, seckillID uint) (int64, error) {
	return r.CountByUserAndSeckillTx(nil, userID, seckillID)
}

func (r *orderRepository) CountByUserAndSeckillTx(tx *gorm.DB, userID uint, seckillID uint) (int64, error) {
	var count int64
	if err := dbOrTx(r.db, tx).Model(&model.Order{}).Where("user_id = ? AND seckill_activity_id = ?", userID, seckillID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count orders by user seckill: %w", err)
	}
	return count, nil
}

func (r *orderRepository) MonthlySummary(month string) (*BillSummary, error) {
	start := month + "-01"
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, fmt.Errorf("monthly summary parse month: %w", err)
	}
	end := startTime.AddDate(0, 1, 0).Format("2006-01-02")
	summary := &BillSummary{}
	q := r.db.Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", start, end)
	if err := q.Count(&summary.TotalOrders).Error; err != nil {
		return nil, fmt.Errorf("monthly total orders: %w", err)
	}
	if err := q.Select("COALESCE(SUM(points_cost * quantity),0)").Scan(&summary.TotalPoints).Error; err != nil {
		return nil, fmt.Errorf("monthly total points: %w", err)
	}
	if err := q.Where("status = ?", constants.OrderShipped).Count(&summary.ShippedCount).Error; err != nil {
		return nil, fmt.Errorf("monthly shipped: %w", err)
	}
	if err := q.Where("status = ?", constants.OrderCompleted).Count(&summary.CompletedCount).Error; err != nil {
		return nil, fmt.Errorf("monthly completed: %w", err)
	}
	if err := q.Where("status = ?", constants.OrderCancelled).Count(&summary.CancelledCount).Error; err != nil {
		return nil, fmt.Errorf("monthly cancelled: %w", err)
	}
	return summary, nil
}
