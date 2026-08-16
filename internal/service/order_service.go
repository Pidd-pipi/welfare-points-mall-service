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

// OrderService 兑换订单业务逻辑。
type OrderService interface {
	Exchange(userID, productID uint, quantity int, receiverInfo string) (*model.Order, error)
	Cancel(userID, orderID uint) (*model.Order, error)
	Ship(orderID uint, logisticsCompany, logisticsNo string) (*model.Order, error)
	Complete(orderID uint) (*model.Order, error)
	List(userID uint, page, pageSize int, status constants.OrderStatus) ([]model.Order, int64, error)
	Bill(month string) (*repository.BillSummary, error)
}

type orderService struct {
	orderRepo  repository.OrderRepository
	productSvc ProductService
	accountSvc PointsAccountService
	db         *gorm.DB
	logger     *slog.Logger
}

// NewOrderService 构造订单服务。
func NewOrderService(orderRepo repository.OrderRepository, productSvc ProductService, accountSvc PointsAccountService, db *gorm.DB, logger *slog.Logger) OrderService {
	return &orderService{orderRepo: orderRepo, productSvc: productSvc, accountSvc: accountSvc, db: db, logger: logger}
}

func (s *orderService) Exchange(userID, productID uint, quantity int, receiverInfo string) (*model.Order, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("exchange quantity[%d]: %w", quantity, util.ErrValidation)
	}
	product, err := s.productSvc.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("exchange product[id=%d]: %w", productID, err)
	}
	if product.Status != "on" {
		return nil, fmt.Errorf("exchange product[id=%d] off shelf: %w", productID, util.ErrConflict)
	}
	cost := util.ExchangeCost(product.PointsCost, quantity)

	var order *model.Order
	err = s.db.Transaction(func(tx *gorm.DB) error {
		count, countErr := s.orderRepo.CountByUserAndProductTx(tx, userID, productID)
		if countErr != nil {
			return fmt.Errorf("exchange count user[%d] product[%d]: %w", userID, productID, countErr)
		}
		if count+int64(quantity) > int64(product.ExchangeLimit) {
			return fmt.Errorf("exchange user[%d] product[%d] limit[%d]: %w", userID, productID, product.ExchangeLimit, util.ErrExchangeLimit)
		}
		if _, deductErr := s.accountSvc.DeductTx(tx, userID, cost, "积分兑换："+product.Name, nil); deductErr != nil {
			return fmt.Errorf("exchange deduct points user[%d]: %w", userID, deductErr)
		}
		if stockErr := s.productSvc.DecrementStockTx(tx, productID, quantity); stockErr != nil {
			return fmt.Errorf("exchange decrement stock product[id=%d]: %w", productID, stockErr)
		}
		order = &model.Order{
			OrderNo: generateOrderNo(), UserID: userID, ProductID: productID,
			PointsCost: product.PointsCost, Quantity: quantity, Status: constants.OrderPending,
			ReceiverInfo: receiverInfo,
		}
		if createErr := s.orderRepo.CreateTx(tx, order); createErr != nil {
			return fmt.Errorf("exchange create order: %w", createErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogOrderCreateSuccess, "order_id", order.ID, "user", userID, "product", productID, "cost", cost)
	return order, nil
}

func (s *orderService) Cancel(userID, orderID uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("cancel order[id=%d]: %w", orderID, err)
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("cancel order[id=%d] user[%d] not owner: %w", orderID, userID, util.ErrForbidden)
	}
	if !constants.CanOrderTransition(order.Status, constants.OrderCancelled) {
		return nil, fmt.Errorf("cancel order[id=%d] status[%s]: %w", orderID, order.Status, util.ErrConflict)
	}
	refund := util.RefundPoints(order.PointsCost, order.Quantity)
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if _, refundErr := s.accountSvc.RefundTx(tx, userID, refund, "取消订单返还："+order.Product.Name, &orderID); refundErr != nil {
			return fmt.Errorf("cancel order[id=%d] refund: %w", orderID, refundErr)
		}
		if restockErr := s.productSvc.IncrementStockTx(tx, order.ProductID, order.Quantity); restockErr != nil {
			return fmt.Errorf("cancel order[id=%d] restock: %w", orderID, restockErr)
		}
		if transitionErr := s.orderRepo.TransitionStatusTx(tx, orderID, order.Status, constants.OrderCancelled); transitionErr != nil {
			return fmt.Errorf("cancel order[id=%d] status[%s]: %w", orderID, order.Status, transitionErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	order.Status = constants.OrderCancelled
	s.logger.Info(constants.LogOrderCancelSuccess, "order_id", orderID, "refund", refund)
	return order, nil
}

func (s *orderService) Ship(orderID uint, logisticsCompany, logisticsNo string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("ship order[id=%d]: %w", orderID, err)
	}
	if !constants.CanOrderTransition(order.Status, constants.OrderShipped) {
		return nil, fmt.Errorf("ship order[id=%d] status[%s]: %w", orderID, order.Status, util.ErrConflict)
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if transitionErr := s.orderRepo.TransitionStatusTx(tx, orderID, order.Status, constants.OrderShipped); transitionErr != nil {
			return fmt.Errorf("ship order[id=%d] status[%s]: %w", orderID, order.Status, transitionErr)
		}
		order.Status = constants.OrderShipped
		order.LogisticsCompany = logisticsCompany
		order.LogisticsNo = logisticsNo
		if updateErr := s.orderRepo.UpdateTx(tx, order); updateErr != nil {
			return fmt.Errorf("ship order[id=%d]: %w", orderID, updateErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogOrderShipSuccess, "order_id", orderID, "logistics_no", logisticsNo)
	return order, nil
}

func (s *orderService) Complete(orderID uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("complete order[id=%d]: %w", orderID, err)
	}
	if !constants.CanOrderTransition(order.Status, constants.OrderCompleted) {
		return nil, fmt.Errorf("complete order[id=%d] status[%s]: %w", orderID, order.Status, util.ErrConflict)
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if transitionErr := s.orderRepo.TransitionStatusTx(tx, orderID, order.Status, constants.OrderCompleted); transitionErr != nil {
			return fmt.Errorf("complete order[id=%d] status[%s]: %w", orderID, order.Status, transitionErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	order.Status = constants.OrderCompleted
	s.logger.Info(constants.LogOrderCompleteSuccess, "order_id", orderID)
	return order, nil
}

func (s *orderService) List(userID uint, page, pageSize int, status constants.OrderStatus) ([]model.Order, int64, error) {
	orders, total, err := s.orderRepo.List(userID, page, pageSize, status)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	s.logger.Info(constants.LogOrderListQueried, "user", userID, "total", total)
	return orders, total, nil
}

func (s *orderService) Bill(month string) (*repository.BillSummary, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	summary, err := s.orderRepo.MonthlySummary(month)
	if err != nil {
		return nil, fmt.Errorf("welfare bill[%s]: %w", month, err)
	}
	s.logger.Info(constants.LogBillGenerated, "month", month, "orders", summary.TotalOrders)
	return summary, nil
}

// generateOrderNo 生成订单号。
func generateOrderNo() string {
	return fmt.Sprintf("WM%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000000)
}
