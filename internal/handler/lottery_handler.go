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

// LotteryHandler 抽奖活动接口处理器。
type LotteryHandler struct {
	lotterySvc service.LotteryService
}

// NewLotteryHandler 构造抽奖活动处理器。
func NewLotteryHandler(lotterySvc service.LotteryService) *LotteryHandler {
	return &LotteryHandler{lotterySvc: lotterySvc}
}

// Create 创建抽奖活动（管理员）。
func (h *LotteryHandler) Create(c *gin.Context) {
	var req dto.LotteryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	activity, err := h.lotterySvc.Create(req.Name, req.CostPoints, req.PrizePool, req.StartTime, req.EndTime)
	if err != nil {
		c.Error(fmt.Errorf("handler create lottery: %w", err))
		return
	}
	util.OK(c, activity)
}

// List 抽奖活动列表。
func (h *LotteryHandler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	list, total, err := h.lotterySvc.List(q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler list lotteries: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// Draw 抽奖（限流）。
func (h *LotteryHandler) Draw(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的抽奖活动ID", err))
		return
	}
	record, err := h.lotterySvc.Draw(claims.UserID, uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler lottery draw: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgLotterySuccess, record)
}

// Records 中奖名单。
func (h *LotteryHandler) Records(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的抽奖活动ID", err))
		return
	}
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	list, total, err := h.lotterySvc.ListRecords(uint(id), q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler lottery records: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}
