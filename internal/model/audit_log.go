package model

import "time"

// AuditLog 操作审计日志。
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:64" json:"username"`
	Action    string    `gorm:"size:64;index" json:"action"`
	Module    string    `gorm:"size:64" json:"module"`
	Detail    string    `gorm:"size:500" json:"detail"`
	IP        string    `gorm:"size:64" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
