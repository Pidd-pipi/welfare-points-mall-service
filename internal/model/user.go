package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
)

// User 员工用户：HR 可配置 PointsRule，员工拥有 PointsAccount。
type User struct {
	ID           uint               `gorm:"primaryKey" json:"id"`
	Username     string             `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string             `gorm:"size:255;not null" json:"-"`
	RealName     string             `gorm:"size:64" json:"real_name"`
	Email        string             `gorm:"size:128" json:"email"`
	Phone        string             `gorm:"size:32" json:"phone"`
	Department   string             `gorm:"size:64" json:"department"`
	Role         constants.UserRole `gorm:"size:32;index;not null" json:"role"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// BeforeCreate 校验角色。
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if !u.Role.Valid() {
		return constants.ErrInvalidUserRole
	}
	return nil
}
