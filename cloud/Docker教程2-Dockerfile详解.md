# Docker 教程2 - Dockerfile 详解

> **重要程度：⭐⭐⭐ 必须掌握**  
> 构建自定义镜像的核心技能

## 1. 什么是 Dockerfile？⭐⭐⭐

Dockerfile 是一个文本文件，包含了构建 Docker 镜像的所有指令。通过 Dockerfile，你可以：
- 定义基础镜像
- 安装依赖
- 复制代码
- 配置环境
- 指定启动命令

**流程：**
```
Dockerfile → docker build → 镜像 → docker run → 容器
```

---

## 2. Dockerfile 基本指令 ⭐⭐⭐

> 必须掌握：FROM、WORKDIR、COPY、RUN、ENV、EXPOSE、CMD

### 2.1 FROM - 基础镜像

```dockerfile
# 指定基础镜像（必须是第一条指令）
FROM ubuntu:22.04

# 使用官方 Go 镜像
FROM golang:1.21

# 使用轻量级 Alpine 镜像
FROM alpine:3.18

# 多阶段构建时可以有多个 FROM
FROM golang:1.21 AS builder
```

### 2.2 WORKDIR - 工作目录

```dockerfile
# 设置工作目录（如果不存在会自动创建）
WORKDIR /app

# 后续的 RUN、CMD、COPY 等指令都在这个目录下执行
```

### 2.3 COPY - 复制文件

```dockerfile
# 复制单个文件
COPY main.go .

# 复制目录
COPY src/ ./src/

# 复制多个文件
COPY package.json package-lock.json ./

# 复制所有文件
COPY . .
```

### 2.4 ADD - 复制文件（增强版）

```dockerfile
# ADD 比 COPY 多两个功能：
# 1. 可以解压 tar 文件
ADD app.tar.gz /app/

# 2. 可以从 URL 下载
ADD https://example.com/file.txt /app/

# 一般推荐用 COPY，除非需要上述功能
```

### 2.5 RUN - 执行命令

```dockerfile
# 执行 shell 命令
RUN apt-get update && apt-get install -y curl

# 多条命令用 && 连接，减少镜像层数
RUN apt-get update && \
    apt-get install -y \
    curl \
    vim \
    git && \
    rm -rf /var/lib/apt/lists/*

# exec 格式
RUN ["apt-get", "install", "-y", "curl"]
```

### 2.6 ENV - 环境变量

```dockerfile
# 设置环境变量
ENV APP_ENV=production
ENV PORT=8080

# 一次设置多个
ENV APP_ENV=production \
    PORT=8080 \
    DB_HOST=localhost
```

### 2.7 EXPOSE - 声明端口

```dockerfile
# 声明容器监听的端口（只是声明，不会自动映射）
EXPOSE 8080

# 声明多个端口
EXPOSE 80 443
```

### 2.8 CMD - 启动命令

```dockerfile
# 容器启动时执行的命令（只能有一个 CMD）
CMD ["./app"]

# shell 格式
CMD ./app

# 带参数
CMD ["nginx", "-g", "daemon off;"]
```

### 2.9 ENTRYPOINT - 入口点

```dockerfile
# 与 CMD 类似，但更难被覆盖
ENTRYPOINT ["./app"]

# ENTRYPOINT + CMD 组合使用
ENTRYPOINT ["./app"]
CMD ["--port", "8080"]
# 相当于执行：./app --port 8080
# 运行时可以覆盖 CMD：docker run myapp --port 9090
```

### 2.10 ARG - 构建参数

```dockerfile
# 定义构建时的变量
ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}

# 构建时传入
# docker build --build-arg GO_VERSION=1.22 .
```

### 2.11 VOLUME - 数据卷

```dockerfile
# 声明数据卷挂载点
VOLUME /data
VOLUME ["/data", "/logs"]
```

---

## 3. 实战：构建 Go 应用镜像

### 3.1 准备 Go 应用

创建一个简单的 Go Web 应用：

```go
// main.go
package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, Docker! 🐳\n")
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "OK\n")
    })

    fmt.Printf("Server starting on port %s...\n", port)
    http.ListenAndServe(":"+port, nil)
}
```

### 3.2 简单版 Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.21

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
```

构建并运行：
```bash
# 构建镜像
docker build -t my-go-app .

# 查看镜像大小
docker images my-go-app
# 大约 800MB+ （因为包含了整个 Go 编译环境）

# 运行容器
docker run -d -p 8080:8080 --name my-app my-go-app

# 测试
curl http://localhost:8080
```

### 3.3 优化版 Dockerfile（多阶段构建）

```dockerfile
# 阶段1：构建
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 先复制 go.mod 和 go.sum，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 阶段2：运行
FROM alpine:3.18

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```

构建并对比：
```bash
# 构建优化版
docker build -t my-go-app:optimized .

# 对比大小
docker images | grep my-go-app
# my-go-app          latest      xxx   800MB
# my-go-app          optimized   xxx   15MB  ← 小了 50 倍！
```

### 3.4 最佳实践版 Dockerfile

```dockerfile
# 阶段1：构建
FROM golang:1.21-alpine AS builder

# 安装必要工具
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源码
COPY . .

# 编译（禁用 CGO，静态链接）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o main .

# 阶段2：运行
FROM scratch

# 复制时区信息和 CA 证书
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 复制二进制文件
COPY --from=builder /app/main /main

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8080

ENTRYPOINT ["/main"]
```

**优化点说明：**
- `scratch`：空镜像，最小化体积
- `-ldflags="-w -s"`：去除调试信息，减小体积
- `ca-certificates`：支持 HTTPS 请求
- `tzdata`：支持时区设置

---

## 4. 练习题

### 练习1：构建第一个镜像

创建一个简单的 Dockerfile，基于 `alpine` 镜像，安装 `curl`，并设置默认命令为 `curl --version`。

<details>
<summary>点击查看答案</summary>

```dockerfile
FROM alpine:3.18

RUN apk add --no-cache curl

CMD ["curl", "--version"]
```

```bash
# 构建
docker build -t my-curl .

# 运行
docker run my-curl
# 输出 curl 版本信息
```
</details>

---

### 练习2：构建 Python 应用

创建一个 Python Flask 应用并容器化：

1. 创建 `app.py`：
```python
from flask import Flask
app = Flask(__name__)

@app.route('/')
def hello():
    return 'Hello from Flask in Docker!'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

2. 创建 `requirements.txt`：
```
flask==3.0.0
```

3. 编写 Dockerfile 并构建运行

<details>
<summary>点击查看答案</summary>

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 5000

CMD ["python", "app.py"]
```

```bash
# 构建
docker build -t my-flask-app .

# 运行
docker run -d -p 5000:5000 --name flask-app my-flask-app

# 测试
curl http://localhost:5000
```
</details>

---

### 练习3：构建 Node.js 应用

创建一个 Node.js Express 应用并容器化：

1. 创建 `package.json`：
```json
{
  "name": "docker-node-app",
  "version": "1.0.0",
  "main": "index.js",
  "scripts": {
    "start": "node index.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}
```

2. 创建 `index.js`：
```javascript
const express = require('express');
const app = express();
const PORT = process.env.PORT || 3000;

app.get('/', (req, res) => {
  res.json({ message: 'Hello from Node.js in Docker!' });
});

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
```

3. 编写 Dockerfile（使用多阶段构建）

<details>
<summary>点击查看答案</summary>

```dockerfile
# 阶段1：安装依赖
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

# 阶段2：运行
FROM node:20-alpine

WORKDIR /app

COPY --from=builder /app/node_modules ./node_modules
COPY . .

ENV PORT=3000
EXPOSE 3000

CMD ["npm", "start"]
```

```bash
# 构建
docker build -t my-node-app .

# 运行
docker run -d -p 3000:3000 --name node-app my-node-app

# 测试
curl http://localhost:3000
```
</details>

---

### 练习4：使用 ARG 和 ENV

创建一个 Dockerfile，要求：
1. 使用 ARG 定义 Go 版本，默认 1.21
2. 使用 ENV 设置应用端口，默认 8080
3. 构建时可以通过 `--build-arg` 修改 Go 版本

<details>
<summary>点击查看答案</summary>

```dockerfile
ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .

ENV PORT=8080
EXPOSE ${PORT}

CMD ["./main"]
```

```bash
# 使用默认版本构建
docker build -t my-app .

# 使用指定版本构建
docker build --build-arg GO_VERSION=1.22 -t my-app:go1.22 .

# 运行时修改端口
docker run -d -p 9090:9090 -e PORT=9090 my-app
```
</details>

---

### 练习5：优化镜像大小

给定以下 Dockerfile，请优化它以减小镜像大小：

```dockerfile
FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN pip3 install flask

COPY app.py /app/app.py

WORKDIR /app

CMD ["python3", "app.py"]
```

<details>
<summary>点击查看答案</summary>

```dockerfile
# 优化版
FROM python:3.11-slim

WORKDIR /app

# 合并 RUN 命令，清理缓存
RUN pip install --no-cache-dir flask

COPY app.py .

CMD ["python", "app.py"]
```

**优化点：**
1. 使用 `python:3.11-slim` 代替 `ubuntu`（已包含 Python）
2. 合并 RUN 命令减少层数
3. 使用 `--no-cache-dir` 不保存 pip 缓存
4. 删除不必要的 apt 安装步骤

```bash
# 对比大小
# 优化前：约 500MB
# 优化后：约 150MB
```
</details>

---

## 5. .dockerignore 文件

类似 `.gitignore`，用于排除不需要复制到镜像的文件：

```
# .dockerignore
.git
.gitignore
*.md
Dockerfile
docker-compose.yml
node_modules
__pycache__
*.pyc
.env
.DS_Store
```

**好处：**
- 减小构建上下文大小
- 加快构建速度
- 避免敏感文件进入镜像

---

## 6. 构建缓存

Docker 会缓存每一层的构建结果，合理利用缓存可以加快构建速度。

**缓存原则：**
- 把不常变化的指令放前面
- 把经常变化的指令放后面

```dockerfile
# ❌ 不好的写法：每次代码变化都要重新安装依赖
COPY . .
RUN npm install

# ✅ 好的写法：只有 package.json 变化才重新安装
COPY package*.json ./
RUN npm install
COPY . .
```

**强制不使用缓存：**
```bash
docker build --no-cache -t my-app .
```

---

## 7. 本章小结

**常用指令：**
| 指令 | 作用 |
|------|------|
| `FROM` | 指定基础镜像 |
| `WORKDIR` | 设置工作目录 |
| `COPY` | 复制文件 |
| `RUN` | 执行命令 |
| `ENV` | 设置环境变量 |
| `EXPOSE` | 声明端口 |
| `CMD` | 启动命令 |
| `ENTRYPOINT` | 入口点 |
| `ARG` | 构建参数 |

**最佳实践：**
1. 使用多阶段构建减小镜像体积
2. 合理利用构建缓存
3. 使用 `.dockerignore` 排除无关文件
4. 合并 RUN 命令减少层数
5. 使用轻量级基础镜像（alpine、slim、scratch）

**下一章预告：** Docker Compose 多容器编排
