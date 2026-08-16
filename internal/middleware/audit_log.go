package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

// AuditLog 操作审计日志中间件：为已登录用户的写操作记录审计。
func AuditLog(logRepo repository.AuditLogRepository) gin.HandlerFunc {
	mutating := map[string]bool{"POST": true, "PUT": true, "DELETE": true, "PATCH": true}
	return func(c *gin.Context) {
		c.Next()
		if !mutating[c.Request.Method] {
			return
		}
		claims, err := CurrentUser(c)
		if err != nil {
			return
		}
		entry := &model.AuditLog{
			UserID:    claims.UserID,
			Username:  claims.Username,
			Action:    c.Request.Method + " " + strings.TrimPrefix(c.FullPath(), "/api/v1"),
			Module:    moduleOf(c.FullPath()),
			Detail:    c.Request.URL.Path,
			IP:        c.ClientIP(),
			CreatedAt: time.Now(),
		}
		if err := logRepo.Create(entry); err != nil {
			util.GetLogger().Warn(constants.LogAuditRecorded, "error", err)
		}
	}
}

func moduleOf(path string) string {
	trimmed := strings.TrimPrefix(path, "/api/v1/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 0 {
		return "misc"
	}
	return parts[0]
}
