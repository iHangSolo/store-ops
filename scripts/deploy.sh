#!/bin/bash

# 门店运维系统部署脚本
# 用法: ./scripts/deploy.sh [environment]
# environment: production | staging

set -e

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
ENVIRONMENT=${1:-production}

echo "=========================================="
echo "  门店运维系统部署脚本"
echo "  环境: $ENVIRONMENT"
echo "=========================================="

cd "$PROJECT_ROOT"

# 检查 Docker 和 Docker Compose
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    exit 1
fi

# 使用 docker compose 或 docker-compose
COMPOSE_CMD="docker compose"
if ! docker compose version &> /dev/null; then
    COMPOSE_CMD="docker-compose"
fi

echo ""
echo "[1/4] 拉取最新镜像..."
$COMPOSE_CMD pull

echo ""
echo "[2/4] 构建镜像..."
$COMPOSE_CMD build

echo ""
echo "[3/4] 停止旧容器..."
$COMPOSE_CMD down

echo ""
echo "[4/4] 启动新容器..."
$COMPOSE_CMD up -d

echo ""
echo "等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "服务状态:"
$COMPOSE_CMD ps

echo ""
echo "=========================================="
echo "  部署完成!"
echo "=========================================="
echo ""
echo "访问地址:"
echo "  - Web 管理平台: http://localhost"
echo "  - API 服务:     http://localhost:8080"
echo "  - RustDesk:     localhost:21115"
echo ""
echo "查看日志:"
echo "  $COMPOSE_CMD logs -f"
echo ""