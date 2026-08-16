package model

import (
	"time"

	"github.com/ld/welfaremall/internal/constants"
)

// Product 积分商品。
type Product struct {
	ID            uint                      `gorm:"primaryKey" json:"id"`
	Name          string                    `gorm:"size:128;not null" json:"name"`
	Category      constants.ProductCategory `gorm:"size:32;index;not null" json:"category"`
	PointsCost    int                       `gorm:"not null" json:"points_cost"`
	Stock         int                       `gorm:"not null;default:0" json:"stock"`
	ExchangeLimit int                       `gorm:"not null;default:1" json:"exchange_limit"`
	Status        string                    `gorm:"size:16;default:on" json:"status"`
	CoverImage    string                    `gorm:"size:255" json:"cover_image"`
	Description   string                    `gorm:"size:500" json:"description"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}
