# BUG_REPRO

## Bug 是什么
录音上传成功但数据库关联失败时，错误分支跳过对象清理，MinIO 堆积孤儿文件。

## 如何触发
录音上传成功但数据库关联失败时，错误分支跳过对象清理，MinIO 堆积孤儿文件。

复现命令（定向测试）：

```
go test ./internal/handler -run '^TestUploadAudioAttachFailCleanupP601$' -count=1
go test ./internal/service -run '^TestStorageRemoveEmptyKeyP602$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 上传成功但关联失败的错误分支跳过对象清理，defer 释放与错误分支漏释放交织导致孤儿对象堆积
