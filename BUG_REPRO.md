# BUG_REPRO

## Bug 是什么
用户列表分页复用底层数组，后续分页覆盖调用方持有的前一页切片。

## 如何触发
用户列表分页复用底层数组，后续分页覆盖调用方持有的前一页切片。

复现命令（定向测试）：

```
go test ./internal/service -run '^TestUserListBackingStableP801$' -count=1
go test ./internal/repository -run '^TestUserRepoListBackingStableP802$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 分页查询复用底层数组，后续分页写回污染调用方持有的旧切片
