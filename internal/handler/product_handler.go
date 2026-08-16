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

// ProductHandler 商品接口处理器。
type ProductHandler struct {
	productSvc service.ProductService
}

// NewProductHandler 构造商品处理器。
func NewProductHandler(productSvc service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

// Create 新增商品。
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	product, err := h.productSvc.Create(req.Name, req.Category, req.PointsCost, req.Stock, req.ExchangeLimit, req.CoverImage, req.Description)
	if err != nil {
		c.Error(fmt.Errorf("handler create product: %w", err))
		return
	}
	util.OK(c, product)
}

// Update 修改商品。
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的商品ID", err))
		return
	}
	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	product, err := h.productSvc.Update(uint(id), req.Name, req.Category, req.PointsCost, req.Stock, req.ExchangeLimit, req.CoverImage, req.Description)
	if err != nil {
		c.Error(fmt.Errorf("handler update product: %w", err))
		return
	}
	util.OK(c, product)
}

// Delete 删除商品。
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的商品ID", err))
		return
	}
	if err := h.productSvc.Delete(uint(id)); err != nil {
		c.Error(fmt.Errorf("handler delete product: %w", err))
		return
	}
	util.OKMessage(c, "商品已删除", nil)
}

// Toggle 上下架。
func (h *ProductHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的商品ID", err))
		return
	}
	product, err := h.productSvc.ToggleStatus(uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler toggle product: %w", err))
		return
	}
	util.OK(c, product)
}

// List 商品列表。
func (h *ProductHandler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	category := constants.ProductCategory(c.Query("category"))
	keyword := c.Query("keyword")
	products, total, err := h.productSvc.List(q.Page, q.PageSize, category, keyword)
	if err != nil {
		c.Error(fmt.Errorf("handler list products: %w", err))
		return
	}
	util.OK(c, gin.H{"list": products, "total": total, "page": q.Page, "page_size": q.PageSize})
}
