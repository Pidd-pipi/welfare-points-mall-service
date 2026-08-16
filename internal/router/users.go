package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/handler"
)

func registerUserRoutes(v1 *gin.RouterGroup, h *handler.UserHandler, auth gin.HandlerFunc) {
	users := v1.Group("/users", auth)
	users.GET("/me", h.Me)
	users.PUT("/me", h.UpdateProfile)
}
