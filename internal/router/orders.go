package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
)

func registerOrderRoutes(v1 *gin.RouterGroup, h *handler.OrderHandler, auth gin.HandlerFunc, adminRoles, hrRoles []constants.UserRole, limiter gin.HandlerFunc) {
	orders := v1.Group("/orders", auth)
	orders.GET("", h.List)
	orders.POST("", limiter, h.Create)
	orders.POST("/:id/cancel", h.Cancel)
	orders.PUT("/:id/ship", middleware.RequireRole(adminRoles...), h.Ship)
	orders.PUT("/:id/complete", middleware.RequireRole(adminRoles...), h.Complete)
	bills := v1.Group("/admin/bills", auth)
	bills.GET("", middleware.RequireRole(hrRoles...), h.Bill)
}
