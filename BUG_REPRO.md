# BUG_REPRO

## Bug 是什么
问题列表过滤与仓库查询复用底层数组，后续调用写回污染调用方持有的旧切片。

## 如何触发
问题列表过滤与仓库查询复用底层数组，后续调用写回污染调用方持有的旧切片。

复现命令（定向测试）：

```
go test ./internal/service -run '^TestQuestionListBackingStableP401$' -count=1
go test ./internal/repository -run '^TestQuestionRepoListBackingStableP402$' -count=1
go test ./internal/handler -run '^TestQuestionHandlerListBackingStableP403$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 列表过滤与仓库查询复用底层数组，后续调用写回污染调用方持有的旧切片
