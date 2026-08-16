package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
)

func registerPointsRuleRoutes(v1 *gin.RouterGroup, h *handler.PointsRuleHandler, auth gin.HandlerFunc, hrRoles []constants.UserRole) {
	rules := v1.Group("/points-rules", auth)
	rules.GET("", h.List)
	rules.POST("", middleware.RequireRole(hrRoles...), h.Create)
	rules.PUT("/:id", middleware.RequireRole(hrRoles...), h.Update)
	rules.DELETE("/:id", middleware.RequireRole(hrRoles...), h.Delete)
	rules.POST("/:id/execute", middleware.RequireRole(hrRoles...), h.Execute)
}
