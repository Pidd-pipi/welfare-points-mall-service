package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/dto"
	"github.com/ld/welfaremall/internal/service"
	"github.com/ld/welfaremall/internal/util"
)

// PointsRuleHandler 积分规则接口处理器。
type PointsRuleHandler struct {
	ruleSvc service.PointsRuleService
}

// NewPointsRuleHandler 构造积分规则处理器。
func NewPointsRuleHandler(ruleSvc service.PointsRuleService) *PointsRuleHandler {
	return &PointsRuleHandler{ruleSvc: ruleSvc}
}

// Create 新增规则。
func (h *PointsRuleHandler) Create(c *gin.Context) {
	var req dto.PointsRuleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	rule, err := h.ruleSvc.Create(req.Name, req.RuleType, req.Points, req.EffectiveDate, req.Enabled, req.Description)
	if err != nil {
		c.Error(fmt.Errorf("handler create points rule: %w", err))
		return
	}
	util.OK(c, rule)
}

// Update 修改规则。
func (h *PointsRuleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的规则ID", err))
		return
	}
	var req dto.PointsRuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	rule, err := h.ruleSvc.Update(uint(id), req.Name, req.RuleType, req.Points, req.EffectiveDate, enabled, req.Description)
	if err != nil {
		c.Error(fmt.Errorf("handler update points rule: %w", err))
		return
	}
	util.OK(c, rule)
}

// Delete 删除规则。
func (h *PointsRuleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的规则ID", err))
		return
	}
	if err := h.ruleSvc.Delete(uint(id)); err != nil {
		c.Error(fmt.Errorf("handler delete points rule: %w", err))
		return
	}
	util.OKMessage(c, "规则已删除", nil)
}

// List 规则列表。
func (h *PointsRuleHandler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	rules, total, err := h.ruleSvc.List(q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler list points rules: %w", err))
		return
	}
	util.OK(c, gin.H{"list": rules, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// Execute 手动执行发放。
func (h *PointsRuleHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的规则ID", err))
		return
	}
	granted, err := h.ruleSvc.Execute(uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler execute points rule: %w", err))
		return
	}
	util.OKMessage(c, fmt.Sprintf("%s，共发放 %d 人", constants.MsgRuleExecuteSuccess, granted), gin.H{"granted": granted})
}
