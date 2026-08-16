package dto

import "time"

// SeckillCreateRequest 秒杀活动创建请求。
type SeckillCreateRequest struct {
	ProductID     uint      `json:"product_id" binding:"required"`
	SeckillPoints int       `json:"seckill_points" binding:"required,min=1"`
	SeckillStock  int       `json:"seckill_stock" binding:"required,min=1"`
	LimitPerUser  int       `json:"limit_per_user" binding:"min=1"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
}

// LotteryCreateRequest 抽奖活动创建请求。
type LotteryCreateRequest struct {
	Name       string    `json:"name" binding:"required,max=128"`
	CostPoints int       `json:"cost_points" binding:"required,min=1"`
	PrizePool  string    `json:"prize_pool"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
}
