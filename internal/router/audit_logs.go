package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
)

func registerAuditLogRoutes(v1 *gin.RouterGroup, h *handler.AuditLogHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	logs := v1.Group("/audit-logs", auth)
	logs.GET("", middleware.RequireRole(adminRoles...), h.List)
}
