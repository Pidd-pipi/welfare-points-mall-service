# BUG 复现说明（welfare-mall__003）

## Bug 是什么
分页参数归一化与商品搜索过滤同时失效：Page<=0 未纠正为 1，PageSize<=0 未纠正为默认值，仓库 offset 错乱，商品列表搜索关键词被丢弃。

## 如何触发
```bash
go test ./internal/dto -run 'TestListQueryNormalizeDefaults' -count=1
```

## 错误信息
```
--- FAIL: TestListQueryNormalizeDefaults
    pagination_combination_test.go:9: Page = 0, want 1
```
