# BUG_REPRO

## Bug 是什么
JWT 解析用 %v 断链丢失底层 sentinel，鉴权无法区分过期、格式错误与签名无效。

## 如何触发
JWT 解析用 %v 断链丢失底层 sentinel，鉴权无法区分过期、格式错误与签名无效。

复现命令（定向测试）：

```
go test ./internal/util -run '^TestParseTokenExpiredChainP201$' -count=1
go test ./internal/util -run '^TestParseTokenMalformedChainP202$' -count=1
go test ./internal/middleware -run '^TestAuthExpiredTokenDistinctMessageP203$' -count=1
go test ./internal/middleware -run '^TestAuthSignatureInvalidDistinctMessageP204$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 JWT 解析用 %v 断链丢失底层 sentinel，errors.Is 失效，鉴权把过期与格式错误混为一谈
