# HTTP 协议详解

> **重要程度：⭐⭐⭐ 必须精通**  
> HTTP 是 Web 开发的核心！作为后端开发，你每天都在和 HTTP 打交道！

## 📚 本章学习目标

学完本章，你将能够：

- 理解 HTTP 请求和响应的格式
- 掌握各种 HTTP 方法的使用
- 熟记常用的 HTTP 状态码
- 理解 HTTP 头部的作用
- 掌握 Cookie 和 Session 机制

---

## 1. HTTP 是什么？⭐⭐⭐

### 1.1 HTTP 的全称

```
HTTP = HyperText Transfer Protocol
       超文本传输协议

超文本 = 不仅仅是文本，还包括图片、视频、音频等
传输协议 = 规定了数据如何传输的规则
```

### 1.2 HTTP 的特点

```
1. 简单快速
   - 请求方法简单（GET、POST 等）
   - 协议简单，容易实现

2. 请求独立
   - HTTP/1.0 默认短连接（每次请求建立新连接）
   - HTTP/1.1 默认持久连接（Keep-Alive）
   - 但每个请求仍是独立的，服务器不记录请求间的关系

3. 无状态
   - 服务器不记住你是谁
   - 每次请求都是独立的
   - 需要 Cookie/Session 来保持状态

4. 基于 TCP
   - 可靠传输
   - 有序传输
```

### 1.3 HTTP 的工作流程

```
客户端（浏览器/App）          服务器
     │                          │
     │  1. 建立 TCP 连接         │
     ├─────────────────────────→│
     │                          │
     │  2. 发送 HTTP 请求        │
     ├─────────────────────────→│
     │                          │
     │  3. 服务器处理请求        │
     │                          │
     │  4. 返回 HTTP 响应        │
     │←─────────────────────────┤
     │                          │
     │  5. 断开 TCP 连接         │
     │←─────────────────────────┤
     │                          │
```

---

## 2. HTTP 请求 ⭐⭐⭐

### 2.1 HTTP 请求的格式

```
GET /api/users HTTP/1.1              ← 请求行
Host: api.example.com                ← 请求头
User-Agent: Mozilla/5.0              ← 请求头
Accept: application/json             ← 请求头
                                     ← 空行
{"name": "张三", "age": 18}          ← 请求体（可选）

分为四部分：
1. 请求行：方法 + 路径 + 协议版本
2. 请求头：键值对，提供额外信息
3. 空行：分隔头部和主体
4. 请求体：实际数据（GET 请求通常没有）
```

### 2.2 请求行详解

```
GET /api/users?page=1&size=10 HTTP/1.1
│   │                         │
│   │                         └─ 协议版本
│   └─ 请求路径（包含查询参数）
└─ 请求方法

常见方法：
GET    - 获取资源
POST   - 创建资源
PUT    - 更新资源（完整更新）
PATCH  - 更新资源（部分更新）
DELETE - 删除资源
HEAD   - 获取响应头（不要响应体）
OPTIONS - 查询支持的方法
```

### 2.3 常用请求头

```
Host: api.example.com
- 目标服务器的域名
- 必需的请求头

User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)
- 客户端信息（浏览器、操作系统）
- 服务器可以根据这个返回不同内容

Accept: application/json
- 客户端能接受的数据类型
- application/json、text/html、image/png 等

Content-Type: application/json
- 请求体的数据类型
- 告诉服务器如何解析请求体

Content-Length: 123
- 请求体的长度（字节）

Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
- 认证信息（Token）

Cookie: session_id=abc123; user_id=456
- 客户端存储的数据

Referer: https://www.example.com/page1
- 从哪个页面跳转过来的

Connection: keep-alive
- 是否保持连接
```

### 📝 练习题 2.1

**问题**：写出一个完整的 HTTP GET 请求，获取用户列表。

<details>
<summary>点击查看答案</summary>

```http
GET /api/users?page=1&size=20 HTTP/1.1
Host: api.example.com
User-Agent: Mozilla/5.0
Accept: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Connection: keep-alive

```

注意：

- GET 请求没有请求体
- 参数放在 URL 的查询字符串中
- 最后有一个空行

</details>

---

## 3. HTTP 响应 ⭐⭐⭐

### 3.1 HTTP 响应的格式

```
HTTP/1.1 200 OK                      ← 状态行
Content-Type: application/json       ← 响应头
Content-Length: 58                   ← 响应头
Date: Mon, 19 Jan 2026 10:00:00 GMT  ← 响应头
                                     ← 空行
{"id": 1, "name": "张三", "age": 18} ← 响应体

分为四部分：
1. 状态行：协议版本 + 状态码 + 状态描述
2. 响应头：键值对
3. 空行
4. 响应体：实际数据
```

### 3.2 状态行详解

```
HTTP/1.1 200 OK
│        │   │
│        │   └─ 状态描述（OK、Not Found 等）
│        └─ 状态码（200、404 等）
└─ 协议版本
```

### 3.3 HTTP 状态码 ⭐⭐⭐

#### 1xx - 信息响应

```
100 Continue
- 客户端应继续请求
- 很少用到
```

#### 2xx - 成功

```
200 OK
- 请求成功
- 最常见的状态码

201 Created
- 资源创建成功
- POST 请求常用

204 No Content
- 请求成功，但没有返回内容
- DELETE 请求常用
```

#### 3xx - 重定向

```
301 Moved Permanently
- 永久重定向
- 资源永久移动到新位置

302 Found
- 临时重定向
- 资源临时移动

304 Not Modified
- 资源未修改
- 可以使用缓存
```

#### 4xx - 客户端错误

```
400 Bad Request
- 请求格式错误
- 参数错误

401 Unauthorized
- 未认证
- 需要登录

403 Forbidden
- 已认证，但没有权限
- 禁止访问

404 Not Found
- 资源不存在
- 最常见的错误

405 Method Not Allowed
- 请求方法不支持
- 例如：用 GET 访问只支持 POST 的接口

409 Conflict
- 请求冲突
- 例如：用户名已存在

429 Too Many Requests
- 请求过多
- 触发限流
```

#### 5xx - 服务器错误

```
500 Internal Server Error
- 服务器内部错误
- 代码出 bug 了

502 Bad Gateway
- 网关错误
- 上游服务器无响应

503 Service Unavailable
- 服务不可用
- 服务器过载或维护中

504 Gateway Timeout
- 网关超时
- 上游服务器响应超时
```

### 📝 练习题 3.1

**问题**：以下场景应该返回什么状态码？

1. 用户登录成功
2. 用户未登录访问需要认证的接口
3. 删除用户成功
4. 创建用户时，用户名已存在
5. 服务器代码出现空指针异常

<details>
<summary>点击查看答案</summary>

```
1. 200 OK
   - 登录成功，返回 Token

2. 401 Unauthorized
   - 未认证，需要登录

3. 204 No Content 或 200 OK
   - 删除成功，通常不返回内容

4. 409 Conflict 或 400 Bad Request
   - 资源冲突

5. 500 Internal Server Error
   - 服务器内部错误
```

</details>

---

## 4. HTTP 方法详解 ⭐⭐⭐

### 4.1 GET - 获取资源

```http
GET /api/users/123 HTTP/1.1
Host: api.example.com

特点：
✅ 幂等：多次请求结果相同
✅ 安全：不会修改服务器数据
✅ 可缓存
❌ 参数在 URL 中，有长度限制
❌ 参数可见，不适合敏感数据

使用场景：
- 获取用户信息
- 获取文章列表
- 搜索
```

```go
// Go 示例
func GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := db.GetUserByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }
    c.JSON(200, user)
}
```

### 4.2 POST - 创建资源

```http
POST /api/users HTTP/1.1
Host: api.example.com
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com",
  "age": 18
}

特点：
❌ 不幂等：多次请求会创建多个资源
❌ 不安全：会修改服务器数据
❌ 不可缓存
✅ 参数在请求体中，无长度限制
✅ 参数不可见，适合敏感数据

使用场景：
- 创建用户
- 提交表单
- 上传文件
```

```go
// Go 示例
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    if err := db.Create(&user).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(201, user)
}
```

### 4.3 PUT - 完整更新资源

```http
PUT /api/users/123 HTTP/1.1
Host: api.example.com
Content-Type: application/json

{
  "name": "李四",
  "email": "lisi@example.com",
  "age": 20
}

特点：
✅ 幂等：多次请求结果相同
❌ 不安全：会修改服务器数据
✅ 需要提供完整的资源数据

使用场景：
- 完整更新用户信息
- 替换整个资源
```

### 4.4 PATCH - 部分更新资源

```http
PATCH /api/users/123 HTTP/1.1
Host: api.example.com
Content-Type: application/json

{
  "age": 21
}

特点：
✅ 幂等：多次请求结果相同
❌ 不安全：会修改服务器数据
✅ 只需要提供要修改的字段

使用场景：
- 只更新用户的某个字段
- 部分更新
```

### 4.5 DELETE - 删除资源

```http
DELETE /api/users/123 HTTP/1.1
Host: api.example.com

特点：
✅ 幂等：多次删除结果相同（资源不存在）
❌ 不安全：会修改服务器数据

使用场景：
- 删除用户
- 删除文章
```

```go
// Go 示例
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    if err := db.Delete(&User{}, id).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to delete user"})
        return
    }
    c.JSON(204, nil)
}
```

### 4.6 方法对比

| 方法   | 幂等 | 安全 | 缓存 | 请求体 | 使用场景 |
| ------ | ---- | ---- | ---- | ------ | -------- |
| GET    | ✅   | ✅   | ✅   | ❌     | 查询     |
| POST   | ❌   | ❌   | ❌   | ✅     | 创建     |
| PUT    | ✅   | ❌   | ❌   | ✅     | 完整更新 |
| PATCH  | ✅   | ❌   | ❌   | ✅     | 部分更新 |
| DELETE | ✅   | ❌   | ❌   | ❌     | 删除     |

---

## 5. Cookie 和 Session ⭐⭐⭐

### 5.1 为什么需要 Cookie？

```
问题：HTTP 是无状态的

你登录了网站 → 刷新页面 → 服务器不记得你是谁 ❌

解决：用 Cookie 保存状态

你登录了网站 → 服务器返回 Cookie → 浏览器保存 Cookie
→ 下次请求带上 Cookie → 服务器知道你是谁 ✅
```

### 5.2 Cookie 的工作流程

```
1. 用户登录
   客户端 → 服务器：POST /login {username, password}

2. 服务器验证成功，返回 Cookie
   服务器 → 客户端：Set-Cookie: session_id=abc123

3. 浏览器保存 Cookie

4. 后续请求自动带上 Cookie
   客户端 → 服务器：GET /api/profile
                    Cookie: session_id=abc123

5. 服务器根据 Cookie 识别用户
```

### 5.3 Cookie 的属性

```http
Set-Cookie: session_id=abc123; Path=/; Domain=example.com; Max-Age=3600; HttpOnly; Secure; SameSite=Lax

session_id=abc123  ← Cookie 的值
Path=/             ← 生效路径
Domain=example.com ← 生效域名
Max-Age=3600       ← 有效期（秒）
HttpOnly           ← 禁止 JavaScript 访问（防 XSS）
Secure             ← 只在 HTTPS 下发送
SameSite=Lax       ← 跨站请求限制（防 CSRF）
```

### 5.4 Session 机制

```
Cookie 存在客户端，不安全
Session 存在服务器，更安全

工作流程：
1. 用户登录成功
2. 服务器创建 Session，生成 Session ID
3. 把 Session ID 通过 Cookie 发给客户端
4. 客户端保存 Session ID
5. 后续请求带上 Session ID
6. 服务器根据 Session ID 查找 Session 数据

┌─────────────┐                    ┌─────────────┐
│   客户端     │                    │   服务器     │
│             │                    │             │
│ Cookie:     │                    │ Session:    │
│ session_id  │ ←─────────────────→│ {user_id,   │
│ = abc123    │                    │  username,  │
│             │                    │  ...}       │
└─────────────┘                    └─────────────┘
```

### 5.5 Go 中使用 Session

```go
package main

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 使用 Cookie 存储 Session
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", store))

    // 登录接口
    r.POST("/login", func(c *gin.Context) {
        session := sessions.Default(c)

        // 验证用户名密码...

        // 保存到 Session
        session.Set("user_id", 123)
        session.Set("username", "zhangsan")
        session.Save()

        c.JSON(200, gin.H{"message": "Login success"})
    })

    // 获取用户信息
    r.GET("/profile", func(c *gin.Context) {
        session := sessions.Default(c)

        userID := session.Get("user_id")
        if userID == nil {
            c.JSON(401, gin.H{"error": "Not logged in"})
            return
        }

        c.JSON(200, gin.H{
            "user_id": userID,
            "username": session.Get("username"),
        })
    })

    r.Run(":8080")
}
```

### 📝 练习题 5.1

**问题**：Cookie 和 Session 有什么区别？

<details>
<summary>点击查看答案</summary>

| 对比项       | Cookie           | Session            |
| ------------ | ---------------- | ------------------ |
| **存储位置** | 客户端（浏览器） | 服务器             |
| **安全性**   | 较低（可被篡改） | 较高               |
| **存储大小** | 4KB 左右         | 无限制             |
| **性能**     | 不占服务器资源   | 占用服务器内存     |
| **生命周期** | 可设置过期时间   | 通常会话结束就失效 |

**使用建议**：

- 不敏感的数据：用 Cookie（如：语言偏好、主题）
- 敏感数据：用 Session（如：用户 ID、权限）
- 现代做法：用 JWT Token（下一章讲）

</details>

---

## 6. 本章总结

### 必须记住的内容

**HTTP 方法**：

- GET：查询
- POST：创建
- PUT：完整更新
- PATCH：部分更新
- DELETE：删除

**常用状态码**：

- 200：成功
- 201：创建成功
- 400：请求错误
- 401：未认证
- 403：无权限
- 404：不存在
- 500：服务器错误

**Cookie vs Session**：

- Cookie：客户端存储
- Session：服务器存储

### 下一章预告

下一章我们将学习 **HTTPS 和加密**，包括：

- 为什么需要 HTTPS
- 对称加密和非对称加密
- SSL/TLS 握手过程
- 数字证书

让你的网站更安全！🔒
