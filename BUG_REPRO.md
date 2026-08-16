# BUG 复现说明（welfare-mall__002）

## Bug 是什么
用户查询与登录的 nil 处理缺失：仓库层未找到用户返回 (nil,nil)，服务层未判空，导致查用户返回空数据、登录未知账号直接 panic。

## 如何触发
```bash
go test ./internal/service -run 'TestGetByIDNilUserReturnsError|TestLoginNilUserUnauthorized' -count=1
```

## 错误信息
```
--- FAIL: TestGetByIDNilUserReturnsError
    user_nil_combination_test.go:31: expected error for nil user, got user=<nil>
panic: runtime error: invalid memory address or nil pointer dereference
```
