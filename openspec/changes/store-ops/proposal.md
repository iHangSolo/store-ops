## 为什么

服装公司需要远程管理全国 80 家直营门店的 Windows 电脑，目前缺乏统一的远程运维工具。门店电脑仅用于收银和办公，需要总部署 IT 人员进行远程控制、资源监控和软件维护。现有方案无法满足批量管理、审计合规、低延迟远程桌面等需求。

## 变更内容

构建一套门店远程运维管理系统（store-ops），包含：

- **门店客户端**：安装在门店电脑上，支持远程桌面、资源监控、命令执行
- **总部服务端**：提供管理 API、WebSocket 网关、RustDesk 服务器
- **Web 管理平台**：供总部 IT 人员进行远程运维操作
- **IT 客户端**：RustDesk 客户端，通过 Web 唤起进行远程桌面连接

## 功能

### 新增功能

- `remote-desktop`: 远程桌面控制，支持 Web 唤起本地 RustDesk 客户端直连门店电脑
- `remote-command`: 远程命令执行，支持在门店电脑上执行命令/脚本，需二次确认
- `software-deploy`: 软件推送安装/更新，支持向门店电脑推送软件包
- `process-monitor`: 进程/服务状态查看，实时获取门店电脑进程列表
- `resource-monitor`: 资源监控，实时采集门店电脑 CPU、内存、磁盘使用情况
- `audit-log`: 审计日志，记录远程桌面连接和命令执行操作，保留 3 个月
- `client-management`: 客户端管理，支持门店主机注册、审批、状态监控

### 修改功能

无（全新系统）

## 影响

### 新增系统组件

| 组件 | 技术栈 | 部署位置 |
|------|--------|---------|
| store-ops-server | Go + Gin + WebSocket | 总部 Linux 服务器 |
| store-ops-web | Vue 3 + Element Plus | 总部 Linux 服务器 |
| store-ops-client | Wails (Go + Vue) + SQLite + RustDesk | 门店 Windows 电脑 |
| RustDesk Server | hbbs + hbbr | 总部 Linux 服务器 |
| PostgreSQL | 15.x | 总部 Linux 服务器 |

### 网络端口

- **总部服务器入站**：443（HTTPS/WSS）、21115/21117/21119（RustDesk）
- **门店客户端出站**：443、21115/21117/21119

### 运维要求

- 门店客户端需开机自启
- 门店客户端需配置防火墙白名单（如被拦截）
- IT 人员需安装 RustDesk 客户端