# Task Manager - 全栈任务管理系统

一个使用 Go + React + TypeScript 构建的现代化全栈任务管理应用。

## 技术栈

### 后端
- **Go** - 高性能后端语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 开发数据库（可切换到 PostgreSQL）
- **JWT** - 用户认证
- **bcrypt** - 密码加密

### 前端
- **React 18** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 快速构建工具
- **TailwindCSS** - 样式框架
- **React Router** - 路由管理
- **Axios** - HTTP 客户端

## 功能特性

- 🔐 用户注册和登录（JWT 认证）
- ✅ 任务 CRUD（创建、查看、更新、删除）
- 📋 任务状态管理（待处理、进行中、已完成）
- 🎯 任务优先级（低、中、高）
- 📅 任务截止日期
- 🔍 任务搜索和筛选
- 📄 分页支持
- 📱 响应式设计

## 项目结构

```
project1/
├── backend/                 # Go 后端
│   ├── main.go             # 入口文件
│   ├── config/             # 配置
│   ├── database/           # 数据库连接
│   ├── handlers/           # HTTP 处理器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   └── utils/              # 工具函数
├── frontend/               # React 前端
│   ├── src/
│   │   ├── components/     # UI 组件
│   │   ├── contexts/       # React Context
│   │   ├── hooks/          # 自定义 Hooks
│   │   ├── pages/          # 页面组件
│   │   ├── services/       # API 服务
│   │   └── types/          # TypeScript 类型
│   ├── package.json
│   └── vite.config.ts
└── README.md
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- pnpm (推荐) 或 npm

### 1. 启动后端

```bash
# 进入后端目录
cd project1/backend

# 安装依赖
go mod tidy

# 运行服务器
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 2. 启动前端

```bash
# 进入前端目录
cd project1/frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

前端将在 `http://localhost:5173` 启动。

## API 文档

### 认证相关

#### 注册
```
POST /api/auth/register
Content-Type: application/json

{
  "name": "用户名",
  "email": "user@example.com",
  "password": "password123"
}
```

#### 登录
```
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### 任务相关（需要认证）

所有任务相关 API 需要在请求头中添加 JWT Token：
```
Authorization: Bearer <token>
```

#### 获取任务列表
```
GET /api/tasks?page=1&page_size=10&status=pending&priority=high&search=关键词
```

#### 创建任务
```
POST /api/tasks
Content-Type: application/json

{
  "title": "任务标题",
  "description": "任务描述",
  "priority": "high",
  "due_date": "2024-12-31T00:00:00Z"
}
```

#### 获取单个任务
```
GET /api/tasks/:id
```

#### 更新任务
```
PUT /api/tasks/:id
Content-Type: application/json

{
  "title": "新标题",
  "status": "completed"
}
```

#### 删除任务
```
DELETE /api/tasks/:id
```

## 环境变量

### 后端 (可选)
```bash
PORT=8080                                    # 服务端口
DATABASE_DSN=./taskmanager.db                # 数据库路径
JWT_SECRET=your-secret-key                   # JWT 密钥（生产环境请更改）
```

### 前端
在 `frontend/.env` 文件中配置：
```bash
VITE_API_URL=http://localhost:8080/api       # 后端 API 地址
```

## 开发指南

### 代码结构说明

#### 后端
- `handlers/` - HTTP 请求处理器，处理业务逻辑
- `middleware/` - 中间件，包括 JWT 认证和 CORS
- `models/` - 数据模型和 DTO
- `database/` - 数据库连接和迁移
- `utils/` - 工具函数（JWT 生成/验证、响应格式化）

#### 前端
- `components/` - 可复用 UI 组件
- `pages/` - 路由页面组件
- `services/` - API 调用封装
- `hooks/` - 自定义 React Hooks
- `contexts/` - React Context（全局状态管理）
- `types/` - TypeScript 类型定义

### 添加新功能

1. **后端**：在 `handlers/` 添加处理器，在 `main.go` 注册路由
2. **前端**：在 `services/` 添加 API 调用，创建对应组件

## 部署

### 后端部署

```bash
cd backend
go build -o taskmanager
./taskmanager
```

### 前端构建

```bash
cd frontend
pnpm build
# 构建产物在 dist/ 目录
```

## License

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！
