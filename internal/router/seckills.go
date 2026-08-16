package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
)

func registerSeckillRoutes(v1 *gin.RouterGroup, h *handler.SeckillHandler, auth gin.HandlerFunc, adminRoles, employeeRoles []constants.UserRole, limiter gin.HandlerFunc) {
	seckills := v1.Group("/seckills", auth)
	seckills.GET("", h.List)
	seckills.POST("", middleware.RequireRole(adminRoles...), h.Create)
	seckills.POST("/:id/purchase", middleware.RequireRole(employeeRoles...), limiter, h.Purchase)
}
