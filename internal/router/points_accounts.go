package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/handler"
)

func registerPointsAccountRoutes(v1 *gin.RouterGroup, h *handler.PointsAccountHandler, auth gin.HandlerFunc) {
	account := v1.Group("/points-account", auth)
	account.GET("/me", h.Me)
	account.GET("/me/transactions", h.Transactions)
	account.GET("/me/stats", h.Stats)
}
