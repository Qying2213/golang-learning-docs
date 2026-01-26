# 练习项目：用户管理系统

> 综合练习 01.md 的所有知识点

---

## 项目功能

1. 用户注册
2. 用户登录（JWT）
3. 获取用户信息
4. 修改用户信息
5. 用户列表（分页）
6. 修改密码

---

## 技术栈

- Gin 框架
- GORM + MySQL
- Redis 缓存
- JWT 认证
- CORS 跨域

---

## 项目结构

```
test1/
├── main.go              # 入口文件
├── config/
│   └── config.go        # 配置
├── model/
│   └── user.go          # 用户模型
├── handler/
│   └── user.go          # 处理函数
├── service/
│   └── user.go          # 业务逻辑
├── middleware/
│   ├── jwt.go           # JWT 中间件
│   └── cors.go          # CORS 中间件
├── utils/
│   ├── jwt.go           # JWT 工具
│   └── response.go      # 统一响应
└── README.md
```

---

## API 设计

| 方法 | 路径 | 说明 | 是否需要登录 |
|------|------|------|-------------|
| POST | /api/register | 注册 | 否 |
| POST | /api/login | 登录 | 否 |
| GET | /api/user/profile | 获取当前用户信息 | 是 |
| PUT | /api/user/profile | 修改用户信息 | 是 |
| PUT | /api/user/password | 修改密码 | 是 |
| GET | /api/users | 用户列表（分页） | 是 |

---

## 你需要完成的任务

### 任务 1：完成 model/user.go
- 定义 User 结构体
- 包含字段：ID, Username, Password, Email, CreatedAt, UpdatedAt

### 任务 2：完成 utils/jwt.go
- 实现 GenerateToken 函数
- 实现 ParseToken 函数

### 任务 3：完成 middleware/jwt.go
- 实现 JWT 认证中间件
- 从 Header 获取 Token
- 验证 Token 并把用户信息存到 Context

### 任务 4：完成 handler/user.go
- 实现 Register 处理函数
- 实现 Login 处理函数
- 实现 GetProfile 处理函数
- 实现 UpdateProfile 处理函数
- 实现 UpdatePassword 处理函数
- 实现 ListUsers 处理函数

### 任务 5：完成 main.go
- 初始化数据库连接
- 注册路由
- 启动服务器

---

## 运行项目

```bash
# 1. 安装依赖
go mod tidy

# 2. 创建数据库
mysql -u root -p
CREATE DATABASE test1;

# 3. 运行
go run main.go
```

---

## 测试 API

```bash
# 注册
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"秦阳","password":"123456","email":"qinyang@test.com"}'

# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"秦阳","password":"123456"}'

# 获取用户信息（需要 Token）
curl http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer <your-token>"

# 用户列表
curl "http://localhost:8080/api/users?page=1&limit=10" \
  -H "Authorization: Bearer <your-token>"
```

---

## 知识点覆盖

| 知识点 | 在项目中的应用 |
|--------|---------------|
| 结构体 | User 模型定义 |
| 接口 | error 接口、自定义错误 |
| 并发 | 数据库连接池 |
| Context | 请求上下文传递用户信息 |
| 错误处理 | 统一错误响应 |
| HTTP/RESTful | API 设计 |
| Gin | 路由、中间件、参数绑定 |
| JWT | 用户认证 |
| CORS | 跨域处理 |
| MySQL/GORM | 数据持久化 |

---

## 提示

1. 先看懂骨架代码的结构
2. 按任务顺序完成
3. 每完成一个任务就测试一下
4. 遇到问题先自己想，想不出来再问我
