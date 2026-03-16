# 数据库设计

## ER 图

```
┌─────────────────┐       ┌─────────────────┐
│     admins      │       │     stores      │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │
│ username        │       │ name            │
│ password_hash   │       │ device_id       │
│ created_at      │       │ device_fingerprint│
│ updated_at      │       │ rustdesk_id     │
└─────────────────┘       │ status          │
                          │ approved        │
                          │ last_seen       │
                          │ created_at      │
                          │ updated_at      │
                          └────────┬────────┘
                                   │
                                   │ 1:N
                                   │
                          ┌────────▼────────┐
                          │  audit_logs     │
                          ├─────────────────┤
                          │ id (PK)         │
                          │ store_id (FK)   │
                          │ action_type     │
                          │ operator        │
                          │ details (JSONB) │
                          │ created_at      │
                          └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│ resource_history│       │ desktop_sessions│
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │
│ store_id (FK)   │       │ store_id (FK)   │
│ cpu_percent     │       │ operator        │
│ memory_percent  │       │ started_at      │
│ memory_available│       │ ended_at        │
│ disk_data (JSONB)│       │ duration_seconds│
│ recorded_at     │       │ status          │
└─────────────────┘       └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│ command_logs    │       │ software_tasks  │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │
│ store_id (FK)   │       │ store_id (FK)   │
│ operator        │       │ operator        │
│ command         │       │ file_name       │
│ output          │       │ install_args    │
│ status          │       │ status          │
│ executed_at     │       │ created_at      │
│ duration_ms     │       │ completed_at    │
└─────────────────┘       └─────────────────┘

┌─────────────────┐
│ desktop_queue   │
├─────────────────┤
│ id (PK)         │
│ store_id (FK)   │
│ operator        │
│ position        │
│ created_at      │
└─────────────────┘
```

---

## 表结构详情

### admins（管理员表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| username | VARCHAR(50) | UNIQUE, NOT NULL | 用户名 |
| password_hash | VARCHAR(255) | NOT NULL | 密码哈希（bcrypt） |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |

```sql
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### stores（门店表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| name | VARCHAR(100) | NOT NULL | 门店名称 |
| device_id | UUID | UNIQUE, NOT NULL | 设备唯一标识 |
| device_fingerprint | VARCHAR(64) | NOT NULL | 设备指纹（SHA256） |
| rustdesk_id | VARCHAR(20) | | RustDesk ID |
| status | VARCHAR(20) | NOT NULL | 状态：online/offline |
| approved | BOOLEAN | NOT NULL DEFAULT false | 是否已审批 |
| last_seen | TIMESTAMP | | 最后在线时间 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |

```sql
CREATE TABLE stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    device_id UUID UNIQUE NOT NULL,
    device_fingerprint VARCHAR(64) NOT NULL,
    rustdesk_id VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'offline',
    approved BOOLEAN NOT NULL DEFAULT false,
    last_seen TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stores_status ON stores(status);
CREATE INDEX idx_stores_approved ON stores(approved);
```

### audit_logs（审计日志表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| store_id | UUID | FK | 关联门店 |
| action_type | VARCHAR(50) | NOT NULL | 操作类型 |
| operator | VARCHAR(50) | NOT NULL | 操作人 |
| details | JSONB | | 操作详情 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID REFERENCES stores(id) ON DELETE SET NULL,
    action_type VARCHAR(50) NOT NULL,
    operator VARCHAR(50) NOT NULL,
    details JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_store_id ON audit_logs(store_id);
CREATE INDEX idx_audit_logs_action_type ON audit_logs(action_type);
```

**action_type 枚举值：**
- `desktop_connect` - 远程桌面连接
- `desktop_disconnect` - 远程桌面断开
- `command_execute` - 命令执行
- `software_push` - 软件推送
- `process_kill` - 终止进程
- `service_control` - 服务控制

### resource_history（资源历史表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| store_id | UUID | FK | 关联门店 |
| cpu_percent | DECIMAL(5,2) | NOT NULL | CPU 使用率 |
| memory_percent | DECIMAL(5,2) | NOT NULL | 内存使用率 |
| memory_available_gb | DECIMAL(10,2) | NOT NULL | 可用内存 |
| disk_data | JSONB | | 磁盘数据 |
| recorded_at | TIMESTAMP | NOT NULL | 记录时间 |

```sql
CREATE TABLE resource_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    cpu_percent DECIMAL(5,2) NOT NULL,
    memory_percent DECIMAL(5,2) NOT NULL,
    memory_available_gb DECIMAL(10,2) NOT NULL,
    disk_data JSONB,
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_resource_history_store_time ON resource_history(store_id, recorded_at);
```

**disk_data JSONB 结构：**
```json
[
  {"name": "C:", "percent": 45.0, "available_gb": 256.5}
]
```

### desktop_sessions（远程桌面会话表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| store_id | UUID | FK | 关联门店 |
| operator | VARCHAR(50) | NOT NULL | 操作人 |
| started_at | TIMESTAMP | NOT NULL | 开始时间 |
| ended_at | TIMESTAMP | | 结束时间 |
| duration_seconds | INT | | 持续时长 |
| status | VARCHAR(20) | NOT NULL | 状态 |

```sql
CREATE TABLE desktop_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    operator VARCHAR(50) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_seconds INT,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);
```

### command_logs（命令日志表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| store_id | UUID | FK | 关联门店 |
| operator | VARCHAR(50) | NOT NULL | 操作人 |
| command | TEXT | NOT NULL | 命令内容 |
| output | TEXT | | 执行输出 |
| status | VARCHAR(20) | NOT NULL | 执行状态 |
| executed_at | TIMESTAMP | NOT NULL | 执行时间 |
| duration_ms | INT | | 执行耗时 |

```sql
CREATE TABLE command_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID REFERENCES stores(id) ON DELETE SET NULL,
    operator VARCHAR(50) NOT NULL,
    command TEXT NOT NULL,
    output TEXT,
    status VARCHAR(20) NOT NULL,
    executed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration_ms INT
);

CREATE INDEX idx_command_logs_store_id ON command_logs(store_id);
```

### software_tasks（软件任务表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| store_id | UUID | FK | 关联门店 |
| operator | VARCHAR(50) | NOT NULL | 操作人 |
| file_name | VARCHAR(255) | NOT NULL | 文件名 |
| file_size | BIGINT | | 文件大小 |
| install_args | VARCHAR(500) | | 安装参数 |
| status | VARCHAR(20) | NOT NULL | 状态 |
| message | TEXT | | 状态消息 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| completed_at | TIMESTAMP | | 完成时间 |

```sql
CREATE TABLE software_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    operator VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    install_args VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP
);
```

### desktop_queue（远程桌面排队表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | UUID | PK | 主键 |
| store_id | UUID | FK | 关联门店 |
| operator | VARCHAR(50) | NOT NULL | 操作人 |
| position | INT | NOT NULL | 队列位置 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |

```sql
CREATE TABLE desktop_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    operator VARCHAR(50) NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_desktop_queue_store ON desktop_queue(store_id, position);
```

---

## 数据清理策略

### 审计日志（3 个月）

```sql
-- 定时任务：每天凌晨执行
DELETE FROM audit_logs
WHERE created_at < NOW() - INTERVAL '3 months';
```

### 资源历史（1 小时详情，聚合保留）

```sql
-- 保留最近 1 小时的详细数据
DELETE FROM resource_history
WHERE recorded_at < NOW() - INTERVAL '1 hour';
```

---

## 初始化脚本

```sql
-- 创建默认管理员
INSERT INTO admins (username, password_hash)
VALUES ('admin', '$2a$10$...');  -- 密码: admin123
```