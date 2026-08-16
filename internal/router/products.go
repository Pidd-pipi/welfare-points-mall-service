package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/handler"
	"github.com/ld/welfaremall/internal/middleware"
)

func registerProductRoutes(v1 *gin.RouterGroup, h *handler.ProductHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	products := v1.Group("/products", auth)
	products.GET("", h.List)
	products.POST("", middleware.RequireRole(adminRoles...), h.Create)
	products.PUT("/:id", middleware.RequireRole(adminRoles...), h.Update)
	products.DELETE("/:id", middleware.RequireRole(adminRoles...), h.Delete)
	products.PUT("/:id/toggle", middleware.RequireRole(adminRoles...), h.Toggle)
}
