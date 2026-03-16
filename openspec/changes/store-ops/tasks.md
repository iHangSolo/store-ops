# 实施任务清单

## 1. 项目初始化

- [x] 1.1 创建项目根目录结构 (server/, web/, client/, docs/, scripts/)
- [x] 1.2 初始化 Go 后端模块 (server/go.mod, Makefile)
- [x] 1.3 初始化 Vue 3 前端项目 (web/package.json, vite.config.ts)
- [x] 1.4 初始化 Wails 客户端项目 (client/wails.json, go.mod)
- [x] 1.5 创建 Docker 编排文件 (docker-compose.yml, Dockerfile.server, Dockerfile.web)
- [x] 1.6 创建构建脚本 (scripts/build.sh, scripts/deploy.sh)
- [x] 1.7 创建 GitHub 仓库并推送代码

## 2. CI/CD 配置 (GitHub Actions)

- [x] 2.1 创建 CI 工作流配置 (.github/workflows/ci.yml)
- [x] 2.2 配置后端测试任务 (go test)
- [x] 2.3 配置前端检查任务 (npm run lint, type-check)
- [x] 2.4 配置后端构建任务 (go build, Docker 镜像)
- [x] 2.5 配置前端构建任务 (npm run build)
- [x] 2.6 配置 Windows 客户端构建任务 (wails build, 仅 tag 触发)
- [x] 2.7 创建 CD 部署工作流 (.github/workflows/deploy.yml)
- [x] 2.8 配置 SSH 部署到服务器 (appleboy/ssh-action)
- [x] 2.9 在 GitHub 仓库配置 Secrets (SERVER_HOST, SERVER_USER, SSH_KEY)
- [ ] 2.10 测试 CI 流程（提交代码验证自动测试）
- [ ] 2.11 测试 CD 流程（验证自动部署到本机）
- [ ] 2.12 测试客户端构建流程（创建 tag 验证 .exe 生成）

## 3. 服务器环境准备

- [ ] 3.1 安装 Docker 和 Docker Compose
- [ ] 3.2 配置本机防火墙开放端口 (443, 5432, 8080, 21115, 21117, 21119)
- [ ] 3.3 配置 SSH 服务（允许 GitHub Actions 部署）
- [ ] 3.4 生成 SSH 密钥对（用于 CI/CD 部署）
- [ ] 3.5 创建项目目录 /opt/store-ops
- [ ] 3.6 配置 Nginx 反向代理和 HTTPS 证书

## 4. 数据库设计与迁移

- [ ] 4.1 创建数据库迁移脚本 (migrations/001_init.up.sql)
- [ ] 4.2 创建回滚脚本 (migrations/001_init.down.sql)
- [ ] 4.3 创建数据库连接模块 (internal/pkg/database/)
- [ ] 4.4 初始化默认管理员账户

## 5. 服务端 - 基础架构

- [ ] 5.1 实现配置管理模块 (internal/config/)
- [ ] 5.2 实现日志模块 (pkg/logger/)
- [ ] 5.3 实现统一响应格式 (internal/pkg/response/)
- [ ] 5.4 实现加密工具 (internal/pkg/crypto/)
- [ ] 5.5 实现 JWT 认证中间件 (internal/middleware/auth.go)
- [ ] 5.6 实现 CORS 中间件 (internal/middleware/cors.go)
- [ ] 5.7 实现请求日志中间件 (internal/middleware/logger.go)

## 6. 服务端 - 数据模型与仓库

- [ ] 6.1 创建 Admin 模型和仓库 (internal/model/admin.go, internal/repository/admin_repo.go)
- [ ] 6.2 创建 Store 模型和仓库 (internal/model/store.go, internal/repository/store_repo.go)
- [ ] 6.3 创建 AuditLog 模型和仓库 (internal/model/audit_log.go, internal/repository/audit_log_repo.go)
- [ ] 6.4 创建 CommandLog 模型和仓库 (internal/model/command_log.go, internal/repository/command_log_repo.go)
- [ ] 6.5 创建 SoftwareTask 模型和仓库 (internal/model/software_task.go, internal/repository/software_task_repo.go)
- [ ] 6.6 创建 DesktopSession 模型和仓库 (internal/model/desktop_session.go, internal/repository/desktop_session_repo.go)
- [ ] 6.7 创建 ResourceHistory 模型和仓库 (internal/model/resource_history.go, internal/repository/resource_history_repo.go)
- [ ] 6.8 创建 DesktopQueue 模型和仓库

## 7. 服务端 - WebSocket 网关

- [ ] 7.1 实现 WebSocket Hub 连接管理器 (internal/ws/hub.go)
- [ ] 7.2 实现 WebSocket 客户端连接 (internal/ws/client.go)
- [ ] 7.3 定义 WebSocket 消息类型 (internal/ws/message.go)
- [ ] 7.4 实现 WebSocket 消息处理器 (internal/ws/handler.go)
- [ ] 7.5 实现 TOKEN 认证逻辑
- [ ] 7.6 实现设备指纹校验
- [ ] 7.7 实现心跳检测机制

## 8. 服务端 - 业务服务

- [ ] 8.1 实现认证服务 (internal/service/auth_service.go)
- [ ] 8.2 实现门店管理服务 (internal/service/store_service.go)
- [ ] 8.3 实现远程桌面服务 (internal/service/desktop_service.go)
- [ ] 8.4 实现命令执行服务 (internal/service/command_service.go)
- [ ] 8.5 实现软件推送服务 (internal/service/software_service.go)
- [ ] 8.6 实现进程监控服务 (internal/service/process_service.go)
- [ ] 8.7 实现资源监控服务 (internal/service/resource_service.go)
- [ ] 8.8 实现审计日志服务 (internal/service/audit_service.go)

## 9. 服务端 - HTTP 处理器

- [ ] 9.1 实现认证处理器 (internal/handler/auth.go) - POST /auth/login, POST /auth/logout
- [ ] 9.2 实现门店管理处理器 (internal/handler/store.go) - GET/POST/DELETE /stores
- [ ] 9.3 实现远程桌面处理器 (internal/handler/desktop.go) - 连接状态、排队
- [ ] 9.4 实现命令执行处理器 (internal/handler/command.go) - 执行命令、获取结果
- [ ] 9.5 实现软件推送处理器 (internal/handler/software.go) - 推送安装、查询状态
- [ ] 9.6 实现进程监控处理器 (internal/handler/process.go) - 进程列表、终止进程、服务管理
- [ ] 9.7 实现资源监控处理器 (internal/handler/resource.go) - 当前资源、历史数据
- [ ] 9.8 实现审计日志处理器 (internal/handler/audit.go) - 日志列表、详情

## 10. 服务端 - 路由与启动

- [ ] 10.1 配置 HTTP 路由 (cmd/server/main.go)
- [ ] 10.2 配置 WebSocket 路由
- [ ] 10.3 实现优雅关闭
- [ ] 10.4 集成所有模块并测试启动

## 11. 服务端 - 单元测试

- [ ] 11.1 创建测试工具函数 (internal/testutil/setup.go)
- [ ] 11.2 编写密码加密/比对测试 (internal/pkg/crypto/crypto_test.go)
- [ ] 11.3 编写 TOKEN 生成/验证测试 (internal/pkg/auth/auth_test.go)
- [ ] 11.4 编写设备指纹生成测试 (internal/pkg/device/device_test.go)
- [ ] 11.5 编写 WebSocket 消息解析测试 (internal/ws/message_test.go)
- [ ] 11.6 编写统一响应格式测试 (internal/pkg/response/response_test.go)

## 12. 服务端 - 集成测试

- [ ] 12.1 编写认证 API 测试 - 登录成功/失败场景 (internal/handler/auth_test.go)
- [ ] 12.2 编写门店管理 API 测试 - 列表/详情/审批 (internal/handler/store_test.go)
- [ ] 12.3 编写命令执行 API 测试 - 执行/结果/离线场景 (internal/handler/command_test.go)
- [ ] 12.4 编写进程监控 API 测试 (internal/handler/process_test.go)
- [ ] 12.5 编写资源监控 API 测试 (internal/handler/resource_test.go)
- [ ] 12.6 编写审计日志 API 测试 (internal/handler/audit_test.go)
- [ ] 12.7 配置 CI 测试覆盖率报告

## 13. Web 前端 - 基础架构

- [ ] 13.1 配置路由 (src/router/index.ts)
- [ ] 13.2 配置 Pinia 状态管理 (src/stores/)
- [ ] 13.3 配置 Axios 请求封装 (src/utils/request.ts)
- [ ] 13.4 实现布局组件 (src/components/Layout/)
- [ ] 13.5 配置 Element Plus 主题
- [ ] 13.6 实现登录页面 (src/views/login/)

## 14. Web 前端 - API 封装

- [ ] 14.1 实现认证 API (src/api/auth.ts)
- [ ] 14.2 实现门店管理 API (src/api/store.ts)
- [ ] 14.3 实现远程桌面 API (src/api/desktop.ts)
- [ ] 14.4 实现命令执行 API (src/api/command.ts)
- [ ] 14.5 实现软件推送 API (src/api/software.ts)
- [ ] 14.6 实现进程监控 API (src/api/process.ts)
- [ ] 14.7 实现资源监控 API (src/api/resource.ts)
- [ ] 14.8 实现审计日志 API (src/api/audit.ts)

## 15. Web 前端 - 页面开发

- [ ] 15.1 实现仪表盘页面 (src/views/dashboard/) - 概览统计
- [ ] 15.2 实现门店列表页面 (src/views/stores/) - 列表、审批、状态
- [ ] 15.3 实现远程桌面页面 (src/views/desktop/) - 连接、排队
- [ ] 15.4 实现命令执行页面 (src/views/command/) - 执行、结果展示
- [ ] 15.5 实现软件推送页面 (src/views/software/) - 上传、进度
- [ ] 15.6 实现进程监控页面 (src/views/process/) - 进程列表、服务列表
- [ ] 15.7 实现资源监控页面 (src/views/resource/) - 实时数据、图表
- [ ] 15.8 实现审计日志页面 (src/views/audit/) - 日志列表、筛选

## 16. Web 前端 - 通用组件

- [ ] 16.1 实现门店状态组件 (src/components/StoreStatus/)
- [ ] 16.2 实现命令对话框组件 (src/components/CommandDialog/)
- [ ] 16.3 实现资源图表组件 (src/components/ResourceChart/)

## 17. 门店客户端 - 基础架构

- [ ] 17.1 初始化 Wails 项目结构
- [ ] 17.2 实现配置管理 (internal/config/)
- [ ] 17.3 实现 SQLite 本地存储 (internal/storage/sqlite.go)
- [ ] 17.4 实现设备 ID 生成 (internal/device/id.go)
- [ ] 17.5 实现设备指纹采集 (internal/device/fingerprint.go)
- [ ] 17.6 实现开机自启功能 (internal/autostart/windows.go)
- [ ] 17.7 实现崩溃自动重启

## 18. 门店客户端 - WebSocket 连接

- [ ] 18.1 实现 WebSocket 客户端 (internal/connection/client.go)
- [ ] 18.2 实现断线重连机制 (internal/connection/reconnect.go)
- [ ] 18.3 实现心跳发送 (internal/connection/heartbeat.go)
- [ ] 18.4 实现消息收发处理
- [ ] 18.5 实现注册/认证流程

## 19. 门店客户端 - 功能模块

- [ ] 19.1 实现命令执行模块 (internal/executor/command.go)
- [ ] 19.2 实现脚本执行模块 (internal/executor/script.go)
- [ ] 19.3 实现软件下载模块 (internal/software/downloader.go)
- [ ] 19.4 实现软件安装模块 (internal/software/installer.go)
- [ ] 19.5 实现进程列表获取 (internal/process/list.go)
- [ ] 19.6 实现进程终止功能 (internal/process/kill.go)
- [ ] 19.7 实现服务列表获取 (internal/service/list.go)
- [ ] 19.8 实现服务控制功能 (internal/service/control.go)
- [ ] 19.9 实现 CPU 数据采集 (internal/resource/cpu.go)
- [ ] 19.10 实现内存数据采集 (internal/resource/memory.go)
- [ ] 19.11 实现磁盘数据采集 (internal/resource/disk.go)
- [ ] 19.12 实现资源定时上报

## 20. 门店客户端 - RustDesk 集成

- [ ] 20.1 下载并集成 RustDesk 服务端组件
- [ ] 20.2 实现 RustDesk 配置生成 (internal/rustdesk/config.go)
- [ ] 20.3 实现 RustDesk 进程管理 (internal/rustdesk/embed.go)
- [ ] 20.4 实现自动获取 RustDesk ID
- [ ] 20.5 配置无密码连接

## 21. 门店客户端 - 界面开发

- [ ] 21.1 实现状态显示界面 (frontend/src/components/StatusDisplay.vue)
- [ ] 21.2 实现门店名称输入对话框（首次启动）
- [ ] 21.3 实现系统托盘图标
- [ ] 21.4 实现右键菜单（退出、查看状态）

## 22. RustDesk Server 部署

- [ ] 22.1 配置 hbbs 服务
- [ ] 22.2 配置 hbbr 服务
- [ ] 22.3 生成并配置密钥
- [ ] 22.4 测试 ID 注册和中继功能

## 23. 系统集成测试

- [ ] 23.1 测试门店客户端连接服务端
- [ ] 23.2 测试门店审批流程
- [ ] 23.3 测试远程桌面连接（直连模式）
- [ ] 23.4 测试远程桌面连接（中继模式）
- [ ] 23.5 测试命令执行功能
- [ ] 23.6 测试软件推送功能
- [ ] 23.7 测试进程监控功能
- [ ] 23.8 测试资源监控功能
- [ ] 23.9 测试审计日志记录
- [ ] 23.10 测试断线重连
- [ ] 23.11 测试多用户排队

## 24. 部署与文档

- [ ] 24.1 编写部署文档
- [ ] 24.2 编写用户手册（门店客户端安装）
- [ ] 24.3 编写 IT 操作手册（Web 平台使用）
- [ ] 24.4 打包门店客户端安装程序
- [ ] 24.5 部署服务端到生产环境
- [ ] 24.6 部署 Web 前端到生产环境
- [ ] 24.7 配置 Nginx 反向代理
- [ ] 24.8 配置 HTTPS 证书