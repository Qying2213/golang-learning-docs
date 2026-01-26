# Docker 教程1 - 基础入门

> **重要程度：⭐⭐⭐ 必须掌握**  
> 本章是 Docker 的核心基础，必须全部掌握

## 1. Docker 简介 ⭐⭐⭐

### 1.1 什么是 Docker？

Docker 是一个**容器化平台**，可以把应用程序和它的依赖打包到一个轻量级、可移植的容器中。

**传统部署的痛点：**
- "在我电脑上能跑啊" - 环境不一致
- 安装依赖麻烦 - 不同项目依赖冲突
- 部署复杂 - 每台服务器都要配置环境

**Docker 解决方案：**
- 把应用 + 环境 + 依赖打包成镜像
- 在任何安装了 Docker 的机器上都能运行
- 一次构建，到处运行

### 1.2 核心概念

```
+------------------+
|     镜像 Image    |  ← 只读模板，类似"安装包"
+------------------+
         ↓ 运行
+------------------+
|    容器 Container |  ← 镜像的运行实例，类似"运行中的程序"
+------------------+

+------------------+
|  仓库 Registry   |  ← 存放镜像的地方，如 Docker Hub
+------------------+
```

**类比理解：**
- 镜像 = 类（Class）
- 容器 = 对象（Object）
- 一个镜像可以创建多个容器

### 1.3 Docker vs 虚拟机

| 特性 | Docker 容器 | 虚拟机 |
|------|------------|--------|
| 启动速度 | 秒级 | 分钟级 |
| 资源占用 | MB 级 | GB 级 |
| 性能 | 接近原生 | 有损耗 |
| 隔离性 | 进程级 | 系统级 |
| 系统支持 | 共享宿主机内核 | 完整 OS |

---

## 2. 安装 Docker ⭐⭐⭐

### 2.1 macOS 安装

```bash
# 方式1：使用 Homebrew
brew install --cask docker

# 方式2：下载 Docker Desktop
# 访问 https://www.docker.com/products/docker-desktop 下载安装

# 安装后启动 Docker Desktop 应用
```

### 2.2 验证安装

```bash
# 查看版本
docker --version
# Docker version 24.0.7, build afdd53b

# 查看详细信息
docker info

# 运行测试容器
docker run hello-world
```

如果看到 "Hello from Docker!" 说明安装成功！

---

## 3. 镜像操作 ⭐⭐⭐

### 3.1 搜索镜像 ⭐

> 了解即可，实际工作中很少用，一般直接去 Docker Hub 网站搜索

```bash
# 在 Docker Hub 搜索镜像
docker search nginx

# 结果说明
# NAME: 镜像名
# DESCRIPTION: 描述
# STARS: 星标数（越多越受欢迎）
# OFFICIAL: 是否官方镜像
```

### 3.2 拉取镜像 ⭐⭐⭐

```bash
# 拉取最新版本
docker pull nginx

# 拉取指定版本（标签）
docker pull nginx:1.24

# 拉取指定平台
docker pull --platform linux/amd64 nginx
```

### 3.3 查看本地镜像 ⭐⭐⭐

```bash
docker images
# 或
docker image ls

# 输出示例：
# REPOSITORY   TAG       IMAGE ID       CREATED        SIZE
# nginx        latest    a6bd71f48f68   2 weeks ago    187MB
# nginx        1.24      5e4a2f0e0c8b   3 months ago   142MB
```

**字段说明：**
- `REPOSITORY`: 镜像名
- `TAG`: 标签/版本
- `IMAGE ID`: 镜像唯一 ID
- `SIZE`: 镜像大小

### 3.4 删除镜像 ⭐⭐

```bash
# 按名称删除
docker rmi nginx

# 按 ID 删除
docker rmi a6bd71f48f68

# 强制删除（即使有容器在使用）
docker rmi -f nginx

# 删除所有未使用的镜像
docker image prune -a
```

---

## 4. 容器操作 ⭐⭐⭐

> **核心重点！** 这是 Docker 最常用的操作，必须熟练掌握

### 4.1 创建并运行容器 ⭐⭐⭐

> **最重要的命令！** `docker run` 是使用频率最高的命令

```bash
# 基本语法
docker run [选项] 镜像名 [命令]

# 运行 nginx 容器
docker run nginx

# 后台运行（-d = detach）
docker run -d nginx

# 指定容器名称
docker run -d --name my-nginx nginx

# 端口映射（-p 宿主机端口:容器端口）
docker run -d -p 8080:80 --name my-nginx nginx
# 访问 http://localhost:8080 就能看到 nginx 页面
```

### 4.2 常用运行选项 ⭐⭐⭐

> 必须掌握：`-d`、`-p`、`--name`、`-e`、`-v`

| 选项 | 说明 | 示例 |
|------|------|------|
| `-d` | 后台运行 | `docker run -d nginx` |
| `-p` | 端口映射 | `-p 8080:80` |
| `--name` | 容器名称 | `--name my-app` |
| `-e` | 环境变量 | `-e MYSQL_ROOT_PASSWORD=123456` |
| `-v` | 挂载卷 | `-v /host/path:/container/path` |
| `--rm` | 退出后自动删除 | `docker run --rm nginx` |
| `-it` | 交互式终端 | `docker run -it ubuntu bash` |

### 4.3 查看容器 ⭐⭐⭐

```bash
# 查看运行中的容器
docker ps

# 查看所有容器（包括已停止的）
docker ps -a

# 只显示容器 ID
docker ps -q

# 输出示例：
# CONTAINER ID   IMAGE   COMMAND                  STATUS          PORTS                  NAMES
# abc123def456   nginx   "/docker-entrypoint.…"   Up 2 minutes    0.0.0.0:8080->80/tcp   my-nginx
```

### 4.4 容器生命周期 ⭐⭐⭐

```bash
# 停止容器
docker stop my-nginx

# 启动已停止的容器
docker start my-nginx

# 重启容器
docker restart my-nginx

# 暂停容器
docker pause my-nginx

# 恢复容器
docker unpause my-nginx

# 删除容器（必须先停止）
docker rm my-nginx

# 强制删除运行中的容器
docker rm -f my-nginx

# 删除所有已停止的容器
docker container prune
```

### 4.5 进入容器 ⭐⭐⭐

> 调试必备技能

```bash
# 进入运行中的容器
docker exec -it my-nginx bash

# 如果容器没有 bash，用 sh
docker exec -it my-nginx sh

# 执行单条命令
docker exec my-nginx cat /etc/nginx/nginx.conf
```

### 4.6 查看容器日志 ⭐⭐⭐

> 排查问题的第一步，必须掌握

```bash
# 查看日志
docker logs my-nginx

# 实时查看日志（类似 tail -f）
docker logs -f my-nginx

# 查看最后 100 行
docker logs --tail 100 my-nginx

# 显示时间戳
docker logs -t my-nginx
```

### 4.7 查看容器详情 ⭐⭐

> 了解即可，偶尔用于查看容器配置

```bash
# 查看容器详细信息
docker inspect my-nginx

# 查看容器 IP 地址
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' my-nginx

# 查看容器资源使用
docker stats my-nginx
```

---

## 5. 练习题 ⭐⭐⭐

> 必须动手做！只看不练等于没学

### 练习1：运行第一个容器

1. 拉取 `nginx` 镜像
2. 运行一个 nginx 容器，命名为 `web-server`，将容器的 80 端口映射到宿主机的 8888 端口
3. 在浏览器访问 `http://localhost:8888`，确认能看到 nginx 欢迎页
4. 查看容器日志
5. 停止并删除容器

<details>
<summary>点击查看答案</summary>

```bash
# 1. 拉取镜像
docker pull nginx

# 2. 运行容器
docker run -d -p 8888:80 --name web-server nginx

# 3. 浏览器访问 http://localhost:8888

# 4. 查看日志
docker logs web-server

# 5. 停止并删除
docker stop web-server
docker rm web-server
```
</details>

---

### 练习2：运行 MySQL 容器

运行一个 MySQL 容器，要求：
- 容器名：`my-mysql`
- 后台运行
- 端口映射：3306 -> 3306
- 设置 root 密码为 `123456`
- 设置默认数据库为 `testdb`

提示：MySQL 镜像的环境变量
- `MYSQL_ROOT_PASSWORD`: root 密码
- `MYSQL_DATABASE`: 默认创建的数据库

<details>
<summary>点击查看答案</summary>

```bash
docker run -d \
  --name my-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=testdb \
  mysql:8.0

# 验证：连接 MySQL
docker exec -it my-mysql mysql -uroot -p123456

# 在 MySQL 中查看数据库
SHOW DATABASES;
```
</details>

---

### 练习3：运行 Redis 容器

1. 运行一个 Redis 容器，命名为 `my-redis`，端口映射 6379
2. 进入容器，使用 `redis-cli` 测试
3. 设置一个 key：`name` = `docker-test`
4. 获取这个 key 的值

<details>
<summary>点击查看答案</summary>

```bash
# 1. 运行 Redis
docker run -d --name my-redis -p 6379:6379 redis

# 2. 进入容器
docker exec -it my-redis redis-cli

# 3. 设置 key（在 redis-cli 中）
SET name docker-test

# 4. 获取 key
GET name
# 输出：docker-test

# 退出
exit
```
</details>

---

### 练习4：容器管理

1. 查看当前所有运行中的容器
2. 查看所有容器（包括已停止的）
3. 停止所有运行中的容器
4. 删除所有已停止的容器
5. 删除所有未使用的镜像

<details>
<summary>点击查看答案</summary>

```bash
# 1. 查看运行中的容器
docker ps

# 2. 查看所有容器
docker ps -a

# 3. 停止所有运行中的容器
docker stop $(docker ps -q)

# 4. 删除所有已停止的容器
docker container prune
# 或
docker rm $(docker ps -aq)

# 5. 删除所有未使用的镜像
docker image prune -a
```
</details>

---

### 练习5：交互式容器

1. 运行一个 Ubuntu 容器，进入交互式终端
2. 在容器中安装 `curl`
3. 使用 curl 访问 `https://www.baidu.com`
4. 退出容器

<details>
<summary>点击查看答案</summary>

```bash
# 1. 运行 Ubuntu 并进入终端
docker run -it ubuntu bash

# 2. 安装 curl（在容器内）
apt update && apt install -y curl

# 3. 访问百度
curl https://www.baidu.com

# 4. 退出
exit
```
</details>

---

## 6. 数据持久化 ⭐⭐⭐

> 重要！容器删除后数据会丢失，必须学会持久化

容器删除后，里面的数据也会丢失。要持久化数据，需要使用**数据卷（Volume）**或**绑定挂载（Bind Mount）**。

### 6.1 数据卷（Volume）⭐⭐⭐

> 推荐方式，Docker 自动管理

```bash
# 创建数据卷
docker volume create my-data

# 查看所有数据卷
docker volume ls

# 查看数据卷详情
docker volume inspect my-data

# 使用数据卷运行容器
docker run -d \
  --name my-mysql \
  -v my-data:/var/lib/mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  mysql:8.0

# 删除数据卷
docker volume rm my-data

# 删除所有未使用的数据卷
docker volume prune
```

### 6.2 绑定挂载（Bind Mount）⭐⭐

> 开发环境常用，可以实时同步代码

```bash
# 将宿主机目录挂载到容器
docker run -d \
  --name my-nginx \
  -p 8080:80 \
  -v $(pwd)/html:/usr/share/nginx/html \
  nginx

# 现在修改宿主机的 ./html 目录，容器内会同步变化
```

### 6.3 练习：数据持久化

1. 创建一个目录 `~/docker-test/mysql-data`
2. 运行 MySQL 容器，将数据目录挂载到这个目录
3. 创建一个数据库和表，插入数据
4. 删除容器
5. 重新运行容器（使用相同的挂载），验证数据还在

<details>
<summary>点击查看答案</summary>

```bash
# 1. 创建目录
mkdir -p ~/docker-test/mysql-data

# 2. 运行 MySQL
docker run -d \
  --name my-mysql \
  -p 3306:3306 \
  -v ~/docker-test/mysql-data:/var/lib/mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  mysql:8.0

# 3. 创建数据
docker exec -it my-mysql mysql -uroot -p123456 -e "
CREATE DATABASE testdb;
USE testdb;
CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50));
INSERT INTO users VALUES (1, 'Docker');
"

# 4. 删除容器
docker rm -f my-mysql

# 5. 重新运行
docker run -d \
  --name my-mysql \
  -p 3306:3306 \
  -v ~/docker-test/mysql-data:/var/lib/mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  mysql:8.0

# 验证数据
docker exec -it my-mysql mysql -uroot -p123456 -e "SELECT * FROM testdb.users;"
# 输出：1 | Docker
```
</details>

---

## 7. 本章小结

**核心概念：**
- 镜像（Image）：只读模板
- 容器（Container）：镜像的运行实例
- 仓库（Registry）：存放镜像的地方

**常用命令：**
| 操作 | 命令 |
|------|------|
| 拉取镜像 | `docker pull` |
| 查看镜像 | `docker images` |
| 删除镜像 | `docker rmi` |
| 运行容器 | `docker run` |
| 查看容器 | `docker ps` |
| 停止容器 | `docker stop` |
| 删除容器 | `docker rm` |
| 进入容器 | `docker exec -it` |
| 查看日志 | `docker logs` |

**下一章预告：** Dockerfile 构建自定义镜像
