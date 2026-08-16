package model

import "time"

// PointsRule 积分发放规则：由 HR 配置，可被积分发放任务引用。
type PointsRule struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:128;not null" json:"name"`
	RuleType      string    `gorm:"size:32;index;not null" json:"rule_type"`
	Points        int       `gorm:"not null" json:"points"`
	EffectiveDate string    `gorm:"size:16" json:"effective_date"`
	Enabled       bool      `gorm:"not null;default:true" json:"enabled"`
	Description   string    `gorm:"size:255" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
