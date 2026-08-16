package model

import (
	"time"

	"github.com/ld/welfaremall/internal/constants"
)

// Order 兑换订单。
type Order struct {
	ID                uint                  `gorm:"primaryKey" json:"id"`
	OrderNo           string                `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID            uint                  `gorm:"index;not null" json:"user_id"`
	User              *User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID         uint                  `gorm:"index;not null" json:"product_id"`
	Product           *Product              `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SeckillActivityID *uint                 `gorm:"index" json:"seckill_activity_id"`
	PointsCost        int                   `gorm:"not null" json:"points_cost"`
	Quantity          int                   `gorm:"not null" json:"quantity"`
	Status            constants.OrderStatus `gorm:"size:16;index;not null" json:"status"`
	ReceiverInfo      string                `gorm:"size:255" json:"receiver_info"`
	LogisticsCompany  string                `gorm:"size:64" json:"logistics_company"`
	LogisticsNo       string                `gorm:"size:64" json:"logistics_no"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}
