package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

// PointsRuleRepository 积分规则仓储。
type PointsRuleRepository interface {
	Create(rule *model.PointsRule) error
	FindByID(id uint) (*model.PointsRule, error)
	List(page, pageSize int) ([]model.PointsRule, int64, error)
	ListEnabled() ([]model.PointsRule, error)
	Update(rule *model.PointsRule) error
	Delete(id uint) error
}

type pointsRuleRepository struct {
	db *gorm.DB
}

// NewPointsRuleRepository 构造积分规则仓储。
func NewPointsRuleRepository(db *gorm.DB) PointsRuleRepository {
	return &pointsRuleRepository{db: db}
}

func (r *pointsRuleRepository) Create(rule *model.PointsRule) error {
	if err := r.db.Create(rule).Error; err != nil {
		return fmt.Errorf("create points rule: %w", err)
	}
	return nil
}

func (r *pointsRuleRepository) FindByID(id uint) (*model.PointsRule, error) {
	var rule model.PointsRule
	err := r.db.First(&rule, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find points rule by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find points rule by id: %w", err)
	}
	return &rule, nil
}

func (r *pointsRuleRepository) List(page, pageSize int) ([]model.PointsRule, int64, error) {
	var rules []model.PointsRule
	var total int64
	q := r.db.Model(&model.PointsRule{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count points rules: %w", err)
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("list points rules: %w", err)
	}
	return rules, total, nil
}

func (r *pointsRuleRepository) ListEnabled() ([]model.PointsRule, error) {
	var rules []model.PointsRule
	if err := r.db.Where("enabled = ?", true).Order("id desc").Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("list enabled points rules: %w", err)
	}
	return rules, nil
}

func (r *pointsRuleRepository) Update(rule *model.PointsRule) error {
	if err := r.db.Save(rule).Error; err != nil {
		return fmt.Errorf("update points rule: %w", err)
	}
	return nil
}

func (r *pointsRuleRepository) Delete(id uint) error {
	res := r.db.Delete(&model.PointsRule{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete points rule: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete points rule: %w", util.ErrNotFound)
	}
	return nil
}
