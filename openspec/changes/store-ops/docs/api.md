# REST API 设计

## 通用说明

### 基础路径

```
https://<server>/api/v1
```

### 认证方式

- 管理员 API：JWT Token（Header: `Authorization: Bearer <token>`）
- 所有响应使用 JSON 格式

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误码

| code | 含义 |
|------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 认证失败 |
| 1003 | 权限不足 |
| 2001 | 门店离线 |
| 2002 | 命令执行超时 |
| 5000 | 服务器内部错误 |

---

## 认证 API

### POST /auth/login

管理员登录

**请求：**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2024-01-02T00:00:00Z"
  }
}
```

### POST /auth/logout

登出

---

## 门店管理 API

### GET /stores

获取门店列表

**查询参数：**
| 参数 | 类型 | 说明 |
|------|------|------|
| status | string | 筛选状态：online/offline/pending |
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应：**
```json
{
  "code": 0,
  "data": {
    "total": 80,
    "items": [
      {
        "id": "uuid-xxx",
        "name": "上海南京东路店",
        "device_id": "device-uuid",
        "status": "online",
        "last_seen": "2024-01-01T12:00:00Z",
        "approved": true,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### GET /stores/{id}

获取门店详情

**响应：**
```json
{
  "code": 0,
  "data": {
    "id": "uuid-xxx",
    "name": "上海南京东路店",
    "device_id": "device-uuid",
    "device_fingerprint": "hash-xxx",
    "status": "online",
    "last_seen": "2024-01-01T12:00:00Z",
    "approved": true,
    "rustdesk_id": "123456789",
    "created_at": "2024-01-01T00:00:00Z",
    "resource": {
      "cpu_percent": 35.5,
      "memory_percent": 62.3,
      "memory_available_gb": 3.2,
      "disk_percent": 45.0,
      "disk_available_gb": 256.5
    }
  }
}
```

### POST /stores/{id}/approve

审批门店

**响应：**
```json
{
  "code": 0,
  "message": "门店已审批通过"
}
```

### POST /stores/{id}/reject

拒绝门店

**响应：**
```json
{
  "code": 0,
  "message": "门店已拒绝"
}
```

### DELETE /stores/{id}

删除门店（解绑）

---

## 远程桌面 API

### GET /stores/{id}/desktop/status

获取远程桌面状态

**响应：**
```json
{
  "code": 0,
  "data": {
    "rustdesk_id": "123456789",
    "is_busy": false,
    "current_user": null,
    "queue_count": 0
  }
}
```

### POST /stores/{id}/desktop/connect

请求远程桌面连接

**响应：**
```json
{
  "code": 0,
  "data": {
    "rustdesk_id": "123456789",
    "rustdesk_url": "rustdesk://connect?id=123456789&server=xxx"
  }
}
```

### POST /stores/{id}/desktop/queue

加入排队

**响应：**
```json
{
  "code": 0,
  "data": {
    "queue_position": 2
  }
}
```

---

## 命令执行 API

### POST /stores/{id}/command

执行命令

**请求：**
```json
{
  "command": "ipconfig /all",
  "timeout": 60
}
```

或执行脚本：

```json
{
  "script_type": "bat",
  "script_content": "@echo off\necho Hello",
  "timeout": 60
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "command_id": "cmd-uuid",
    "status": "executing"
  }
}
```

### GET /stores/{id}/command/{command_id}

获取命令执行结果

**响应：**
```json
{
  "code": 0,
  "data": {
    "command_id": "cmd-uuid",
    "status": "success",
    "output": "Windows IP Configuration...",
    "executed_at": "2024-01-01T12:00:00Z",
    "duration_ms": 150
  }
}
```

---

## 软件推送 API

### POST /stores/{id}/software

推送软件安装

**请求：**
```
Content-Type: multipart/form-data

file: software.exe
install_args: /S
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "task_id": "task-uuid",
    "status": "uploading"
  }
}
```

### GET /stores/{id}/software/{task_id}

获取推送状态

**响应：**
```json
{
  "code": 0,
  "data": {
    "task_id": "task-uuid",
    "status": "installed",
    "progress": 100,
    "message": "安装成功"
  }
}
```

---

## 进程监控 API

### GET /stores/{id}/processes

获取进程列表

**响应：**
```json
{
  "code": 0,
  "data": {
    "processes": [
      {
        "name": "chrome.exe",
        "pid": 1234,
        "cpu_percent": 5.2,
        "memory_percent": 3.5,
        "status": "running"
      }
    ]
  }
}
```

### DELETE /stores/{id}/processes/{pid}

终止进程

**响应：**
```json
{
  "code": 0,
  "message": "进程已终止"
}
```

### GET /stores/{id}/services

获取服务列表

**响应：**
```json
{
  "code": 0,
  "data": {
    "services": [
      {
        "name": "Winmgmt",
        "display_name": "Windows Management Instrumentation",
        "status": "running"
      }
    ]
  }
}
```

### POST /stores/{id}/services/{name}/start

启动服务

### POST /stores/{id}/services/{name}/stop

停止服务

---

## 资源监控 API

### GET /stores/{id}/resource

获取当前资源状态

**响应：**
```json
{
  "code": 0,
  "data": {
    "cpu_percent": 35.5,
    "memory_percent": 62.3,
    "memory_total_gb": 8.0,
    "memory_available_gb": 3.2,
    "disks": [
      {
        "name": "C:",
        "percent": 45.0,
        "total_gb": 500.0,
        "available_gb": 256.5
      }
    ]
  }
}
```

### GET /stores/{id}/resource/history

获取资源历史数据

**查询参数：**
| 参数 | 类型 | 说明 |
|------|------|------|
| duration | string | 时间范围：1h/6h/24h，默认 1h |

**响应：**
```json
{
  "code": 0,
  "data": {
    "cpu": [
      {"time": "2024-01-01T11:00:00Z", "value": 35.5},
      {"time": "2024-01-01T11:01:00Z", "value": 42.3}
    ],
    "memory": [...]
  }
}
```

---

## 审计日志 API

### GET /audit-logs

获取审计日志列表

**查询参数：**
| 参数 | 类型 | 说明 |
|------|------|------|
| start_date | string | 开始日期 |
| end_date | string | 结束日期 |
| action_type | string | 操作类型：desktop/command |
| store_id | string | 门店 ID |
| page | int | 页码 |
| page_size | int | 每页数量 |

**响应：**
```json
{
  "code": 0,
  "data": {
    "total": 150,
    "items": [
      {
        "id": "log-uuid",
        "action_type": "desktop_connect",
        "operator": "admin",
        "store_id": "store-uuid",
        "store_name": "上海南京东路店",
        "details": {
          "duration_seconds": 300
        },
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

### GET /audit-logs/{id}

获取审计日志详情