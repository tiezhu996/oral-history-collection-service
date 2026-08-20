# cy-180 口述历史采集工具 · 执行验证报告

- 项目编号：cy-180
- 项目短名：oralhistory
- 项目路径：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/新奇领域项目提示词/cy-180`
- 验证时间：2026-08-17 03:10 ~ 03:26（Asia/Shanghai）
- 启动命令：`docker compose up -d --build`
- 关闭命令：`docker compose down -v --remove-orphans`

## 端口

| 服务 | 宿主机端口 | 容器内端口 |
| --- | --- | --- |
| 前端（Nginx） | 8180 | 80 |
| 后端（Go/Gin） | 9180 | 8080 |
| MySQL | 10180 | 3306 |
| Redis | 46323 | 6379 |
| MinIO API | 47026 | 9000 |
| MinIO Console | 47027 | 9001 |

## docker compose ps（全部 healthy）

```
NAME                   IMAGE                  COMMAND                  SERVICE    STATUS
oralhistory-backend    oralhistory-backend    "/app/server"            backend    Up (healthy)
oralhistory-db         mysql:8.0              "docker-entrypoint.s…"   db         Up (healthy)
oralhistory-frontend   oralhistory-frontend   "/docker-entrypoint.…"   frontend   Up
oralhistory-minio      minio/minio:latest     "/usr/bin/docker-ent…"   minio      Up (healthy)
oralhistory-redis      redis:7-alpine         "docker-entrypoint.s…"   redis      Up (healthy)
```

后端健康检查：`GET http://localhost:9180/healthz` → 200 `{"code":0,"message":"ok","data":{"db":true,"service":"oralhistory"}}`

## curl 接口冒烟清单（20/20 PASS）

| # | 方法 | 路径 | 状态码/业务码 | 说明 |
| --- | --- | --- | --- | --- |
| 1 | GET | /healthz | 200 | 健康检查 |
| 2 | POST | /api/v1/auth/login | 0 | 管理员登录（admin/admin123456） |
| 3 | GET | /api/v1/auth/me | 0 | 当前用户 |
| 4 | POST | /api/v1/auth/register | 0 | 注册采访员 |
| 5 | POST | /api/v1/projects | 0 | 创建采访项目 |
| 6 | GET | /api/v1/projects | 0 | 项目列表 |
| 7 | GET | /api/v1/projects/1 | 0 | 项目详情 |
| 8 | POST | /api/v1/projects/1/questions | 0 | 添加采访问题 |
| 9 | GET | /api/v1/projects/1/questions | 0 | 问题列表 |
| 10 | POST | /api/v1/recordings | 0 | 创建录音记录 |
| 11 | GET | /api/v1/recordings?project_id=1 | 0 | 录音列表（按项目，复用 service） |
| 12 | POST | /api/v1/recordings/1/audio | 0 | 上传音频到 MinIO |
| 13 | GET | /api/v1/recordings/1/audio | 200 | 音频流播放 |
| 14 | PUT | /api/v1/recordings/1/summary | 0 | 一句话摘要 |
| 15 | POST | /api/v1/timeline-markers | 0 | 标注时间轴节点 |
| 16 | GET | /api/v1/timeline-markers?project_id=1 | 0 | 时间轴节点列表（按项目，复用 service） |
| 17 | PUT | /api/v1/projects/1/status | 0 | 项目状态流转 draft→in_progress |
| 18 | GET | /api/v1/audit-logs | 0 | 审计日志（管理员） |
| 19 | GET | /api/v1/projects | 401 | 未带 token 被拒绝 |
| 20 | PUT | /api/v1/projects/1/status | 40902 | 非法状态流转被拒绝 |

## 浏览器验证（内置 playwright 包装脚本，独立 session `oralhistory`）

| 页面 | 验证点 | 截图 |
| --- | --- | --- |
| 登录页 | 渲染正常，登录/注册切换 | output/01_login_and_list.png |
| 登录后项目列表 | admin 登录成功，列表展示已有项目与状态徽标 | output/01_login_and_list.png |
| 新建采访项目 | UI 创建「渡江战役亲历者访谈」成功，列表变 2 项 | output/02_project_list_created.png |
| 项目详情页 | 项目信息、采访问题添加、时间线区域渲染 | output/03_project_detail.png |
| 采访工作台 | 项目下拉选择、问题列表渲染 | output/04_interview_workspace.png |
| 录音面板 | 选择问题后出现「开始录音」与已录片段区 | output/05_recorder_panel.png |
| 审计日志页 | 管理员可见，展示 project.create/status、recording.upload、marker.create 等审计记录 | output/06_audit_page.png |
| 时间线详情 | 片段展示问题关联、一句话摘要、时间轴节点（⏱ 00:06 · 讲到胡同捉迷藏）、音频播放器渲染并可点击 | output/07_timeline_detail.png |
| 采访工作台复检 | 重新进入正常 | output/08_interview_recheck.png |

## 代码质量检查

- `cd backend && go build ./...` ✅
- `go vet ./...` ✅
- `go test ./...` ✅（service 表驱动测试 + repository sqlmock 表驱动测试全部通过）
- `cd frontend && npm run build`（tsc -b && vite build）✅
- `docker compose config --quiet` ✅（中文目录名下正常）

## 修复过程摘要

1. Gin 路由通配符冲突：`/projects/:projectId/questions` 与 `/projects/:id` 冲突，统一改为 `:id`，同步修改 handler 参数解析。
2. 前端 TS 构建错误：补充 `@types/node`、`vite-env.d.ts`、组件缺失 import、未使用变量清理。
3. 前端 hash 路由下 401 拦截器使用 `window.location.href` 导致跳转异常，改为 `window.location.hash = '#/login'`。
4. stores 中未捕获的 Promise rejection 导致控制台报错，为 fetch 方法补充 try/catch。
5. axios 响应拦截器误判 Blob 音频响应，`responseType === 'blob'` 时跳过统一响应校验（否则 AudioPlayer 无法加载）。

## git commit hash

主提交：`f2c9d0b000cb64dd5ffad084e0329427fc0afbce`
RBAC 日志强化复验后提交：见 git log（`98a0cd4` docs + RBAC 强化提交）

## 偏离说明

- 无重大偏离。任务清单端口 8180/9180/10180/46323/47026 全部按表使用。
- 提示词原文要求"暴露端口 …数据库 5432:5432"，实际按任务清单统一改为 MySQL `10180:3306`（任务清单优先级更高）。
- 前端框架原提示词未指定，采用 React 18 + Vite + TypeScript（README 技术栈已按规范标注"原前端框架保持不变"）。
- Swagger 采用 OpenAPI 3.0 文件（backend/api/openapi.yaml）方式提供，未接入 gin-swagger 运行时 UI。
