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

// SeckillHandler 秒杀活动接口处理器。
type SeckillHandler struct {
	seckillSvc service.SeckillService
}

// NewSeckillHandler 构造秒杀活动处理器。
func NewSeckillHandler(seckillSvc service.SeckillService) *SeckillHandler {
	return &SeckillHandler{seckillSvc: seckillSvc}
}

// Create 创建秒杀活动（管理员）。
func (h *SeckillHandler) Create(c *gin.Context) {
	var req dto.SeckillCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	if req.LimitPerUser <= 0 {
		req.LimitPerUser = 1
	}
	activity, err := h.seckillSvc.Create(req.ProductID, req.SeckillPoints, req.SeckillStock, req.LimitPerUser, req.StartTime, req.EndTime)
	if err != nil {
		c.Error(fmt.Errorf("handler create seckill: %w", err))
		return
	}
	util.OK(c, activity)
}

// List 秒杀活动列表。
func (h *SeckillHandler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	list, total, err := h.seckillSvc.List(q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler list seckills: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// Purchase 秒杀抢购（限流）。
func (h *SeckillHandler) Purchase(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的秒杀活动ID", err))
		return
	}
	order, err := h.seckillSvc.Purchase(claims.UserID, uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler seckill purchase: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgSeckillSuccess, order)
}
