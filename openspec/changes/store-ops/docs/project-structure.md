# 项目目录结构

## 整体结构

```
store_tool/
├── CLAUDE.md                    # AI 编码指导文件
├── openspec/                    # OpenSpec 变更管理
│   ├── config.yaml
│   └── changes/
│       └── store-ops/           # 当前变更
│
├── server/                      # 后端服务 (Go)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # 入口文件
│   ├── internal/
│   │   ├── config/              # 配置管理
│   │   ├── handler/             # HTTP 处理器
│   │   │   ├── auth.go
│   │   │   ├── store.go
│   │   │   ├── command.go
│   │   │   ├── software.go
│   │   │   ├── process.go
│   │   │   ├── resource.go
│   │   │   ├── desktop.go
│   │   │   └── audit.go
│   │   ├── middleware/          # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   ├── model/               # 数据模型
│   │   │   ├── admin.go
│   │   │   ├── store.go
│   │   │   ├── audit_log.go
│   │   │   ├── command_log.go
│   │   │   ├── software_task.go
│   │   │   ├── desktop_session.go
│   │   │   └── resource_history.go
│   │   ├── repository/          # 数据访问层
│   │   │   ├── admin_repo.go
│   │   │   ├── store_repo.go
│   │   │   ├── audit_log_repo.go
│   │   │   └── ...
│   │   ├── service/             # 业务逻辑层
│   │   │   ├── auth_service.go
│   │   │   ├── store_service.go
│   │   │   ├── command_service.go
│   │   │   ├── software_service.go
│   │   │   ├── process_service.go
│   │   │   ├── resource_service.go
│   │   │   ├── desktop_service.go
│   │   │   └── audit_service.go
│   │   ├── ws/                  # WebSocket 相关
│   │   │   ├── hub.go           # 连接管理
│   │   │   ├── client.go        # 客户端连接
│   │   │   ├── message.go       # 消息类型定义
│   │   │   └── handler.go       # 消息处理
│   │   └── pkg/                 # 内部工具包
│   │       ├── response/        # 统一响应
│   │       ├── crypto/          # 加密工具
│   │       └── device/          # 设备指纹
│   ├── pkg/                     # 公共工具包
│   │   ├── logger/
│   │   └── config/
│   ├── migrations/              # 数据库迁移
│   │   ├── 001_init.up.sql
│   │   └── 001_init.down.sql
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
│
├── web/                         # Web 管理平台 (Vue 3)
│   ├── src/
│   │   ├── api/                 # API 调用
│   │   │   ├── auth.ts
│   │   │   ├── store.ts
│   │   │   ├── command.ts
│   │   │   ├── software.ts
│   │   │   ├── process.ts
│   │   │   ├── resource.ts
│   │   │   ├── desktop.ts
│   │   │   └── audit.ts
│   │   ├── components/          # 通用组件
│   │   │   ├── Layout/
│   │   │   ├── StoreStatus/
│   │   │   ├── CommandDialog/
│   │   │   └── ResourceChart/
│   │   ├── views/               # 页面视图
│   │   │   ├── login/
│   │   │   ├── dashboard/
│   │   │   ├── stores/
│   │   │   ├── command/
│   │   │   ├── software/
│   │   │   ├── process/
│   │   │   ├── resource/
│   │   │   ├── desktop/
│   │   │   └── audit/
│   │   ├── stores/              # Pinia 状态管理
│   │   │   ├── auth.ts
│   │   │   ├── store.ts
│   │   │   └── ...
│   │   ├── router/              # 路由配置
│   │   │   └── index.ts
│   │   ├── utils/               # 工具函数
│   │   ├── styles/              # 样式文件
│   │   ├── App.vue
│   │   └── main.ts
│   ├── public/
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   └── .env
│
├── client/                      # 门店客户端 (Wails)
│   ├── main.go                  # Go 后端入口
│   ├── app.go                   # Wails 应用
│   ├── internal/
│   │   ├── config/              # 配置管理
│   │   ├── connection/          # WebSocket 连接
│   │   │   ├── client.go
│   │   │   ├── reconnect.go
│   │   │   └── heartbeat.go
│   │   ├── executor/            # 命令/脚本执行
│   │   │   ├── command.go
│   │   │   └── script.go
│   │   ├── software/            # 软件安装
│   │   │   ├── downloader.go
│   │   │   └── installer.go
│   │   ├── process/             # 进程管理
│   │   │   ├── list.go
│   │   │   └── kill.go
│   │   ├── service/             # 服务管理
│   │   │   ├── list.go
│   │   │   └── control.go
│   │   ├── resource/            # 资源采集
│   │   │   ├── cpu.go
│   │   │   ├── memory.go
│   │   │   └── disk.go
│   │   ├── rustdesk/            # RustDesk 集成
│   │   │   ├── embed.go
│   │   │   └── config.go
│   │   ├── device/              # 设备信息
│   │   │   ├── id.go
│   │   │   └── fingerprint.go
│   │   ├── storage/             # 本地存储
│   │   │   └── sqlite.go
│   │   └── autostart/           # 开机自启
│   │       └── windows.go
│   ├── frontend/                # Wails 前端 (Vue 3)
│   │   ├── src/
│   │   │   ├── App.vue
│   │   │   ├── main.ts
│   │   │   └── components/
│   │   │       └── StatusDisplay.vue
│   │   ├── index.html
│   │   ├── package.json
│   │   └── vite.config.ts
│   ├── build/                   # 构建配置
│   │   └── appicon.png
│   ├── wails.json               # Wails 配置
│   ├── go.mod
│   └── go.sum
│
├── docs/                        # 项目文档
│   ├── architecture.md
│   ├── api.md
│   ├── websocket.md
│   ├── database.md
│   └── project-structure.md
│
├── scripts/                     # 脚本
│   ├── build.sh                 # 构建脚本
│   └── deploy.sh                # 部署脚本
│
├── docker-compose.yml           # Docker 编排
├── Dockerfile.server            # 服务端镜像
├── Dockerfile.web               # Web 镜像
├── .gitignore
└── Makefile
```

## 模块职责

### server/ (Go 后端)

| 目录 | 职责 |
|------|------|
| cmd/server | 程序入口 |
| internal/config | 配置加载和验证 |
| internal/handler | HTTP 请求处理 |
| internal/middleware | 请求中间件 |
| internal/model | 数据模型定义 |
| internal/repository | 数据库操作 |
| internal/service | 业务逻辑 |
| internal/ws | WebSocket 连接和消息处理 |
| internal/pkg | 内部工具 |
| pkg | 可复用公共包 |
| migrations | 数据库迁移脚本 |

### web/ (Vue 3 前端)

| 目录 | 职责 |
|------|------|
| src/api | API 调用封装 |
| src/components | 可复用组件 |
| src/views | 页面视图 |
| src/stores | Pinia 状态管理 |
| src/router | 路由配置 |
| src/utils | 工具函数 |
| src/styles | 样式文件 |

### client/ (Wails 客户端)

| 目录 | 职责 |
|------|------|
| internal/config | 客户端配置 |
| internal/connection | WebSocket 连接管理 |
| internal/executor | 命令/脚本执行 |
| internal/software | 软件下载和安装 |
| internal/process | 进程管理 |
| internal/service | Windows 服务管理 |
| internal/resource | 资源数据采集 |
| internal/rustdesk | RustDesk 集成 |
| internal/device | 设备 ID 和指纹 |
| internal/storage | SQLite 本地存储 |
| internal/autostart | 开机自启 |
| frontend | Wails 前端界面 |

## 构建产物

```
dist/
├── server              # Linux 可执行文件
├── web/                # Web 静态文件
└── client/             # Windows 客户端
    └── store-ops-client.exe
```