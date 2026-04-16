# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

ITFeeds（itfeeds）— 基于 GoFrame v2.10.0 + Vue 3 的全栈资讯聚合系统，支持 RSS 订阅同步、MCP AI 工具调用。模块名 `github.com/cicbyte/itfeeds`，Go 1.23.2+。

## 常用命令

```bash
# 后端开发（热重载）
gf run

# 代码生成（修改数据库表结构后执行）
gf gen dao        # 生成 DAO/DO/Entity
gf gen service    # 生成 Service 接口
gf gen ctrl       # 生成 Controller

# 前端开发
cd web && npm run dev       # Vite 开发服务器 (port 5173，代理 /api → :8000)
cd web && npm run build     # 生产构建
cd web && npm run lint      # oxlint + eslint
cd web && npm run format    # prettier

# 构建 & 部署
python scripts/build.py             # gf pack + gf build → Linux 二进制
python scripts/build_docker.py      # Docker 镜像构建
python scripts/export_image.py      # 导出镜像 tar
docker-compose up -d --build
```

项目无测试套件。

## 架构

### 分层结构（GoFrame 标准）

```
API 层 (api/v1/)          — 请求/响应结构体、路由元信息（g.Meta tag）
    ↓
Controller (internal/controller/)  — 参数验证、调用 Service
    ↓
Service 接口 (internal/service/)   — 业务逻辑接口定义
    ↓
Logic 实现 (internal/logic/)       — 业务逻辑实现（init() 中 service.RegisterXxx()）
    ↓
DAO (internal/dao/)       — 数据库操作（gf gen dao 自动生成，勿手动改）
```

### 关键入口

- **`internal/cmd/cmd.go`** — HTTP 服务器启动、MCP 路由注册（绕过中间件）、SPA 回退
- **`internal/router/router.go`** — API 路由组（`/api/v1`），绑定 CORS 中间件和 Controller
- **`internal/logic/logic.go`** — 所有 logic 子包的空白导入注册（新增模块必须在此添加）

### MCP 服务（`internal/mcp/`）

基于 mcp-go，HTTP Streamable 模式，端点 `POST /mcp`。路由在 `cmd.go` 中直接绑定到 `ghttp.Request`，绕过 GoFrame 中间件链（避免 JSON-RPC 响应被包装）。5 个工具定义在 `tools.go`，处理函数以 `handle` 前缀命名。

### 前端

Vue 3 + Vite 7 + Ant Design Vue 4 + Pinia + Vue Router 5。`@` 别名指向 `web/src/`。开发时 Vite 代理 `/api` 到后端 `:8000`，`/mcp` 路径未配置代理（直连后端）。生产构建输出到 `resource/public/` 由 GoFrame 静态服务。

## 新增业务模块流程

1. `api/v1/<module>/<module>.go` — 定义 Req/Res 结构体（`g.Meta` 标注路由）
2. `internal/model/<module>.go` — 业务模型
3. `gf gen dao` — 生成 DAO
4. `internal/service/<module>.go` — Service 接口
5. `internal/logic/<module>/<module>.go` — 实现（`init()` 中注册）
6. **`internal/logic/logic.go`** — 添加 `_ "github.com/cicbyte/itfeeds/internal/logic/<module>"`（遗漏则模块不加载）
7. `internal/controller/<module>.go` — Controller
8. `internal/router/router.go` — `group.Bind(controller.Xxx)`

## 命名约定

| 类型 | 格式 | 示例 |
|------|------|------|
| Controller 变量 | 大写 | `RssEntries = rssEntriesController{}` |
| Service 接口 | `I` + 模块名 | `IRssEntries` |
| Logic 实现 | `s` + 模块名 | `sRssEntries` |

## 配置

主配置 `manifest/config/config.yaml`。RSS 同步使用 **6 位 cron**（秒 分 时 日 月 周），支持 `crons` 数组配置多个时段。日志输出到 `resource/log/`。

## 数据库约定

- 操作：`dao.RssEntries.Ctx(ctx)`，写入用 `do.RssEntries`
- 返回模型：`model.RssEntriesInfo`（Info 后缀）
- 字段名：`dao.RssEntries.Columns().Xxx`
- 内容保留原始 HTML，不做 stripHTML

## 错误处理

```go
import liberr "github.com/cicbyte/itfeeds/library/liberr"
liberr.ErrIsNil(ctx, err, "操作失败")  // 失败则 panic，被中间件捕获
```
