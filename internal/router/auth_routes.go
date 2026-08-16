package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/handler"
)

func registerAuthRoutes(v1 *gin.RouterGroup, h *handler.UserHandler, limiter gin.HandlerFunc) {
	auth := v1.Group("/auth")
	auth.POST("/register", limiter, h.Register)
	auth.POST("/login", limiter, h.Login)
}
