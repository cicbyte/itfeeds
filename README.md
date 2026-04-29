# ITFeeds

> 全栈 RSS 资讯聚合系统，内置 MCP AI 工具接口，支持 Web 一键初始化和 Docker 部署。

[English](README_EN.md) | **中文**

![Build](https://img.shields.io/github/actions/workflow/status/cicbyte/itfeeds/docker-publish.yml?branch=master)
![Go Report Card](https://goreportcard.com/badge/github.com/cicbyte/itfeeds)
![License](https://img.shields.io/github/license/cicbyte/itfeeds)
![Docker Pulls](https://img.shields.io/badge/ghcr.io-cicbyte%2Fitfeeds-blue)
![Stars](https://img.shields.io/github/stars/cicbyte/itfeeds?style=social)

## 目录

- [功能特性](#功能特性)
- [截图](#截图)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [API 接口](#api-接口)
- [MCP 配置](#mcp-配置)
- [配置说明](#配置说明)
- [Docker 镜像](#docker-镜像)
- [技术栈](#技术栈)
- [License](#license)

## 功能特性

- **RSS 多源同步** — 内置定时任务，支持配置多个 RSS 源，自动抓取和去重
- **Web 初始化** — Docker 首次部署通过浏览器配置数据库，零命令行操作
- **MCP AI 接口** — 内置 5 个 MCP 工具，支持 Cherry Studio 等 AI 客户端直接调用
- **手动拉取** — 前端一键触发 RSS 同步，无需等待定时任务
- **Bark 推送** — 同步完成后推送新增数量通知（可选）
- **Docker 部署** — GitHub Actions 自动构建镜像，推送至 GHCR

## 截图

<table>
<tr>
<td width="50%"><img src="images/001_init.png" alt="初始化页面"></td>
<td width="50%"><img src="images/002_main.png" alt="主页面"></td>
</tr>
<tr>
<td width="50%"><img src="images/003_mcp_config.png" alt="MCP 配置"></td>
<td width="50%"><img src="images/004_mcp_tools.png" alt="MCP 工具列表"></td>
</tr>
</table>

## 快速开始

### Docker 部署（推荐）

```bash
docker compose up -d
```

浏览器访问 `http://<ip>:<port>/init` 完成数据库配置，初始化完成后自动跳转主页。

### 本地开发

```bash
go mod tidy
cd web && npm install
cp manifest/config/config.example.yaml manifest/config/config.yaml
# 编辑 manifest/config/config.yaml 填入数据库连接信息
gf run
# 另一个终端
cd web && npm run dev
```

### 环境要求

- Go 1.23+
- Node.js 18+
- MySQL 5.7+（仅支持 MySQL）

## 项目结构

```
itfeeds/
├── api/v1/                    # API 定义
├── internal/
│   ├── cmd/                   # 启动入口
│   ├── controller/            # 控制器
│   ├── dao/                   # 数据访问（自动生成）
│   ├── logic/
│   │   ├── rss_entries/       # RSS 条目业务
│   │   ├── rss_sync/          # RSS 同步任务
│   │   └── init/              # 系统初始化
│   ├── mcp/                   # MCP 服务
│   ├── model/                 # 数据模型
│   ├── service/               # Service 接口
│   └── router/                # 路由注册
├── manifest/config/           # 配置文件
├── resource/
│   ├── public/                # 前端静态资源
│   └── sql/mysql/             # SQL 脚本
├── scripts/                   # 构建脚本
├── web/                       # Vue 3 前端
├── Dockerfile
├── docker-compose.yaml
└── .github/workflows/         # CI/CD
```

## API 接口

### REST API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/health | 健康检查 |
| GET | /api/v1/rss_entries/list | 新闻列表（分页） |
| GET | /api/v1/rss_entries/detail?id=N | 新闻详情 |
| POST | /api/v1/rss_entries/sync | 手动拉取 RSS |
| GET | /api/v1/init/status | 初始化状态 |
| POST | /api/v1/init/test-connection | 测试数据库连接 |
| POST | /api/v1/init/setup | 执行初始化 |

### MCP 工具

端点：`POST /mcp`（Streamable HTTP）

| 工具名 | 说明 | 参数 |
|--------|------|------|
| get_server_info | 服务器运行信息 | 无 |
| get_news_list | 新闻列表 | `page_num` `page_size` `title` `start_date` `end_date` |
| get_news_detail | 新闻详情 | `id` *(必填)* |
| search_news | 搜索新闻 | `keyword` *(必填)* `limit` |
| get_statistics | 新闻统计 | 无 |

## MCP 配置

在 AI 客户端中添加 MCP Server，类型选择 **Streamable HTTP**，地址填 `http://localhost:8000/mcp`。

> Docker 部署时将 `localhost:8000` 替换为实际地址。

## 配置说明

RSS 同步配置位于 `manifest/config/config.yaml`：

```yaml
rss:
  enabled: true                    # 是否开启定时同步
  crons:                           # 6 位 cron（秒 分 时 日 月 周）
    - "0 0 8-21 * * *"
    - "0 30 8-20 * * *"
  feeds:                           # RSS 源列表，支持多个
    - "https://www.ithome.com/rss"
  barkPush: ""                     # Bark 推送 key（留空不推送）
```

- **crons** — 6 位 cron 表达式，支持数组配置多个时段
- **feeds** — 可添加任意 RSS 源，同步时逐个拉取，单个失败不阻断
- **barkPush** — 配置 [Bark](https://github.com/Finb/bark-server) 推送 key

## Docker 镜像

镜像托管在 GitHub Container Registry，推送 `v*` tag 自动触发构建：

```bash
docker pull ghcr.io/cicbyte/itfeeds:latest
docker pull ghcr.io/cicbyte/itfeeds:1.0.0
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | GoFrame v2.10.0 / gcron / mcp-go |
| 前端 | Vue 3 / Vite / Ant Design Vue 4 / Pinia / Vue Router |
| 数据库 | MySQL 5.7+ |
| 部署 | Docker / GitHub Actions / GHCR |

## License

[MIT](LICENSE)
