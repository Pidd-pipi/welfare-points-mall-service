package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

// PrizeItem 奖池条目。
type PrizeItem struct {
	Name   string  `json:"name"`
	Level  int     `json:"level"`
	Rate   float64 `json:"rate"`
	Points int     `json:"points,omitempty"`
}

// LotteryService 抽奖活动业务逻辑。
type LotteryService interface {
	Create(name string, costPoints int, prizePool string, startTime, endTime time.Time) (*model.LotteryActivity, error)
	List(page, pageSize int) ([]model.LotteryActivity, int64, error)
	Draw(userID uint, activityID uint) (*model.LotteryRecord, error)
	ListRecords(activityID uint, page, pageSize int) ([]model.LotteryRecord, int64, error)
	ResolveStatus(a *model.LotteryActivity) string
}

type lotteryService struct {
	lotteryRepo repository.LotteryRepository
	accountSvc  PointsAccountService
	db          *gorm.DB
	logger      *slog.Logger
}

// NewLotteryService 构造抽奖活动服务。
func NewLotteryService(lotteryRepo repository.LotteryRepository, accountSvc PointsAccountService, db *gorm.DB, logger *slog.Logger) LotteryService {
	return &lotteryService{lotteryRepo: lotteryRepo, accountSvc: accountSvc, db: db, logger: logger}
}

func (s *lotteryService) Create(name string, costPoints int, prizePool string, startTime, endTime time.Time) (*model.LotteryActivity, error) {
	if name == "" || costPoints <= 0 || startTime.IsZero() || endTime.IsZero() || !endTime.After(startTime) {
		return nil, fmt.Errorf("create lottery: %w", util.ErrValidation)
	}
	activity := &model.LotteryActivity{Name: name, CostPoints: costPoints, PrizePool: prizePool, StartTime: startTime, EndTime: endTime, Status: "draft"}
	if err := s.lotteryRepo.Create(activity); err != nil {
		return nil, fmt.Errorf("create lottery: %w", err)
	}
	s.logger.Info(constants.LogLotteryCreated, "lottery_id", activity.ID, "name", name)
	return activity, nil
}

func (s *lotteryService) List(page, pageSize int) ([]model.LotteryActivity, int64, error) {
	list, total, err := s.lotteryRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list lotteries: %w", err)
	}
	for i := range list {
		list[i].Status = s.ResolveStatus(&list[i])
	}
	s.logger.Info(constants.LogLotteryListQueried, "total", total)
	return list, total, nil
}

func (s *lotteryService) ResolveStatus(a *model.LotteryActivity) string {
	now := time.Now()
	if now.Before(a.StartTime) {
		return "draft"
	}
	if now.After(a.EndTime) {
		return "ended"
	}
	return "active"
}

func (s *lotteryService) Draw(userID uint, activityID uint) (*model.LotteryRecord, error) {
	activity, err := s.lotteryRepo.FindByID(activityID)
	if err != nil {
		return nil, fmt.Errorf("lottery draw activity[id=%d]: %w", activityID, err)
	}
	if s.ResolveStatus(activity) != "active" {
		s.logger.Warn(constants.LogLotteryDrawFailed, "activity_id", activityID, "reason", "not active")
		return nil, fmt.Errorf("lottery draw activity[id=%d]: %w", activityID, util.ErrActivityNotActive)
	}
	prize := rollPrize(activity.PrizePool)
	record := &model.LotteryRecord{UserID: userID, ActivityID: activityID, PrizeName: prize.Name, PrizeLevel: prize.Level}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if _, deductErr := s.accountSvc.DeductTx(tx, userID, activity.CostPoints, "积分抽奖："+activity.Name, nil); deductErr != nil {
			return fmt.Errorf("lottery draw deduct points user[%d]: %w", userID, deductErr)
		}
		if createErr := s.lotteryRepo.CreateRecordTx(tx, record); createErr != nil {
			return fmt.Errorf("lottery draw create record: %w", createErr)
		}
		if prize.Points > 0 {
			if _, grantErr := s.accountSvc.GrantTx(tx, userID, prize.Points, "抽奖奖品积分："+prize.Name, nil); grantErr != nil {
				return fmt.Errorf("lottery draw grant prize points: %w", grantErr)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogLotteryDrawSuccess, "user_id", userID, "activity_id", activityID, "prize", prize.Name)
	return record, nil
}

func (s *lotteryService) ListRecords(activityID uint, page, pageSize int) ([]model.LotteryRecord, int64, error) {
	list, total, err := s.lotteryRepo.ListRecords(activityID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list lottery records: %w", err)
	}
	s.logger.Info(constants.LogLotteryRecordQueried, "activity_id", activityID, "total", total)
	return list, total, nil
}

// rollPrize 按概率抽取奖品。
func rollPrize(poolJSON string) PrizeItem {
	items := []PrizeItem{{Name: "谢谢参与", Level: 0, Rate: 100}}
	if poolJSON != "" {
		var parsed []PrizeItem
		if err := json.Unmarshal([]byte(poolJSON), &parsed); err == nil && len(parsed) > 0 {
			items = parsed
		}
	}
	r := rand.Float64() * 100
	cum := 0.0
	for _, item := range items {
		cum += item.Rate
		if r <= cum {
			return item
		}
	}
	return items[len(items)-1]
}
