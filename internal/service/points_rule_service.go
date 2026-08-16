package service

import (
	"fmt"
	"log/slog"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

// PointsRuleService 积分规则业务逻辑。
type PointsRuleService interface {
	Create(name, ruleType string, points int, effectiveDate string, enabled bool, description string) (*model.PointsRule, error)
	Update(id uint, name, ruleType string, points int, effectiveDate string, enabled bool, description string) (*model.PointsRule, error)
	Delete(id uint) error
	List(page, pageSize int) ([]model.PointsRule, int64, error)
	Execute(id uint) (int, error)
}

type pointsRuleService struct {
	ruleRepo   repository.PointsRuleRepository
	userRepo   repository.UserRepository
	accountSvc PointsAccountService
	logger     *slog.Logger
}

// NewPointsRuleService 构造积分规则服务。
func NewPointsRuleService(ruleRepo repository.PointsRuleRepository, userRepo repository.UserRepository, accountSvc PointsAccountService, logger *slog.Logger) PointsRuleService {
	return &pointsRuleService{ruleRepo: ruleRepo, userRepo: userRepo, accountSvc: accountSvc, logger: logger}
}

func (s *pointsRuleService) Create(name, ruleType string, points int, effectiveDate string, enabled bool, description string) (*model.PointsRule, error) {
	if name == "" || points <= 0 {
		return nil, fmt.Errorf("create points rule: %w", util.ErrValidation)
	}
	rule := &model.PointsRule{Name: name, RuleType: ruleType, Points: points, EffectiveDate: effectiveDate, Enabled: enabled, Description: description}
	if err := s.ruleRepo.Create(rule); err != nil {
		return nil, fmt.Errorf("create points rule: %w", err)
	}
	s.logger.Info(constants.LogRuleCreated, "rule_id", rule.ID, "rule_type", ruleType, "points", points)
	return rule, nil
}

func (s *pointsRuleService) Update(id uint, name, ruleType string, points int, effectiveDate string, enabled bool, description string) (*model.PointsRule, error) {
	rule, err := s.ruleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update points rule[id=%d]: %w", id, err)
	}
	if name != "" {
		rule.Name = name
	}
	if ruleType != "" {
		rule.RuleType = ruleType
	}
	if points > 0 {
		rule.Points = points
	}
	if effectiveDate != "" {
		rule.EffectiveDate = effectiveDate
	}
	rule.Enabled = enabled
	if description != "" {
		rule.Description = description
	}
	if err := s.ruleRepo.Update(rule); err != nil {
		return nil, fmt.Errorf("update points rule[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogRuleUpdated, "rule_id", id)
	return rule, nil
}

func (s *pointsRuleService) Delete(id uint) error {
	if err := s.ruleRepo.Delete(id); err != nil {
		return fmt.Errorf("delete points rule[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogRuleDeleted, "rule_id", id)
	return nil
}

func (s *pointsRuleService) List(page, pageSize int) ([]model.PointsRule, int64, error) {
	rules, total, err := s.ruleRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list points rules: %w", err)
	}
	s.logger.Info(constants.LogRuleListQueried, "total", total)
	return rules, total, nil
}

// Execute 手动执行积分发放任务：给所有员工发放该规则对应积分。
func (s *pointsRuleService) Execute(id uint) (int, error) {
	rule, err := s.ruleRepo.FindByID(id)
	if err != nil {
		return 0, fmt.Errorf("execute points rule[id=%d]: %w", id, err)
	}
	if !rule.Enabled {
		return 0, fmt.Errorf("execute points rule[id=%d] disabled: %w", id, util.ErrConflict)
	}
	employees, err := s.userRepo.ListEmployees()
	if err != nil {
		return 0, fmt.Errorf("execute points rule[id=%d]: %w", id, err)
	}
	granted := 0
	for _, emp := range employees {
		if _, err := s.accountSvc.Grant(emp.ID, rule.Points, fmt.Sprintf("积分规则[%s]发放", rule.Name), nil); err != nil {
			s.logger.Warn(constants.LogRuleExecuteFailed, "rule_id", id, "user_id", emp.ID, "error", err)
			continue
		}
		granted++
	}
	s.logger.Info(constants.LogRuleExecuted, "rule_id", id, "granted", granted)
	return granted, nil
}
