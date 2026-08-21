# BUG_REPRO

## Bug 是什么
审计日志按用户名聚合时内层 map 未初始化，首次写入触发 nil map panic。

## 如何触发
审计日志按用户名聚合时内层 map 未初始化，首次写入触发 nil map panic。

复现命令（定向测试）：

```
go test ./internal/util -run '^TestGroupAuditLogsNoPanicP901$' -count=1
go test ./internal/util -run '^TestGroupAuditLogsEmptyP902$' -count=1
go test ./internal/handler -run '^TestAuditHandlerGroupEmptyP903$' -count=1
go test ./internal/handler -run '^TestAuditHandlerListEmptyP904$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 审计日志分组 map 内层切片零值未初始化，首次写入触发 assignment to entry in nil map panic
