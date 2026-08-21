# BUG_REPRO

## Bug 是什么
审计写入忽略 ctx.Err()，中间件用 context.Background 在请求结束后仍写审计。

## 如何触发
审计写入忽略 ctx.Err()，中间件用 context.Background 在请求结束后仍写审计。

复现命令（定向测试）：

```
go test ./internal/service -run '^TestAuditRecordActiveCancelledCtxP501$' -count=1
go test ./internal/middleware -run '^TestAuditMiddlewareCancelledCtxP502$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 审计写入忽略 ctx.Err()，中间件用 context.Background 在请求结束后仍写审计，取消不向下游传播
