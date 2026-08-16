package model

import "time"

// LotteryActivity 积分抽奖活动。
type LotteryActivity struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:128;not null" json:"name"`
	CostPoints int       `gorm:"not null" json:"cost_points"`
	PrizePool  string    `gorm:"type:text" json:"prize_pool"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `gorm:"size:16;index;not null" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
