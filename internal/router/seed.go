package router

import (
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
)

// seedData 初始化种子数据：用户、积分账户、规则、商品、活动、订单。
func seedData(db *gorm.DB, logger *slog.Logger) error {
	hash := func(pwd string) string {
		h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("seed hash failed", "error", err)
			return ""
		}
		return string(h)
	}

	users := []model.User{
		{Username: "admin", PasswordHash: hash("admin123"), RealName: "系统管理员", Role: constants.RoleAdmin, Department: "信息部"},
		{Username: "hr", PasswordHash: hash("hr123456"), RealName: "HR 小姐姐", Role: constants.RoleHR, Department: "人力资源部"},
		{Username: "employee1", PasswordHash: hash("emp123456"), RealName: "张伟", Role: constants.RoleEmployee, Department: "技术部"},
		{Username: "employee2", PasswordHash: hash("emp123456"), RealName: "李娜", Role: constants.RoleEmployee, Department: "市场部"},
	}
	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	accounts := []model.PointsAccount{
		{UserID: users[0].ID, Balance: 10000, TotalEarned: 10000},
		{UserID: users[1].ID, Balance: 8000, TotalEarned: 8000},
		{UserID: users[2].ID, Balance: 2000, TotalEarned: 3000, TotalSpent: 1000},
		{UserID: users[3].ID, Balance: 1500, TotalEarned: 2000, TotalSpent: 500},
	}
	if err := db.Create(&accounts).Error; err != nil {
		return fmt.Errorf("seed accounts: %w", err)
	}

	rules := []model.PointsRule{
		{Name: "全勤奖励", RuleType: "attendance_full", Points: 200, EffectiveDate: "2026-08-01", Enabled: true, Description: "当月全勤奖励 200 积分"},
		{Name: "绩效优秀", RuleType: "performance_excellent", Points: 500, EffectiveDate: "2026-08-01", Enabled: true, Description: "季度绩效优秀奖励"},
		{Name: "周年纪念", RuleType: "anniversary", Points: 1000, EffectiveDate: "2026-08-15", Enabled: true, Description: "入职周年纪念奖励"},
		{Name: "节日福利", RuleType: "festival", Points: 300, EffectiveDate: "2026-09-10", Enabled: false, Description: "中秋节日福利"},
	}
	if err := db.Create(&rules).Error; err != nil {
		return fmt.Errorf("seed rules: %w", err)
	}

	products := []model.Product{
		{Name: "蓝牙耳机", Category: constants.CategoryPhysical, PointsCost: 5000, Stock: 50, ExchangeLimit: 2, Status: "on", Description: "无线蓝牙耳机"},
		{Name: "保温杯", Category: constants.CategoryPhysical, PointsCost: 1500, Stock: 100, ExchangeLimit: 3, Status: "on", Description: "316 不锈钢保温杯"},
		{Name: "视频会员月卡", Category: constants.CategoryVirtual, PointsCost: 800, Stock: 200, ExchangeLimit: 2, Status: "on", Description: "主流视频平台会员月卡"},
		{Name: "按摩体验券", Category: constants.CategoryService, PointsCost: 3000, Stock: 20, ExchangeLimit: 1, Status: "on", Description: "肩颈按摩 60 分钟体验券"},
	}
	if err := db.Create(&products).Error; err != nil {
		return fmt.Errorf("seed products: %w", err)
	}

	now := time.Now()
	seckills := []model.SeckillActivity{
		{ProductID: products[0].ID, SeckillPoints: 3000, SeckillStock: 10, StartTime: now.Add(-time.Hour), EndTime: now.Add(24 * time.Hour), LimitPerUser: 1, Status: "active"},
		{ProductID: products[2].ID, SeckillPoints: 400, SeckillStock: 20, StartTime: now.Add(-time.Hour), EndTime: now.Add(24 * time.Hour), LimitPerUser: 1, Status: "active"},
	}
	if err := db.Create(&seckills).Error; err != nil {
		return fmt.Errorf("seed seckills: %w", err)
	}

	lotteries := []model.LotteryActivity{
		{
			Name: "幸运积分转盘", CostPoints: 100,
			PrizePool: `[{"name":"一等奖 5000积分","level":1,"rate":1,"points":5000},{"name":"二等奖 1000积分","level":2,"rate":4,"points":1000},{"name":"三等奖 200积分","level":3,"rate":15,"points":200},{"name":"谢谢参与","level":0,"rate":80}]`,
			StartTime: now.Add(-time.Hour), EndTime: now.Add(72 * time.Hour), Status: "active",
		},
	}
	if err := db.Create(&lotteries).Error; err != nil {
		return fmt.Errorf("seed lotteries: %w", err)
	}

	orders := []model.Order{
		{OrderNo: "WM202608010001", UserID: users[2].ID, ProductID: products[1].ID, PointsCost: 1500, Quantity: 1, Status: constants.OrderCompleted, ReceiverInfo: "张伟 13800000001 北京市朝阳区"},
		{OrderNo: "WM202608020002", UserID: users[3].ID, ProductID: products[2].ID, PointsCost: 800, Quantity: 1, Status: constants.OrderPending, ReceiverInfo: "李娜 13800000002 上海市浦东新区"},
	}
	if err := db.Create(&orders).Error; err != nil {
		return fmt.Errorf("seed orders: %w", err)
	}
	logger.Info("seed data created", "users", len(users), "products", len(products), "orders", len(orders))
	return nil
}
