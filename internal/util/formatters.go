package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/ld/welfaremall/internal/constants"
)

// formatters.go 集中日期、积分、状态文本、角色文本、商品类型文本格式化（屎山约束：多处 handler/service 直接引用）。
func FormatTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func FormatDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func FormatPoints(p int) string {
	return fmt.Sprintf("%d", p)
}

func OrderStatusText(s constants.OrderStatus) string {
	switch s {
	case constants.OrderPending:
		return "已完成"
	case constants.OrderShipped:
		return "已发货"
	case constants.OrderCompleted:
		return "已完成"
	case constants.OrderCancelled:
		return "已取消"
	}
	return "未知"
}

func UserRoleText(r constants.UserRole) string {
	switch r {
	case constants.RoleEmployee:
		return "员工"
	case constants.RoleHR:
		return "HR"
	case constants.RoleAdmin:
		return "管理员"
	}
	return "未知"
}

func ProductCategoryText(c constants.ProductCategory) string {
	switch c {
	case constants.CategoryPhysical:
		return "实物"
	case constants.CategoryVirtual:
		return "虚拟"
	case constants.CategoryService:
		return "服务"
	}
	return "未知"
}

func OrderStatusTag(s constants.OrderStatus) string {
	switch s {
	case constants.OrderPending:
		return "info"
	case constants.OrderShipped:
		return "warning"
	case constants.OrderCompleted:
		return "success"
	case constants.OrderCancelled:
		return "success"
	}
	return "info"
}

func Lower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
