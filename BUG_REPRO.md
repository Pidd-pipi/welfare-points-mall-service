# BUG 复现说明（welfare-mall__004）

## Bug 是什么
订单状态机与状态文案/标签多处错位：待发货订单被禁止取消，待发货文案显示为已完成，已取消订单标签显示为 success。

## 如何触发
```bash
go test ./internal/util -run 'TestOrderStateMachineAllowsCancel|TestOrderStatusTextBoundary' -count=1
```

## 错误信息
```
--- FAIL: TestOrderStateMachineAllowsCancel
    order_state_combination_test.go:11: pending -> cancelled should be allowed
--- FAIL: TestOrderStatusTextBoundary
    order_state_combination_test.go:20: OrderStatusText(pending) = 已完成, want 待发货
```
