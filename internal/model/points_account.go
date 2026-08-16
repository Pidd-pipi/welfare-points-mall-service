package model

import "time"

// PointsAccount 员工积分账户。
type PointsAccount struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	User         *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Balance      int       `gorm:"not null;default:0" json:"balance"`
	TotalEarned  int       `gorm:"not null;default:0" json:"total_earned"`
	TotalSpent   int       `gorm:"not null;default:0" json:"total_spent"`
	FrozenPoints int       `gorm:"not null;default:0" json:"frozen_points"`
	UpdatedAt    time.Time `json:"updated_at"`
}
