package model

import "time"

// PointsTransaction 积分收支流水。
type PointsTransaction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	AccountID      uint      `gorm:"index;not null" json:"account_id"`
	UserID         uint      `gorm:"index;not null" json:"user_id"`
	ChangeType     string    `gorm:"size:32;index;not null" json:"change_type"`
	Amount         int       `gorm:"not null" json:"amount"`
	BalanceAfter   int       `gorm:"not null" json:"balance_after"`
	Description    string    `gorm:"size:255" json:"description"`
	RelatedOrderID *uint     `gorm:"index" json:"related_order_id"`
	CreatedAt      time.Time `json:"created_at"`
}
