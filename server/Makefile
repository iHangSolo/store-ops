.PHONY: build run test clean migrate

# 默认目标
all: build

# 构建
build:
	go build -o bin/server ./cmd/server

# 运行
run:
	go run ./cmd/server

# 测试
test:
	go test -v ./...

# 测试覆盖率
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 数据库迁移
migrate-up:
	@echo "Run migrations..."

migrate-down:
	@echo "Rollback migrations..."

# 开发模式（热重载需要安装 air）
dev:
	which air > /dev/null || go install github.com/cosmtrek/air@latest
	air

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# 安装依赖
deps:
	go mod download
	go mod tidy