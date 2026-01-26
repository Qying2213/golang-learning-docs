# Nginx 教程 - 从入门到部署项目

## Nginx 是什么？

**Nginx = 高性能的 Web 服务器 + 反向代理**

它能做什么：
1. **静态文件服务**：托管前端 HTML/CSS/JS
2. **反向代理**：把请求转发给后端服务
3. **负载均衡**：把请求分发到多个服务器
4. **HTTPS**：配置 SSL 证书

```
用户请求
    ↓
  Nginx（80/443端口）
    ├── /           → 前端静态文件
    ├── /api        → 后端服务（8080端口）
    └── /images     → 图片文件
```

---

## 第一步：安装 Nginx

### Ubuntu/Debian
```bash
sudo apt update
sudo apt install nginx
```

### CentOS
```bash
sudo yum install nginx
```

### 验证安装
```bash
nginx -v
# nginx version: nginx/1.18.0
```

### 启动 Nginx
```bash
# 启动
sudo systemctl start nginx

# 开机自启
sudo systemctl enable nginx

# 查看状态
sudo systemctl status nginx
```

**访问服务器 IP，看到 "Welcome to nginx!" 就成功了！**

---

## 第二步：理解 Nginx 目录结构

```
/etc/nginx/                    # Nginx 配置目录
├── nginx.conf                 # 主配置文件
├── sites-available/           # 可用的站点配置
├── sites-enabled/             # 启用的站点配置（软链接）
└── conf.d/                    # 额外配置文件

/var/www/html/                 # 默认网站根目录
/var/log/nginx/                # 日志目录
├── access.log                 # 访问日志
└── error.log                  # 错误日志
```

---

## 第三步：理解配置文件结构

```nginx
# /etc/nginx/nginx.conf 主配置文件结构

# 全局配置
user www-data;
worker_processes auto;

# 事件配置
events {
    worker_connections 1024;
}

# HTTP 配置
http {
    # 通用设置
    include /etc/nginx/mime.types;
    
    # 日志格式
    access_log /var/log/nginx/access.log;
    
    # 引入站点配置
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
```

**重点：我们主要修改 `sites-available` 里的站点配置**

---

## 第四步：配置静态网站（前端）

### 4.1 创建网站目录

```bash
# 创建目录
sudo mkdir -p /var/www/mysite

# 创建测试页面
sudo nano /var/www/mysite/index.html
```

```html
<!DOCTYPE html>
<html>
<head>
    <title>我的网站</title>
</head>
<body>
    <h1>Hello, 秦阳！</h1>
    <p>Nginx 部署成功！</p>
</body>
</html>
```

### 4.2 创建 Nginx 配置

```bash
sudo nano /etc/nginx/sites-available/mysite
```

```nginx
server {
    listen 80;                          # 监听 80 端口
    server_name your-domain.com;        # 你的域名或 IP
    
    root /var/www/mysite;               # 网站根目录
    index index.html;                   # 默认首页
    
    location / {
        try_files $uri $uri/ =404;      # 尝试找文件，找不到返回 404
    }
}
```

### 4.3 启用配置

```bash
# 创建软链接到 sites-enabled
sudo ln -s /etc/nginx/sites-available/mysite /etc/nginx/sites-enabled/

# 测试配置是否正确
sudo nginx -t

# 重新加载配置
sudo systemctl reload nginx
```

**访问你的服务器 IP，就能看到页面了！**

---

## 第五步：反向代理（后端 API）

### 场景
```
前端：Nginx 托管静态文件
后端：Go 服务跑在 8080 端口

用户访问 /api/* → Nginx 转发到 localhost:8080
```

### 5.1 配置反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    root /var/www/mysite;
    index index.html;
    
    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;  # SPA 单页应用支持
    }
    
    # 后端 API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;  # 转发到后端
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 5.2 配置说明

| 配置 | 说明 |
|------|------|
| `proxy_pass` | 转发到哪个地址 |
| `proxy_set_header Host` | 传递原始域名 |
| `proxy_set_header X-Real-IP` | 传递用户真实 IP |
| `proxy_set_header X-Forwarded-For` | 传递代理链 IP |

### 5.3 测试

```bash
# 测试配置
sudo nginx -t

# 重新加载
sudo systemctl reload nginx
```

---

## 第六步：完整的前后端分离配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # ========== 前端配置 ==========
    root /var/www/mysite/dist;      # Vue/React 打包后的目录
    index index.html;
    
    # 前端路由（Vue Router / React Router）
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # ========== 后端 API ==========
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # ========== 静态资源缓存 ==========
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 30d;                    # 缓存 30 天
        add_header Cache-Control "public, immutable";
    }
    
    # ========== 日志 ==========
    access_log /var/log/nginx/mysite.access.log;
    error_log /var/log/nginx/mysite.error.log;
}
```

---

## 第七步：配置 HTTPS（SSL 证书）

### 7.1 安装 Certbot（免费证书）

```bash
# Ubuntu
sudo apt install certbot python3-certbot-nginx
```

### 7.2 获取证书

```bash
sudo certbot --nginx -d your-domain.com
```

**Certbot 会自动修改 Nginx 配置，添加 HTTPS！**

### 7.3 手动配置 HTTPS（如果有证书）

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;  # HTTP 跳转 HTTPS
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    # SSL 证书
    ssl_certificate /etc/nginx/ssl/your-domain.crt;
    ssl_certificate_key /etc/nginx/ssl/your-domain.key;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # 其他配置...
    root /var/www/mysite/dist;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 第八步：常用命令

```bash
# 测试配置文件语法
sudo nginx -t

# 启动
sudo systemctl start nginx

# 停止
sudo systemctl stop nginx

# 重启
sudo systemctl restart nginx

# 重新加载配置（不中断服务）
sudo systemctl reload nginx

# 查看状态
sudo systemctl status nginx

# 查看错误日志
sudo tail -f /var/log/nginx/error.log

# 查看访问日志
sudo tail -f /var/log/nginx/access.log
```

---

## 第九步：实战 - 部署 Go + Vue 项目

### 9.1 项目结构

```
服务器上：
/var/www/myproject/
├── dist/           # Vue 打包后的前端文件
│   ├── index.html
│   ├── css/
│   └── js/
└── backend/        # Go 后端
    └── main        # 编译后的可执行文件
```

### 9.2 上传文件

```bash
# 本地打包前端
cd frontend
npm run build

# 上传到服务器
scp -r dist/* root@your-server:/var/www/myproject/dist/

# 本地编译后端（Linux 版本）
cd backend
GOOS=linux GOARCH=amd64 go build -o main

# 上传到服务器
scp main root@your-server:/var/www/myproject/backend/
```

### 9.3 启动后端服务

```bash
# 在服务器上
cd /var/www/myproject/backend
chmod +x main
nohup ./main > app.log 2>&1 &
```

**更好的方式：用 systemd 管理**

```bash
sudo nano /etc/systemd/system/myproject.service
```

```ini
[Unit]
Description=My Go Project
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/myproject/backend
ExecStart=/var/www/myproject/backend/main
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl start myproject
sudo systemctl enable myproject
```

### 9.4 配置 Nginx

```bash
sudo nano /etc/nginx/sites-available/myproject
```

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端
    root /var/www/myproject/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/myproject /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 常见问题

### 1. 403 Forbidden

```bash
# 检查目录权限
sudo chown -R www-data:www-data /var/www/myproject
sudo chmod -R 755 /var/www/myproject
```

### 2. 502 Bad Gateway

```bash
# 后端服务没启动，检查后端
sudo systemctl status myproject

# 查看后端日志
tail -f /var/www/myproject/backend/app.log
```

### 3. 配置不生效

```bash
# 检查配置语法
sudo nginx -t

# 确保软链接存在
ls -la /etc/nginx/sites-enabled/

# 重新加载
sudo systemctl reload nginx
```

### 4. 端口被占用

```bash
# 查看端口占用
sudo lsof -i :80
sudo lsof -i :8080

# 杀掉进程
sudo kill -9 <PID>
```

---

## 配置模板总结

### 纯静态网站
```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/html;
    index index.html;
    
    location / {
        try_files $uri $uri/ =404;
    }
}
```

### 前后端分离
```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/dist;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 纯后端 API
```nginx
server {
    listen 80;
    server_name api.example.com;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 总结

| 步骤 | 命令/操作 |
|------|----------|
| 安装 | `apt install nginx` |
| 启动 | `systemctl start nginx` |
| 配置文件 | `/etc/nginx/sites-available/xxx` |
| 启用配置 | `ln -s` 到 `sites-enabled` |
| 测试配置 | `nginx -t` |
| 重载配置 | `systemctl reload nginx` |
| 查看日志 | `tail -f /var/log/nginx/error.log` |

**部署流程：**
1. 上传前端文件到 `/var/www/xxx/dist`
2. 上传后端可执行文件
3. 用 systemd 启动后端
4. 配置 Nginx 反向代理
5. 测试并重载 Nginx
