# BUG 复现说明（welfare-mall__005）

## Bug 是什么
库存扣减与兑换积分计算方向错误：扣库存实际调用加库存方法，兑换多件商品的积分算成单价加数量。

## 如何触发
```bash
go test ./internal/service -run 'TestProductDecrementStock|TestExchangeCostBoundary' -count=1
```

## 错误信息
```
--- FAIL: TestProductDecrementStock
    product_stock_combination_test.go:56: stock = 13, want 7
--- FAIL: TestExchangeCostBoundary
    product_stock_combination_test.go:62: ExchangeCost(100,3) = 103, want 300
```
