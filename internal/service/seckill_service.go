package service

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

// SeckillService 秒杀活动业务逻辑。
type SeckillService interface {
	Create(productID uint, seckillPoints, seckillStock, limitPerUser int, startTime, endTime time.Time) (*model.SeckillActivity, error)
	List(page, pageSize int) ([]model.SeckillActivity, int64, error)
	Purchase(userID uint, seckillID uint) (*model.Order, error)
	ResolveStatus(a *model.SeckillActivity) string
}

type seckillService struct {
	seckillRepo repository.SeckillRepository
	orderRepo   repository.OrderRepository
	accountSvc  PointsAccountService
	productSvc  ProductService
	db          *gorm.DB
	logger      *slog.Logger
}

// NewSeckillService 构造秒杀活动服务。
func NewSeckillService(seckillRepo repository.SeckillRepository, orderRepo repository.OrderRepository, accountSvc PointsAccountService, productSvc ProductService, db *gorm.DB, logger *slog.Logger) SeckillService {
	return &seckillService{seckillRepo: seckillRepo, orderRepo: orderRepo, accountSvc: accountSvc, productSvc: productSvc, db: db, logger: logger}
}

func (s *seckillService) Create(productID uint, seckillPoints, seckillStock, limitPerUser int, startTime, endTime time.Time) (*model.SeckillActivity, error) {
	if seckillPoints <= 0 || seckillStock <= 0 || startTime.IsZero() || endTime.IsZero() || !endTime.After(startTime) {
		return nil, fmt.Errorf("create seckill: %w", util.ErrValidation)
	}
	activity := &model.SeckillActivity{
		ProductID: productID, SeckillPoints: seckillPoints, SeckillStock: seckillStock,
		StartTime: startTime, EndTime: endTime, LimitPerUser: limitPerUser, Status: "draft",
	}
	if err := s.seckillRepo.Create(activity); err != nil {
		return nil, fmt.Errorf("create seckill: %w", err)
	}
	s.logger.Info(constants.LogSeckillCreated, "seckill_id", activity.ID, "product", productID)
	return activity, nil
}

func (s *seckillService) List(page, pageSize int) ([]model.SeckillActivity, int64, error) {
	list, total, err := s.seckillRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list seckills: %w", err)
	}
	for i := range list {
		list[i].Status = s.ResolveStatus(&list[i])
	}
	s.logger.Info(constants.LogSeckillListQueried, "total", total)
	return list, total, nil
}

// ResolveStatus 根据时间计算活动状态（draft/active/ended）。
func (s *seckillService) ResolveStatus(a *model.SeckillActivity) string {
	now := time.Now()
	if now.Before(a.StartTime) {
		return "draft"
	}
	if now.After(a.EndTime) {
		s.logger.Info(constants.LogSeckillEnded, "seckill_id", a.ID)
		return "ended"
	}
	return "active"
}

func (s *seckillService) Purchase(userID uint, seckillID uint) (*model.Order, error) {
	activity, err := s.seckillRepo.FindByID(seckillID)
	if err != nil {
		return nil, fmt.Errorf("seckill purchase seckill[id=%d]: %w", seckillID, err)
	}
	if s.ResolveStatus(activity) != "active" {
		s.logger.Warn(constants.LogSeckillDeductFailed, "seckill_id", seckillID, "reason", "not active")
		return nil, fmt.Errorf("seckill purchase seckill[id=%d]: %w", seckillID, util.ErrActivityNotActive)
	}
	limit := activity.LimitPerUser
	if limit <= 0 {
		limit = 1
	}
	var order *model.Order
	err = s.db.Transaction(func(tx *gorm.DB) error {
		count, countErr := s.orderRepo.CountByUserAndSeckillTx(tx, userID, seckillID)
		if countErr != nil {
			return fmt.Errorf("seckill purchase count user[%d]: %w", userID, countErr)
		}
		if count >= int64(limit) {
			return fmt.Errorf("seckill purchase user[%d] seckill[%d]: %w", userID, seckillID, util.ErrExchangeLimit)
		}
		if _, deductErr := s.accountSvc.DeductTx(tx, userID, activity.SeckillPoints, "秒杀抢购："+activity.Product.Name, nil); deductErr != nil {
			return fmt.Errorf("seckill purchase deduct points user[%d]: %w", userID, deductErr)
		}
		if stockErr := s.seckillRepo.DecrementStockTx(tx, seckillID); stockErr != nil {
			return fmt.Errorf("seckill purchase seckill[id=%d]: %w", seckillID, util.ErrSoldOut)
		}
		if productErr := s.productSvc.DecrementStockTx(tx, activity.ProductID, 1); productErr != nil {
			return fmt.Errorf("seckill purchase product[id=%d]: %w", activity.ProductID, productErr)
		}
		order = &model.Order{
			OrderNo: generateOrderNo(), UserID: userID, ProductID: activity.ProductID,
			SeckillActivityID: &seckillID, PointsCost: activity.SeckillPoints, Quantity: 1,
			Status: constants.OrderPending, ReceiverInfo: "",
		}
		if createErr := s.orderRepo.CreateTx(tx, order); createErr != nil {
			return fmt.Errorf("seckill purchase create order: %w", createErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogOrderCreateSuccess, "order_id", order.ID, "user", userID, "seckill", seckillID)
	return order, nil
}
