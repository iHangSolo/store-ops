# 测试策略

## 测试范围

MVP 阶段仅实施后端测试，前端和客户端测试后期迭代。

## 后端单元测试

### 测试框架

- Go 内置 `testing` 包
- 断言库：`testify/assert`（可选）

### 测试内容

| 模块 | 测试项 | 文件位置 |
|------|--------|---------|
| 认证 | TOKEN 生成、TOKEN 验证、TOKEN 过期 | `internal/pkg/auth/auth_test.go` |
| 加密 | 密码哈希、密码比对 | `internal/pkg/crypto/crypto_test.go` |
| 设备 | 设备 ID 生成、设备指纹生成 | `internal/pkg/device/device_test.go` |
| 响应 | 统一响应格式 | `internal/pkg/response/response_test.go` |
| WebSocket | 消息解析、消息序列化 | `internal/ws/message_test.go` |

### 测试示例

```go
// internal/pkg/crypto/crypto_test.go
package crypto

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPasswordHash(t *testing.T) {
    password := "admin123"

    hash, err := HashPassword(password)
    assert.NoError(t, err)
    assert.NotEmpty(t, hash)

    // 验证正确密码
    ok := CheckPassword(password, hash)
    assert.True(t, ok)

    // 验证错误密码
    ok = CheckPassword("wrongpassword", hash)
    assert.False(t, ok)
}

func TestTokenGeneration(t *testing.T) {
    userID := "user-123"

    token, err := GenerateToken(userID)
    assert.NoError(t, err)
    assert.NotEmpty(t, token)

    // 验证 TOKEN
    parsedUserID, err := ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, userID, parsedUserID)
}
```

### 运行命令

```bash
# 运行所有单元测试
go test ./...

# 运行指定包的测试
go test ./internal/pkg/crypto/...

# 显示详细输出
go test -v ./...
```

---

## 后端集成测试

### 测试框架

- Go 内置 `testing` 包
- 测试数据库：SQLite 内存模式 或 PostgreSQL Docker 容器

### 测试内容

| API | 测试场景 | 文件位置 |
|-----|---------|---------|
| POST /auth/login | 登录成功、用户不存在、密码错误 | `internal/handler/auth_test.go` |
| POST /auth/logout | 登出成功 | `internal/handler/auth_test.go` |
| GET /stores | 列表查询、分页、状态筛选 | `internal/handler/store_test.go` |
| GET /stores/:id | 详情查询、不存在返回 404 | `internal/handler/store_test.go` |
| POST /stores/:id/approve | 审批成功、门店不存在 | `internal/handler/store_test.go` |
| POST /stores/:id/command | 执行成功、门店离线 | `internal/handler/command_test.go` |
| GET /stores/:id/processes | 进程列表查询 | `internal/handler/process_test.go` |
| GET /stores/:id/resource | 资源数据查询 | `internal/handler/resource_test.go` |
| GET /audit-logs | 日志列表、筛选 | `internal/handler/audit_test.go` |

### 测试辅助函数

```go
// internal/testutil/setup.go
package testutil

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

// 创建测试服务器
func SetupTestServer(t *testing.T) *httptest.Server {
    // 初始化测试数据库
    // 初始化路由
    // 返回测试服务器
}

// 发送 HTTP 请求
func MakeRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
    var req *http.Request
    if body != nil {
        jsonBody, _ := json.Marshal(body)
        req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
        req.Header.Set("Content-Type", "application/json")
    } else {
        req = httptest.NewRequest(method, path, nil)
    }

    w := httptest.NewRecorder()
    // router.ServeHTTP(w, req)
    return w
}
```

### 测试示例

```go
// internal/handler/auth_test.go
package handler

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()

    // 创建测试用户
    createUser("admin", "password123")

    // 发送登录请求
    resp := MakeRequest("POST", "/api/v1/auth/login", map[string]string{
        "username": "admin",
        "password": "password123",
    })

    // 验证响应
    assert.Equal(t, 200, resp.Code)

    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)
    assert.Equal(t, float64(0), result["code"])
    assert.NotEmpty(t, result["data"].(map[string]interface{})["token"])
}

func TestLoginWrongPassword(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()

    createUser("admin", "password123")

    resp := MakeRequest("POST", "/api/v1/auth/login", map[string]string{
        "username": "admin",
        "password": "wrongpassword",
    })

    assert.Equal(t, 200, resp.Code)

    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)
    assert.Equal(t, float64(1002), result["code"]) // 认证失败
}
```

### 运行命令

```bash
# 运行集成测试（需要数据库）
go test -tags=integration ./...

# 运行所有测试
go test ./...
```

---

## CI 集成

GitHub Actions 自动运行测试：

```yaml
# .github/workflows/ci.yml
- name: Run tests
  working-directory: server
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    files: ./server/coverage.out
```

---

## 测试覆盖率目标

| 模块 | 目标覆盖率 |
|------|-----------|
| internal/pkg/ (工具包) | ≥ 80% |
| internal/handler/ (API) | ≥ 60% |
| internal/service/ (业务) | ≥ 50% |
| 整体 | ≥ 60% |

查看覆盖率：

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```