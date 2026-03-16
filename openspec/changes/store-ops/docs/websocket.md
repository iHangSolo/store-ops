# WebSocket 消息协议

## 连接说明

### 连接地址

```
wss://<server>/api/v1/ws
```

### 认证

连接时通过查询参数携带 TOKEN：

```
wss://<server>/api/v1/ws?token=<jwt_token>
```

### 心跳

- 客户端每 30 秒发送一次心跳
- 服务端收到心跳后回复
- 3 次心跳无响应判定断线

---

## 消息格式

所有消息使用 JSON 格式：

```json
{
  "type": "<message_type>",
  "payload": { ... },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## 客户端 → 服务端

### 心跳

```json
{
  "type": "ping"
}
```

**响应：**
```json
{
  "type": "pong",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 注册/认证

首次连接或重连后发送：

```json
{
  "type": "register",
  "payload": {
    "device_id": "uuid-xxx",
    "device_fingerprint": "hash-xxx",
    "store_name": "上海南京东路店",
    "rustdesk_id": "123456789",
    "client_version": "1.0.0"
  }
}
```

**响应：**
```json
{
  "type": "register_ack",
  "payload": {
    "success": true,
    "approved": true,
    "message": "注册成功"
  }
}
```

### 资源上报

每 30 秒上报一次：

```json
{
  "type": "resource_report",
  "payload": {
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

---

## 服务端 → 客户端

### 执行命令

```json
{
  "type": "execute_command",
  "payload": {
    "command_id": "cmd-uuid",
    "command": "ipconfig /all",
    "timeout": 60
  }
}
```

**响应：**
```json
{
  "type": "command_result",
  "payload": {
    "command_id": "cmd-uuid",
    "status": "success",
    "output": "Windows IP Configuration...",
    "exit_code": 0,
    "duration_ms": 150
  }
}
```

### 执行脚本

```json
{
  "type": "execute_script",
  "payload": {
    "command_id": "cmd-uuid",
    "script_type": "bat",
    "script_content": "@echo off\necho Hello",
    "timeout": 60
  }
}
```

### 推送软件

```json
{
  "type": "software_push",
  "payload": {
    "task_id": "task-uuid",
    "download_url": "https://xxx/software.exe",
    "file_name": "software.exe",
    "file_size": 50000000,
    "install_args": "/S"
  }
}
```

**响应（进度）：**
```json
{
  "type": "software_progress",
  "payload": {
    "task_id": "task-uuid",
    "status": "downloading",
    "progress": 50
  }
}
```

**响应（完成）：**
```json
{
  "type": "software_result",
  "payload": {
    "task_id": "task-uuid",
    "status": "installed",
    "message": "安装成功"
  }
}
```

### 查询进程

```json
{
  "type": "get_processes",
  "payload": {
    "request_id": "req-uuid"
  }
}
```

**响应：**
```json
{
  "type": "processes",
  "payload": {
    "request_id": "req-uuid",
    "processes": [
      {
        "name": "chrome.exe",
        "pid": 1234,
        "cpu_percent": 5.2,
        "memory_percent": 3.5
      }
    ]
  }
}
```

### 终止进程

```json
{
  "type": "kill_process",
  "payload": {
    "request_id": "req-uuid",
    "pid": 1234
  }
}
```

**响应：**
```json
{
  "type": "kill_process_result",
  "payload": {
    "request_id": "req-uuid",
    "success": true,
    "message": "进程已终止"
  }
}
```

### 查询服务

```json
{
  "type": "get_services",
  "payload": {
    "request_id": "req-uuid"
  }
}
```

**响应：**
```json
{
  "type": "services",
  "payload": {
    "request_id": "req-uuid",
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

### 控制服务

```json
{
  "type": "control_service",
  "payload": {
    "request_id": "req-uuid",
    "service_name": "Winmgmt",
    "action": "start"
  }
}
```

action 可选值：`start`, `stop`, `restart`

**响应：**
```json
{
  "type": "control_service_result",
  "payload": {
    "request_id": "req-uuid",
    "success": true,
    "message": "服务已启动"
  }
}
```

---

## 错误消息

```json
{
  "type": "error",
  "payload": {
    "code": "ERROR_CODE",
    "message": "错误描述"
  }
}
```

常见错误码：
| code | 说明 |
|------|------|
| AUTH_FAILED | 认证失败 |
| DEVICE_NOT_APPROVED | 设备未审批 |
| COMMAND_TIMEOUT | 命令执行超时 |
| COMMAND_FAILED | 命令执行失败 |
| PROCESS_NOT_FOUND | 进程不存在 |
| SERVICE_NOT_FOUND | 服务不存在 |