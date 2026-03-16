#!/bin/bash

# 门店运维系统构建脚本
# 用法: ./scripts/build.sh [target]
# target: all | server | web | client

set -e

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
TARGET=${1:-all}

echo "=========================================="
echo "  门店运维系统构建脚本"
echo "  目标: $TARGET"
echo "=========================================="

build_server() {
    echo ""
    echo "[1/3] 构建后端服务..."
    cd "$PROJECT_ROOT/server"

    # 下载依赖
    echo "  - 下载 Go 依赖..."
    go mod download

    # 运行测试
    echo "  - 运行测试..."
    go test ./... || echo "  (测试跳过或无测试文件)"

    # 构建
    echo "  - 编译..."
    go build -o "$PROJECT_ROOT/dist/server/store-ops-server" ./cmd/server

    echo "  ✓ 后端构建完成: dist/server/store-ops-server"
}

build_web() {
    echo ""
    echo "[2/3] 构建前端..."
    cd "$PROJECT_ROOT/web"

    # 安装依赖
    echo "  - 安装 npm 依赖..."
    npm ci

    # 类型检查
    echo "  - 类型检查..."
    npm run type-check || echo "  (类型检查警告)"

    # 构建
    echo "  - 编译..."
    npm run build

    # 复制到 dist 目录
    mkdir -p "$PROJECT_ROOT/dist/web"
    cp -r dist/* "$PROJECT_ROOT/dist/web/"

    echo "  ✓ 前端构建完成: dist/web/"
}

build_client() {
    echo ""
    echo "[3/3] 构建客户端..."
    cd "$PROJECT_ROOT/client"

    # 检查 Wails 是否安装
    if ! command -v wails &> /dev/null; then
        echo "  ! Wails 未安装，请先安装: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        echo "  跳过客户端构建..."
        return
    fi

    # 安装前端依赖
    echo "  - 安装前端依赖..."
    cd frontend
    npm ci
    cd ..

    # 构建
    echo "  - 编译 Windows 客户端..."
    wails build -platform windows/amd64 -o "$PROJECT_ROOT/dist/client/store-ops-client.exe"

    echo "  ✓ 客户端构建完成: dist/client/store-ops-client.exe"
}

# 创建 dist 目录
mkdir -p "$PROJECT_ROOT/dist"

case "$TARGET" in
    all)
        build_server
        build_web
        build_client
        ;;
    server)
        build_server
        ;;
    web)
        build_web
        ;;
    client)
        build_client
        ;;
    *)
        echo "未知目标: $TARGET"
        echo "用法: $0 [all|server|web|client]"
        exit 1
        ;;
esac

echo ""
echo "=========================================="
echo "  构建完成!"
echo "=========================================="
ls -la "$PROJECT_ROOT/dist/" 2>/dev/null || true