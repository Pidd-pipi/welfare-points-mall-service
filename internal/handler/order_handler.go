package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/dto"
	"github.com/ld/welfaremall/internal/middleware"
	"github.com/ld/welfaremall/internal/service"
	"github.com/ld/welfaremall/internal/util"
)

// OrderHandler 订单接口处理器。
type OrderHandler struct {
	orderSvc service.OrderService
}

// NewOrderHandler 构造订单处理器。
func NewOrderHandler(orderSvc service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

// Create 兑换下单。
func (h *OrderHandler) Create(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	order, err := h.orderSvc.Exchange(claims.UserID, req.ProductID, req.Quantity, req.ReceiverInfo)
	if err != nil {
		c.Error(fmt.Errorf("handler create order: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgExchangeSuccess, order)
}

// Cancel 取消订单。
func (h *OrderHandler) Cancel(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的订单ID", err))
		return
	}
	order, err := h.orderSvc.Cancel(claims.UserID, uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler cancel order: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgOrderCancelSuccess, order)
}

// Ship 发货（管理员）。
func (h *OrderHandler) Ship(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的订单ID", err))
		return
	}
	var req dto.ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	order, err := h.orderSvc.Ship(uint(id), req.LogisticsCompany, req.LogisticsNo)
	if err != nil {
		c.Error(fmt.Errorf("handler ship order: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgOrderShipSuccess, order)
}

// Complete 完成订单（管理员）。
func (h *OrderHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的订单ID", err))
		return
	}
	order, err := h.orderSvc.Complete(uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler complete order: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgOrderCompleteSuccess, order)
}

// List 我的订单。
func (h *OrderHandler) List(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	status := constants.OrderStatus(c.Query("status"))
	orders, total, err := h.orderSvc.List(claims.UserID, q.Page, q.PageSize, status)
	if err != nil {
		c.Error(fmt.Errorf("handler list orders: %w", err))
		return
	}
	util.OK(c, gin.H{"list": orders, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// Bill 福利账单（HR/管理员）。
func (h *OrderHandler) Bill(c *gin.Context) {
	month := c.Query("month")
	summary, err := h.orderSvc.Bill(month)
	if err != nil {
		c.Error(fmt.Errorf("handler welfare bill: %w", err))
		return
	}
	util.OK(c, summary)
}
