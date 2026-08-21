# BUG_REPRO

## Bug 是什么
录音状态机新增重试中间态后转换表漏边、回写旧状态，重试成功后一直停在 processing。

## 如何触发
录音状态机新增重试中间态后转换表漏边、回写旧状态，重试成功后一直停在 processing。

复现命令（定向测试）：

```
go test ./internal/constants -run '^TestRecordingTransitionProcessingToFailedP701$' -count=1
go test ./internal/constants -run '^TestRecordingValidRetryingP702$' -count=1
go test ./internal/service -run '^TestRecordingRetryStatusP703$' -count=1
go test ./internal/dto -run '^TestUpdateRecordingStatusAllowsRetryingP704$' -count=1
go test ./internal/service -run '^TestRecordingAttachAudioRetryingToReadyP705$' -count=1
go test ./internal/constants -run '^TestRecordingTransitionRetryingEdgesP706$' -count=1
```

## 错误信息
见测试失败输出；核心机制为 录音状态机新增重试中间态后转换表漏边、worker 回写旧状态，跨层状态错位
