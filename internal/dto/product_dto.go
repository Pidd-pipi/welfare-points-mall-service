package dto

import "github.com/ld/welfaremall/internal/constants"

// ProductCreateRequest 商品创建请求。
type ProductCreateRequest struct {
	Name          string                    `json:"name" binding:"required,max=128"`
	Category      constants.ProductCategory `json:"category" binding:"required"`
	PointsCost    int                       `json:"points_cost" binding:"required,min=1"`
	Stock         int                       `json:"stock" binding:"min=0"`
	ExchangeLimit int                       `json:"exchange_limit" binding:"min=1"`
	CoverImage    string                    `json:"cover_image" binding:"max=255"`
	Description   string                    `json:"description" binding:"max=500"`
}

// ProductUpdateRequest 商品更新请求。
type ProductUpdateRequest struct {
	Name          string                    `json:"name" binding:"max=128"`
	Category      constants.ProductCategory `json:"category"`
	PointsCost    int                       `json:"points_cost" binding:"min=1"`
	Stock         int                       `json:"stock" binding:"min=0"`
	ExchangeLimit int                       `json:"exchange_limit" binding:"min=1"`
	CoverImage    string                    `json:"cover_image" binding:"max=255"`
	Description   string                    `json:"description" binding:"max=500"`
}
