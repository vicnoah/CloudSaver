#!/bin/bash

echo "Building CloudSaver..."

# 检查依赖
command -v go >/dev/null 2>&1 || { echo "需要安装 Go"; exit 1; }
command -v pnpm >/dev/null 2>&1 || { echo "需要安装 pnpm"; exit 1; }

# 构建前端
echo "Step 1: Building frontend..."
cd web
pnpm install
pnpm build
cd ..

# 复制前端产物到embed目录
echo "Step 2: Copying frontend assets to embed..."
mkdir -p embed/dist
cp -r web/dist/* embed/dist/

# 构建Go程序
echo "Step 3: Building Go binary..."
go mod download
go build -ldflags="-s -w" -o bin/cloudsaver cmd/server/main.go

echo "Build completed successfully!"
echo "Run with: ./bin/cloudsaver"
