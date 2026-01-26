#!/bin/bash

# Task Manager - 启动脚本

echo "🚀 启动 Task Manager 全栈应用..."
echo ""

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go"
    exit 1
fi

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js"
    exit 1
fi

# 启动后端
echo "📦 启动后端服务器..."
cd backend
go run main.go &
BACKEND_PID=$!
echo "✅ 后端运行在 http://localhost:8080 (PID: $BACKEND_PID)"

# 等待后端启动
sleep 3

# 启动前端
echo ""
echo "🎨 启动前端开发服务器..."
cd ../frontend

# 检查是否安装了依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    if command -v pnpm &> /dev/null; then
        pnpm install
    else
        npm install
    fi
fi

if command -v pnpm &> /dev/null; then
    pnpm dev
else
    npm run dev
fi

# 清理后台进程
trap "echo ''; echo '🛑 停止服务...'; kill $BACKEND_PID 2>/dev/null; exit" INT TERM
