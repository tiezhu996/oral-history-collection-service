# BUG_REPRO

## Bug 是什么
项目统计接口用进程内 map 缓存计数，无锁读写且直接返回内部引用，并发请求发生 data race。

## 如何触发
项目统计接口用进程内 map 缓存计数，无锁读写且直接返回内部引用，并发请求发生 data race。

复现命令（定向测试）：

```
go test -race ./internal/service -run '^TestStatsConcurrentNoRaceP101$' -count=1
go test -race ./internal/service -run '^TestStatsSnapshotStableP102$' -count=1
go test -race ./internal/handler -run '^TestStatsHandlerConcurrentNoRaceP103$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 进程内统计缓存 map 无锁读写并直接返回内部引用，并发读写发生 data race 且快照被后续写入污染
