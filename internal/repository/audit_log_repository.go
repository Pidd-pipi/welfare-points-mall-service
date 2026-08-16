package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/model"
)

// AuditLogRepository 审计日志仓储。
type AuditLogRepository interface {
	Create(log *model.AuditLog) error
	List(page, pageSize int) ([]model.AuditLog, int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 构造审计日志仓储。
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *model.AuditLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepository) List(page, pageSize int) ([]model.AuditLog, int64, error) {
	var list []model.AuditLog
	var total int64
	q := r.db.Model(&model.AuditLog{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return list, total, nil
}
