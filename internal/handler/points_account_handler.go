package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/dto"
	"github.com/ld/welfaremall/internal/middleware"
	"github.com/ld/welfaremall/internal/service"
	"github.com/ld/welfaremall/internal/util"
)

// PointsAccountHandler 积分账户接口处理器。
type PointsAccountHandler struct {
	accountSvc service.PointsAccountService
}

// NewPointsAccountHandler 构造积分账户处理器。
func NewPointsAccountHandler(accountSvc service.PointsAccountService) *PointsAccountHandler {
	return &PointsAccountHandler{accountSvc: accountSvc}
}

// Me 我的积分账户。
func (h *PointsAccountHandler) Me(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	account, err := h.accountSvc.GetByUserID(claims.UserID)
	if err != nil {
		c.Error(fmt.Errorf("handler points account me: %v", err))
		return
	}
	util.OK(c, account)
}

// Transactions 我的积分流水。
func (h *PointsAccountHandler) Transactions(c *gin.Context) {
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
	list, total, err := h.accountSvc.ListTransactions(claims.UserID, q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler points transactions: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// Stats 本月积分收支统计。
func (h *PointsAccountHandler) Stats(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	stats, err := h.accountSvc.MonthlyStats(claims.UserID, month)
	if err != nil {
		c.Error(fmt.Errorf("handler points stats: %w", err))
		return
	}
	util.OK(c, stats)
}
