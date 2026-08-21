# BUG_REPRO

## Bug 是什么
配置加载与数据库初始化用 %v 断链丢失哨兵错误，启动阶段无法定位错误来源。

## 如何触发
配置加载与数据库初始化用 %v 断链丢失哨兵错误，启动阶段无法定位错误来源。

复现命令（定向测试）：

```
go test ./internal/config -run '^TestConfigLoadInvalidEnvChainP1001$' -count=1
go test ./internal/database -run '^TestSeedAdminChainErrP1002$' -count=1
go test ./internal/database -run '^TestDatabaseNewNilCfgChainErrP1003$' -count=1
go test ./internal/database -run '^TestMigrateChainErrP1004$' -count=1
go test ./internal/config -run '^TestConfigValidateMissingSecretP1005$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 配置加载与数据库初始化用 %v 断链丢失哨兵错误，errors.Is 失效导致启动阶段误判错误来源
