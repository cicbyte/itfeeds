# ITFeeds

基于 GoFrame v2 + Vue 3 的全栈资讯聚合系统，支持 RSS 订阅同步、MCP AI 工具调用。

## 功能特性

- **RSS 订阅同步** — 内置定时任务，自动抓取 IT 之家 RSS 源
- **Web 初始化** — Docker 首次部署通过 Web 页面配置数据库，无需手动编辑配置
- **新闻列表** — 分页、搜索、时间范围筛选
- **MCP 接口** — AI 工具调用支持（5 个工具）
- **手动拉取** — 前端一键拉取最新 RSS 数据
- **Docker 部署** — GitHub Actions 自动构建镜像，推送至 GHCR

## 截图

**初始化页面**

![初始化页面](images/001_init.png)

**主页面**

![主页面](images/002_main.png)

## 快速开始

### Docker 部署（推荐）

```bash
# 1. 启动容器
docker compose up -d

# 2. 访问初始化页面配置数据库
# 浏览器打开 http://<ip>:<port>/init

# 3. 初始化完成后自动跳转到主页
```

### 本地开发

```bash
# 1. 安装依赖
go mod tidy
cd web && npm install

# 2. 复制并修改配置
cp manifest/config/config.example.yaml manifest/config/config.yaml
# 编辑数据库连接信息

# 3. 启动后端
gf run

# 4. 启动前端（另一个终端）
cd web && npm run dev
```

### 环境要求

- Go 1.23+
- Node.js 18+
- MySQL 5.7+

## 项目结构

```
itfeeds/
├── api/v1/                    # API 定义层
├── internal/
│   ├── cmd/                   # 命令入口
│   ├── controller/            # 控制器层
│   ├── dao/                   # 数据访问层
│   ├── logic/
│   │   ├── rss_entries/       # RSS 条目业务逻辑
│   │   ├── rss_sync/          # RSS 同步定时任务
│   │   └── init/              # 系统初始化
│   ├── mcp/                   # MCP 服务
│   ├── model/                 # 数据模型
│   ├── service/               # Service 接口
│   └── router/                # 路由注册
├── manifest/config/           # 配置文件
├── resource/
│   ├── public/                # 前端静态资源
│   └── sql/mysql/             # SQL 脚本
├── scripts/                   # 构建 & 部署脚本
├── web/                       # Vue 3 前端
├── Dockerfile
├── docker-compose.yaml
└── .github/workflows/         # GitHub Actions
```

## API 接口

### REST API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/health | 健康检查 |
| GET | /api/v1/rss_entries/list | 获取新闻列表 |
| GET | /api/v1/rss_entries/detail?id=N | 获取新闻详情 |
| POST | /api/v1/rss_entries/sync | 手动拉取 RSS |
| GET | /api/v1/init/status | 初始化状态 |
| POST | /api/v1/init/test-connection | 测试数据库连接 |
| POST | /api/v1/init/setup | 执行初始化 |

### MCP 工具

ITFeeds 内置 MCP（Model Context Protocol）服务，AI 客户端可通过该接口查询和搜索新闻数据。

**端点:** `POST /mcp`（Streamable HTTP 模式）

| 工具名 | 说明 | 参数 |
|--------|------|------|
| get_server_info | 获取服务器运行信息 | 无 |
| get_news_list | 获取新闻列表 | `page_num`(页码) `page_size`(每页数量) `title`(标题搜索) `start_date`(开始日期) `end_date`(结束日期) |
| get_news_detail | 获取新闻详情 | `id`*(必填)* |
| search_news | 搜索新闻标题 | `keyword`*(必填)* `limit`(返回数量) |
| get_statistics | 获取新闻统计 | 无 |

#### 配置方式

**Cherry Studio** — 在设置 → MCP 服务中添加，类型选择 `Streamable HTTP`，地址填 `http://localhost:8000/mcp`：

![Cherry Studio MCP 配置](images/003_mcp_config.png)

配置成功后，工具列表中会显示 ITFeeds 提供的 5 个工具：

![MCP 工具列表](images/004_mcp_tools.png)

> Docker 部署时将 `localhost:8000` 替换为实际地址。

## 配置说明

RSS 同步相关配置位于 `manifest/config/config.yaml`：

```yaml
rss:
  enabled: true                    # 是否开启定时同步
  crons:                           # 定时规则（6 位 cron：秒 分 时 日 月 周）
    - "0 0 8-21 * * *"             # 每天 8:00-21:00 每整点执行
    - "0 30 8-20 * * *"            # 每天 8:30-20:30 每半点执行
  feeds:                           # RSS 源列表，支持配置多个
    - "https://www.ithome.com/rss"
  barkPush: ""                     # Bark 推送 key（留空则不推送）
```

- **crons** — 使用 6 位 cron 表达式（比标准 5 位多一个秒位），支持数组配置多个时段
- **feeds** — 可添加任意 RSS 源地址，同步时会逐个拉取
- **barkPush** — 配置 [Bark](https://github.com/Finb/bark-server) 推送 key，同步完成后会推送新增数量通知

## Docker 镜像

镜像托管在 GitHub Container Registry，推送 `v*` tag 自动构建：

```bash
# 拉取最新镜像
docker pull ghcr.io/cicbyte/itfeeds:latest

# 指定版本
docker pull ghcr.io/cicbyte/itfeeds:1.0.0
```

## 技术栈

**后端** — GoFrame v2.10.0 / gcron / mcp-go

**前端** — Vue 3 / Vite / Ant Design Vue 4 / Pinia / Vue Router

## License

MIT
