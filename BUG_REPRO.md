# BUG_REPRO

## Bug 是什么
时间轴节点按录音分组时，内层 map 未初始化，首次写入触发 assignment to entry in nil map。

## 如何触发
时间轴节点按录音分组时，内层 map 未初始化，首次写入触发 assignment to entry in nil map。

复现命令（定向测试）：

```
go test ./internal/service -run '^TestMarkerGroupNoPanicP301$' -count=1
go test ./internal/handler -run '^TestMarkerGroupEmptyListP302$' -count=1
go test ./internal/service -run '^TestMarkerGroupEmptyMapP303$' -count=1
go test ./internal/service -run '^TestMarkerGroupMissingProjectP304$' -count=1
go test ./internal/service -run '^TestMarkerListEmptySliceP305$' -count=1
go test ./internal/handler -run '^TestMarkerListEmptySliceP306$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 分组 map 内层切片零值未初始化，首次写入触发 assignment to entry in nil map panic
