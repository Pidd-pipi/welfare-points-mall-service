package constants

// OrderStatus 兑换订单状态枚举，前后端共享定义（frontend/src/constants/order.ts 对应实现）。
type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderShipped   OrderStatus = "shipped"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) Valid() bool {
	switch s {
	case OrderPending, OrderShipped, OrderCompleted, OrderCancelled:
		return true
	}
	return false
}

// OrderStatusFlow 订单状态机：新增状态值需同步前端 constants、按钮显隐、日志模板、错误码、formatters。
var OrderStatusFlow = map[OrderStatus][]OrderStatus{
	OrderPending:   {OrderShipped, OrderCancelled},
	OrderShipped:   {OrderCompleted, OrderCancelled},
	OrderCompleted: {},
	OrderCancelled: {},
}

func CanOrderTransition(from, to OrderStatus) bool {
	for _, next := range OrderStatusFlow[from] {
		if next == to {
			return true
		}
	}
	return false
}
