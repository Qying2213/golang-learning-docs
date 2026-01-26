# Docker 集成 ⭐⭐⭐

> 本章目标：掌握在 CI/CD 中构建、推送、部署 Docker 镜像

---

## 1. Docker 在 CI/CD 中的作用

```
传统部署：
代码 → 编译 → 上传到服务器 → 安装依赖 → 运行
                              ↓
                    每台服务器环境可能不同，容易出问题

Docker 部署：
代码 → 编译 → 打包成镜像 → 推送到仓库 → 服务器拉取运行
                              ↓
                    环境一致，到处运行
```

---

## 2. Docker 镜像仓库

### 2.1 常用仓库

| 仓库 | 说明 | 适用场景 |
|------|------|----------|
| **Docker Hub** | 官方仓库，免费 | 个人项目、开源项目 |
| **GitHub Container Registry** | GitHub 提供 | GitHub 项目 |
| **阿里云容器镜像服务** | 国内访问快 | 国内部署 |
| **Harbor** | 自建私有仓库 | 企业内部 |

### 2.2 Docker Hub 准备

1. 注册账号：https://hub.docker.com
2. 创建 Access Token：
   - Account Settings → Security → New Access Token
   - 保存 Token（只显示一次）

---

## 3. GitHub Actions 中使用 Docker

### 3.1 基础：构建镜像

```yaml
jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: 拉取代码
        uses: actions/checkout@v3
      
      - name: 构建镜像
        run: docker build -t myapp:latest .
```

### 3.2 登录 Docker Hub

```yaml
- name: 登录 Docker Hub
  uses: docker/login-action@v2
  with:
    username: ${{ secrets.DOCKER_USERNAME }}
    password: ${{ secrets.DOCKER_PASSWORD }}
```

**需要配置的 Secrets：**
- `DOCKER_USERNAME`：Docker Hub 用户名
- `DOCKER_PASSWORD`：Access Token（不是密码）

### 3.3 构建并推送

```yaml
- name: 构建并推送
  uses: docker/build-push-action@v4
  with:
    context: .
    push: true
    tags: username/myapp:latest
```

### 3.4 完整示例

```yaml
name: Docker Build and Push

on:
  push:
    branches: [ main ]

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: 拉取代码
        uses: actions/checkout@v3
      
      - name: 设置 Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: 登录 Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: 构建并推送
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/myapp:latest
            ${{ secrets.DOCKER_USERNAME }}/myapp:${{ github.sha }}
```

---

## 4. 镜像标签策略

### 4.1 常用标签

```yaml
tags: |
  # 最新版本
  username/myapp:latest
  
  # Git commit SHA（精确版本）
  username/myapp:abc1234
  
  # 分支名
  username/myapp:main
  
  # 版本号（从 tag 获取）
  username/myapp:v1.0.0
  
  # 日期
  username/myapp:2024-01-20
```

### 4.2 自动生成标签

```yaml
- name: 提取元数据
  id: meta
  uses: docker/metadata-action@v4
  with:
    images: ${{ secrets.DOCKER_USERNAME }}/myapp
    tags: |
      type=ref,event=branch
      type=sha,prefix=
      type=semver,pattern={{version}}
      type=raw,value=latest,enable={{is_default_branch}}

- name: 构建并推送
  uses: docker/build-push-action@v4
  with:
    context: .
    push: true
    tags: ${{ steps.meta.outputs.tags }}
```

**效果：**
- push 到 main：`myapp:main`, `myapp:abc1234`, `myapp:latest`
- push tag v1.0.0：`myapp:1.0.0`, `myapp:abc1234`

---

## 5. 多平台构建

### 5.1 为什么需要多平台？

```
你的 Mac M1/M2 是 ARM 架构
服务器通常是 AMD64 架构
        ↓
需要构建不同架构的镜像
```

### 5.2 配置多平台构建

```yaml
- name: 设置 QEMU（模拟器）
  uses: docker/setup-qemu-action@v2

- name: 设置 Docker Buildx
  uses: docker/setup-buildx-action@v2

- name: 构建并推送（多平台）
  uses: docker/build-push-action@v4
  with:
    context: .
    push: true
    platforms: linux/amd64,linux/arm64
    tags: username/myapp:latest
```

---

## 6. 构建缓存

### 6.1 为什么需要缓存？

```
没有缓存：每次构建都从头开始，耗时 5 分钟
有缓存：只构建变化的部分，耗时 30 秒
```

### 6.2 使用 GitHub Actions 缓存

```yaml
- name: 构建并推送
  uses: docker/build-push-action@v4
  with:
    context: .
    push: true
    tags: username/myapp:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

### 6.3 使用 Registry 缓存

```yaml
- name: 构建并推送
  uses: docker/build-push-action@v4
  with:
    context: .
    push: true
    tags: username/myapp:latest
    cache-from: type=registry,ref=username/myapp:buildcache
    cache-to: type=registry,ref=username/myapp:buildcache,mode=max
```

---

## 7. Dockerfile 最佳实践

### 7.1 Go 项目 Dockerfile

```dockerfile
# ============ 构建阶段 ============
FROM golang:1.21-alpine AS builder

# 安装必要工具
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 先复制依赖文件（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.Version=${VERSION}" \
    -o /app/main .

# ============ 运行阶段 ============
FROM scratch

# 复制时区信息
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai

# 复制 CA 证书（HTTPS 需要）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 复制可执行文件
COPY --from=builder /app/main /main

# 暴露端口
EXPOSE 8080

# 启动命令
ENTRYPOINT ["/main"]
```

### 7.2 关键优化点

```dockerfile
# 1. 多阶段构建：减小镜像体积
FROM golang:1.21-alpine AS builder  # 构建阶段
FROM scratch                         # 运行阶段（最小镜像）

# 2. 利用缓存：先复制依赖文件
COPY go.mod go.sum ./
RUN go mod download
COPY . .  # 源码变化不会影响依赖缓存

# 3. 减小二进制体积
go build -ldflags="-s -w"  # -s 去掉符号表，-w 去掉调试信息

# 4. 静态编译
CGO_ENABLED=0  # 不依赖 C 库，可以在 scratch 中运行
```

### 7.3 镜像大小对比

```
golang:1.21          ~800MB
golang:1.21-alpine   ~250MB
alpine               ~5MB
scratch              ~0MB（只有你的程序）

最终镜像大小：~10-20MB（取决于你的程序）
```

---

## 8. 部署到服务器

### 8.1 SSH 部署

```yaml
deploy:
  runs-on: ubuntu-latest
  needs: docker
  steps:
    - name: 部署到服务器
      uses: appleboy/ssh-action@v1.0.0
      with:
        host: ${{ secrets.SERVER_HOST }}
        username: ${{ secrets.SERVER_USER }}
        key: ${{ secrets.SERVER_SSH_KEY }}
        script: |
          # 拉取最新镜像
          docker pull ${{ secrets.DOCKER_USERNAME }}/myapp:latest
          
          # 停止并删除旧容器
          docker stop myapp || true
          docker rm myapp || true
          
          # 启动新容器
          docker run -d \
            --name myapp \
            -p 8080:8080 \
            -e DATABASE_URL=${{ secrets.DATABASE_URL }} \
            --restart unless-stopped \
            ${{ secrets.DOCKER_USERNAME }}/myapp:latest
          
          # 清理旧镜像
          docker image prune -f
```

### 8.2 Docker Compose 部署

服务器上的 `docker-compose.yml`：
```yaml
version: '3.8'
services:
  app:
    image: username/myapp:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=${DATABASE_URL}
    restart: unless-stopped
```

CI/CD 部署脚本：
```yaml
script: |
  cd /opt/myapp
  docker-compose pull
  docker-compose up -d
  docker image prune -f
```

---

## 9. 使用 GitHub Container Registry

### 9.1 为什么用 GHCR？

- 与 GitHub 深度集成
- 私有仓库免费
- 权限管理方便

### 9.2 配置

```yaml
- name: 登录 GitHub Container Registry
  uses: docker/login-action@v2
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: 构建并推送
  uses: docker/build-push-action@v4
  with:
    context: .
    push: true
    tags: ghcr.io/${{ github.repository }}:latest
```

**注意：** `GITHUB_TOKEN` 是自动提供的，不需要手动配置。

---

## 10. 完整的 Docker CI/CD 配置

```yaml
name: Docker CI/CD

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

env:
  REGISTRY: docker.io
  IMAGE_NAME: ${{ secrets.DOCKER_USERNAME }}/myapp

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
      - name: 拉取代码
        uses: actions/checkout@v3
      
      - name: 设置 QEMU
        uses: docker/setup-qemu-action@v2
      
      - name: 设置 Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: 登录 Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: 提取元数据
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=sha,prefix=
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=raw,value=latest,enable={{is_default_branch}}
      
      - name: 构建并推送
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
  
  deploy:
    runs-on: ubuntu-latest
    needs: build-and-push
    if: github.ref == 'refs/heads/main'
    
    steps:
      - name: 部署到服务器
        uses: appleboy/ssh-action@v1.0.0
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            docker pull ${{ env.IMAGE_NAME }}:latest
            docker stop myapp || true
            docker rm myapp || true
            docker run -d \
              --name myapp \
              -p 8080:8080 \
              --restart unless-stopped \
              ${{ env.IMAGE_NAME }}:latest
            docker image prune -f
```

---

## 11. 练习题

### 练习1：理解标签

**Q：以下配置会生成哪些标签？**
```yaml
tags: |
  type=ref,event=branch
  type=sha,prefix=
  type=raw,value=latest
```
假设：分支是 main，commit SHA 是 abc1234

<details>
<summary>查看答案</summary>

会生成 3 个标签：
- `myapp:main`（分支名）
- `myapp:abc1234`（commit SHA）
- `myapp:latest`（固定值）

</details>

### 练习2：修复 Dockerfile

**以下 Dockerfile 有什么问题？**
```dockerfile
FROM golang:1.21
COPY . .
RUN go build -o main .
CMD ["./main"]
```

<details>
<summary>查看答案</summary>

问题：
1. 没有使用多阶段构建，镜像太大（~800MB）
2. 没有先复制 go.mod，无法利用缓存
3. 没有设置 WORKDIR

改进：
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
```

</details>

### 练习3：添加健康检查

**任务：** 在 docker run 命令中添加健康检查

<details>
<summary>查看答案</summary>

```bash
docker run -d \
  --name myapp \
  -p 8080:8080 \
  --health-cmd="wget -q --spider http://localhost:8080/health || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  --restart unless-stopped \
  myapp:latest
```

或者在 Dockerfile 中：
```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD wget -q --spider http://localhost:8080/health || exit 1
```

</details>

---

## 12. 本章小结

```
Docker CI/CD 流程：
代码 → 构建镜像 → 推送仓库 → 部署服务器

关键 Actions：
├── docker/login-action：登录仓库
├── docker/setup-buildx-action：设置构建器
├── docker/metadata-action：生成标签
├── docker/build-push-action：构建推送
└── appleboy/ssh-action：SSH 部署

优化技巧：
├── 多阶段构建：减小镜像体积
├── 构建缓存：加速构建
├── 多平台构建：支持不同架构
└── 自动标签：版本管理

镜像仓库选择：
├── Docker Hub：通用
├── GHCR：GitHub 项目
└── 阿里云：国内部署
```

---

下一章：[05-部署策略](./05-部署策略.md)
