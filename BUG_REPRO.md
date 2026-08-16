# BUG 复现说明（welfare-mall__001）

## Bug 是什么
积分账户流水方向错乱：发放积分余额不涨、扣减积分余额反而增加、退款积分余额变少，属于服务层调用方向与工具函数语义多处错位。

## 如何触发
```bash
go test ./internal/service -run 'TestPointsFlowGrantAndDeduct|TestPointsFlowRefund' -count=1
```

## 错误信息
```
--- FAIL: TestPointsFlowGrantAndDeduct
    points_flow_combination_test.go:19: balance after grant = 0, want 500
--- FAIL: TestPointsFlowRefund
    points_flow_combination_test.go:43: balance after refund = 70, want 100
```
