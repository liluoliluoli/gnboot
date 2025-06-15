#!/bin/bash

# ================================================
# Docker 镜像自动构建部署脚本
# 功能：构建镜像 → 打标签 → 推送至仓库 → 启动服务 (未启用)
# 版本号格式：年.月.日.小时 (如 2025.06.16.14)
# ================================================

PROJECT_ROOT="../"
DOCKERFILE_PATH="$PROJECT_ROOT/Dockerfile"
COMPOSE_FILE="$PROJECT_ROOT/docker-compose/docker-compose.yml"
IMAGE_NAME="gnboot"
REPO_NAME="sily1/gnboot"

# 获取当前版本号
VERSION=$(date '+%Y.%m.%d.%H')
FULL_IMAGE_NAME="$REPO_NAME:$VERSION"

# 验证前置条件
function validate() {
    echo "验证环境配置..."

    # 检查 Docker 是否可用
    if ! command -v docker &> /dev/null; then
        echo "错误：未检测到 Docker，请先安装 Docker"
        exit 1
    fi

    # 检查 Dockerfile
    if [ ! -f "$DOCKERFILE_PATH" ]; then
        echo "错误：在以下路径未找到 Dockerfile: $DOCKERFILE_PATH"
        exit 1
    fi

    # 检查 docker-compose 文件
    if [ ! -f "$COMPOSE_FILE" ]; then
        echo "警告：未找到 docker-compose.yml 文件，将跳过服务启动步骤"
    fi

    # 检查是否已登录 Docker Hub
    if ! docker info | grep -q "Username: sily1"; then
        echo "提示：请先执行 docker login 登录 Docker Hub"
        # 这里不退出，因为可能使用其他仓库
    fi
}

# 构建 Docker 镜像
function build_image() {
    echo "Step 1: 构建 Docker 镜像 ($IMAGE_NAME)..."
    docker build -t "$IMAGE_NAME" -f "$DOCKERFILE_PATH" "$PROJECT_ROOT"

    if [ $? -ne 0 ]; then
        echo "错误：镜像构建失败！"
        exit 1
    fi
}

# 为镜像打标签
function tag_image() {
    echo "Step 2: 为镜像打标签 ($FULL_IMAGE_NAME)..."
    docker tag "$IMAGE_NAME:latest" "$FULL_IMAGE_NAME"

    if [ $? -ne 0 ]; then
        echo "错误：镜像标签操作失败！"
        exit 1
    fi
}

# 推送镜像到仓库
function push_image() {
    echo "Step 3: 推送镜像到 Docker Hub..."
    docker push "$FULL_IMAGE_NAME"

    if [ $? -ne 0 ]; then
        echo "错误：镜像推送失败！"
        exit 1
    fi
}

# 启动服务
function start_service() {
    if [ -f "$COMPOSE_FILE" ]; then
        echo "Step 4: 启动服务..."
        (cd "$(dirname "$COMPOSE_FILE")" && docker-compose up -d gnboot)
    else
        echo "跳过服务启动步骤（未找到 docker-compose.yml）"
    fi
}

# 主执行流程
echo "========== 开始执行部署脚本 =========="
validate
build_image
tag_image
push_image
#start_service

echo -e "\n========== 部署成功完成 =========="
echo "镜像版本: $FULL_IMAGE_NAME"
#echo "当前时间: $(date '+%Y-%m-%d %H:%M:%S')"