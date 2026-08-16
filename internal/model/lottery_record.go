package model

import "time"

// LotteryRecord 抽奖记录/中奖名单。
type LotteryRecord struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	UserID     uint             `gorm:"index;not null" json:"user_id"`
	User       *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ActivityID uint             `gorm:"index;not null" json:"activity_id"`
	Activity   *LotteryActivity `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	PrizeName  string           `gorm:"size:128" json:"prize_name"`
	PrizeLevel int              `gorm:"not null;default:0" json:"prize_level"`
	CreatedAt  time.Time        `json:"created_at"`
}
