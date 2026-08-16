package router

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ld/welfaremall/internal/config"
	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/service"
)

// Router 装配依赖并注册路由。
type Router struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

// NewRouter 初始化数据库、依赖并返回 gin 引擎。
func NewRouter(cfg *config.Config, logger *slog.Logger) (*gin.Engine, error) {
	r := &Router{cfg: cfg, logger: logger}
	if err := r.connectDB(); err != nil {
		return nil, err
	}
	if err := r.migrate(); err != nil {
		return nil, err
	}
	if err := r.seed(); err != nil {
		return nil, err
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.ErrorHandler())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     r.cfg.CORSOriginsSlice(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := engine.Group("/api")
	api.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.registerV1(api.Group("/v1"))
	return engine, nil
}

func (r *Router) connectDB() error {
	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(r.cfg.DSN()), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Warn)})
		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				sqlDB.SetMaxOpenConns(20)
				sqlDB.SetMaxIdleConns(5)
			}
			r.db = db
			r.logger.Info("database connected", "host", r.cfg.DBHost, "db", r.cfg.DBName)
			return nil
		}
		r.logger.Warn("database not ready, retrying", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("connect database: %w", err)
}

func (r *Router) migrate() error {
	if err := r.db.AutoMigrate(
		&model.User{}, &model.PointsRule{}, &model.PointsAccount{}, &model.PointsTransaction{},
		&model.Product{}, &model.SeckillActivity{}, &model.LotteryActivity{}, &model.LotteryRecord{},
		&model.Order{}, &model.AuditLog{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	r.logger.Info("database migrated")
	return nil
}

func (r *Router) seed() error {
	var count int64
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users for seed: %w", err)
	}
	if count > 0 {
		r.logger.Info("seed skipped, users exist")
		return nil
	}
	if err := seedData(r.db, r.logger); err != nil {
		return fmt.Errorf("seed data: %w", err)
	}
	r.logger.Info("seed data inserted")
	return nil
}

func (r *Router) registerV1(v1 *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(r.db)
	ruleRepo := repository.NewPointsRuleRepository(r.db)
	accountRepo := repository.NewPointsAccountRepository(r.db)
	productRepo := repository.NewProductRepository(r.db)
	orderRepo := repository.NewOrderRepository(r.db)
	seckillRepo := repository.NewSeckillRepository(r.db)
	lotteryRepo := repository.NewLotteryRepository(r.db)
	auditRepo := repository.NewAuditLogRepository(r.db)

	userSvc := service.NewUserService(userRepo, r.logger, r.cfg.JWTSecret, r.cfg.TokenTTLHours)
	accountSvc := service.NewPointsAccountService(accountRepo, r.db, r.logger)
	ruleSvc := service.NewPointsRuleService(ruleRepo, userRepo, accountSvc, r.logger)
	productSvc := service.NewProductService(productRepo, r.db, r.logger)
	orderSvc := service.NewOrderService(orderRepo, productSvc, accountSvc, r.db, r.logger)
	seckillSvc := service.NewSeckillService(seckillRepo, orderRepo, accountSvc, productSvc, r.db, r.logger)
	lotterySvc := service.NewLotteryService(lotteryRepo, accountSvc, r.db, r.logger)

	userHandler := handler.NewUserHandler(userSvc)
	ruleHandler := handler.NewPointsRuleHandler(ruleSvc)
	accountHandler := handler.NewPointsAccountHandler(accountSvc)
	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	seckillHandler := handler.NewSeckillHandler(seckillSvc)
	lotteryHandler := handler.NewLotteryHandler(lotterySvc)
	auditHandler := handler.NewAuditLogHandler(auditRepo)

	auth := middleware.AuthRequired(r.cfg)
	authLimiter := middleware.RateLimitStrict(r.cfg)
	adminRoles := []constants.UserRole{constants.RoleAdmin}
	hrRoles := []constants.UserRole{constants.RoleAdmin, constants.RoleHR}
	employeeRoles := []constants.UserRole{constants.RoleAdmin, constants.RoleHR, constants.RoleEmployee}
	audit := middleware.AuditLog(auditRepo)

	registerAuthRoutes(v1, userHandler, authLimiter)
	registerUserRoutes(v1, userHandler, auth)
	registerPointsAccountRoutes(v1, accountHandler, auth)
	registerPointsRuleRoutes(v1, ruleHandler, auth, hrRoles)
	registerProductRoutes(v1, productHandler, auth, adminRoles)
	registerSeckillRoutes(v1, seckillHandler, auth, adminRoles, employeeRoles, authLimiter)
	registerLotteryRoutes(v1, lotteryHandler, auth, adminRoles, employeeRoles, authLimiter)
	registerOrderRoutes(v1, orderHandler, auth, adminRoles, hrRoles, authLimiter)
	registerAuditLogRoutes(v1, auditHandler, auth, adminRoles)
	_ = audit
}
