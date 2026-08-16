package model

import "time"

// SeckillActivity 秒杀活动。
type SeckillActivity struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	Product       *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SeckillPoints int       `gorm:"not null" json:"seckill_points"`
	SeckillStock  int       `gorm:"not null" json:"seckill_stock"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	LimitPerUser  int       `gorm:"not null;default:1" json:"limit_per_user"`
	Status        string    `gorm:"size:16;index;not null" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
