package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
)

func registerLotteryRoutes(v1 *gin.RouterGroup, h *handler.LotteryHandler, auth gin.HandlerFunc, adminRoles, employeeRoles []constants.UserRole, limiter gin.HandlerFunc) {
	lotteries := v1.Group("/lotteries", auth)
	lotteries.GET("", h.List)
	lotteries.POST("", middleware.RequireRole(adminRoles...), h.Create)
	lotteries.POST("/:id/draw", middleware.RequireRole(employeeRoles...), limiter, h.Draw)
	lotteries.GET("/:id/records", h.Records)
}
