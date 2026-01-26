# Docker 教程3 - Docker Compose 多容器编排

## 1. 什么是 Docker Compose？

Docker Compose 是一个用于定义和运行多容器应用的工具。通过一个 YAML 文件来配置所有服务，然后一条命令启动所有容器。

**使用场景：**
- 本地开发环境（应用 + 数据库 + 缓存）
- 测试环境搭建
- 单机部署简单应用

**没有 Compose 时：**
```bash
# 需要手动启动每个容器
docker run -d --name mysql ...
docker run -d --name redis ...
docker run -d --name app --link mysql --link redis ...
```

**有了 Compose：**
```bash
# 一条命令启动所有服务
docker-compose up -d
```

---

## 2. 安装 Docker Compose

Docker Desktop 已经自带 Compose，验证安装：

```bash
# 新版本（推荐）
docker compose version

# 旧版本
docker-compose --version
```

---

## 3. docker-compose.yml 基础

### 3.1 基本结构

```yaml
version: '3.8'  # Compose 文件版本

services:       # 定义服务
  web:          # 服务名
    image: nginx
    ports:
      - "8080:80"
  
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: 123456

volumes:        # 定义数据卷（可选）
  db-data:

networks:       # 定义网络（可选）
  app-network:
```

### 3.2 常用配置项

```yaml
services:
  app:
    # 镜像名称
    image: nginx:latest
    
    # 或者从 Dockerfile 构建
    build:
      context: .
      dockerfile: Dockerfile
    
    # 容器名称
    container_name: my-app
    
    # 端口映射
    ports:
      - "8080:80"
      - "443:443"
    
    # 环境变量
    environment:
      - NODE_ENV=production
      - DB_HOST=db
    # 或者
    environment:
      NODE_ENV: production
      DB_HOST: db
    
    # 从文件读取环境变量
    env_file:
      - .env
    
    # 数据卷挂载
    volumes:
      - ./data:/app/data           # 绑定挂载
      - db-data:/var/lib/mysql     # 命名卷
    
    # 依赖关系（先启动 db，再启动 app）
    depends_on:
      - db
      - redis
    
    # 重启策略
    restart: always  # no | always | on-failure | unless-stopped
    
    # 网络
    networks:
      - app-network
    
    # 资源限制
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
    
    # 健康检查
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

---

## 4. 常用命令

### 4.1 启动和停止

```bash
# 启动所有服务（前台运行）
docker compose up

# 后台运行
docker compose up -d

# 启动指定服务
docker compose up -d web

# 停止所有服务
docker compose stop

# 停止并删除容器、网络
docker compose down

# 停止并删除容器、网络、数据卷
docker compose down -v
```

### 4.2 查看状态

```bash
# 查看运行中的服务
docker compose ps

# 查看日志
docker compose logs

# 实时查看日志
docker compose logs -f

# 查看指定服务日志
docker compose logs -f web
```

### 4.3 其他操作

```bash
# 重新构建镜像
docker compose build

# 构建并启动
docker compose up -d --build

# 进入容器
docker compose exec web bash

# 执行命令
docker compose exec db mysql -uroot -p

# 重启服务
docker compose restart web

# 查看服务配置
docker compose config
```

---

## 5. 实战案例

### 5.1 案例1：Go + MySQL + Redis

**项目结构：**
```
my-app/
├── docker-compose.yml
├── Dockerfile
├── main.go
├── go.mod
└── .env
```

**docker-compose.yml：**
```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: go-app
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=123456
      - DB_NAME=myapp
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - mysql
      - redis
    restart: unless-stopped
    networks:
      - app-network

  mysql:
    image: mysql:8.0
    container_name: mysql
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: myapp
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    restart: unless-stopped
    networks:
      - app-network

volumes:
  mysql-data:
  redis-data:

networks:
  app-network:
    driver: bridge
```

**Dockerfile：**
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**启动：**
```bash
# 启动所有服务
docker compose up -d

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f app
```

### 5.2 案例2：Nginx + PHP + MySQL（LNMP）

```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    container_name: nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./www:/var/www/html
    depends_on:
      - php
    networks:
      - lnmp

  php:
    image: php:8.2-fpm
    container_name: php
    volumes:
      - ./www:/var/www/html
    networks:
      - lnmp

  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: app
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - lnmp

volumes:
  mysql-data:

networks:
  lnmp:
```

### 5.3 案例3：前后端分离项目

```yaml
version: '3.8'

services:
  # 前端（Vue/React）
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    networks:
      - app-network

  # 后端（Go/Node）
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    networks:
      - app-network

  # 数据库
  postgres:
    image: postgres:15-alpine
    container_name: postgres
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: 123456
      POSTGRES_DB: myapp
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - app-network

  # 缓存
  redis:
    image: redis:7-alpine
    container_name: redis
    networks:
      - app-network

  # Nginx 反向代理
  nginx:
    image: nginx:alpine
    container_name: nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - frontend
      - backend
    networks:
      - app-network

volumes:
  postgres-data:

networks:
  app-network:
```

---

## 6. 练习题

### 练习1：基础 Compose

创建一个 `docker-compose.yml`，启动一个 Nginx 服务：
- 端口映射 8080:80
- 挂载本地 `./html` 目录到 `/usr/share/nginx/html`

<details>
<summary>点击查看答案</summary>

```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "8080:80"
    volumes:
      - ./html:/usr/share/nginx/html
```

```bash
# 创建测试页面
mkdir html
echo "<h1>Hello Docker Compose!</h1>" > html/index.html

# 启动
docker compose up -d

# 访问 http://localhost:8080
```
</details>

---

### 练习2：WordPress 博客

使用 Docker Compose 搭建 WordPress 博客：
- WordPress 服务，端口 8080
- MySQL 数据库
- 数据持久化

<details>
<summary>点击查看答案</summary>

```yaml
version: '3.8'

services:
  wordpress:
    image: wordpress:latest
    container_name: wordpress
    ports:
      - "8080:80"
    environment:
      WORDPRESS_DB_HOST: mysql
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      WORDPRESS_DB_NAME: wordpress
    volumes:
      - wordpress-data:/var/www/html
    depends_on:
      - mysql
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    container_name: wordpress-mysql
    environment:
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped

volumes:
  wordpress-data:
  mysql-data:
```

```bash
docker compose up -d
# 访问 http://localhost:8080 完成 WordPress 安装
```
</details>

---

### 练习3：开发环境

创建一个本地开发环境，包含：
- MySQL 8.0（端口 3306）
- Redis（端口 6379）
- Adminer（数据库管理工具，端口 8081）

<details>
<summary>点击查看答案</summary>

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: dev-mysql
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: devdb
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: dev-redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    restart: unless-stopped

  adminer:
    image: adminer
    container_name: dev-adminer
    ports:
      - "8081:8080"
    depends_on:
      - mysql
    restart: unless-stopped

volumes:
  mysql-data:
  redis-data:
```

```bash
docker compose up -d

# 访问 Adminer: http://localhost:8081
# 服务器: mysql, 用户名: root, 密码: 123456
```
</details>

---

### 练习4：构建自定义镜像

创建一个项目，包含：
1. 一个简单的 Go Web 应用
2. Dockerfile
3. docker-compose.yml（构建并运行应用，连接 MySQL）

<details>
<summary>点击查看答案</summary>

**main.go：**
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Printf("Database connection error: %v", err)
    } else {
        defer db.Close()
        log.Println("Database connected!")
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Go + MySQL!")
    })

    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Dockerfile：**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**docker-compose.yml：**
```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: go-app
    ports:
      - "8080:8080"
    environment:
      DB_HOST: mysql
      DB_USER: root
      DB_PASSWORD: 123456
      DB_NAME: myapp
    depends_on:
      - mysql
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: myapp
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql-data:
```

```bash
# 初始化 Go 模块
go mod init myapp
go get github.com/go-sql-driver/mysql

# 构建并启动
docker compose up -d --build

# 测试
curl http://localhost:8080
```
</details>

---

### 练习5：环境变量管理

修改练习4，使用 `.env` 文件管理环境变量。

<details>
<summary>点击查看答案</summary>

**.env：**
```
DB_HOST=mysql
DB_USER=root
DB_PASSWORD=123456
DB_NAME=myapp
MYSQL_ROOT_PASSWORD=123456
```

**docker-compose.yml：**
```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: go-app
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - mysql
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql-data:
```

```bash
# .env 文件会自动加载
docker compose up -d
```
</details>

---

## 7. 网络详解

### 7.1 默认网络

Compose 会自动创建一个网络，所有服务都在这个网络中，可以通过服务名互相访问。

```yaml
services:
  app:
    image: myapp
    # 可以通过 "db" 访问数据库
    environment:
      DB_HOST: db
  
  db:
    image: mysql
```

### 7.2 自定义网络

```yaml
services:
  frontend:
    networks:
      - frontend-network
  
  backend:
    networks:
      - frontend-network
      - backend-network
  
  db:
    networks:
      - backend-network

networks:
  frontend-network:
  backend-network:
```

---

## 8. 本章小结

**核心概念：**
- `docker-compose.yml`：定义多容器应用
- `services`：定义各个服务
- `volumes`：数据持久化
- `networks`：容器网络

**常用命令：**
| 命令 | 作用 |
|------|------|
| `docker compose up -d` | 后台启动 |
| `docker compose down` | 停止并删除 |
| `docker compose ps` | 查看状态 |
| `docker compose logs -f` | 查看日志 |
| `docker compose exec` | 进入容器 |
| `docker compose build` | 构建镜像 |

**下一章预告：** Kubernetes 入门
