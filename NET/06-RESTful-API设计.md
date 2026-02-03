# RESTful API 设计

> **重要程度：⭐⭐⭐ 必须掌握**
> 作为后端开发，设计好的 API 是基本功！

## 📚 本章学习目标

学完本章，你将能够：
- 理解 REST 的核心原则
- 掌握 RESTful API 设计规范
- 设计出优雅、易用的 API
- 处理 API 版本控制和错误响应

---

## 1. REST 是什么？⭐⭐⭐

### 1.1 REST 的全称

```
REST = Representational State Transfer
       表现层状态转移

RESTful = 符合 REST 风格的

简单说：REST 是一种 API 设计风格
```

### 1.2 REST 的核心原则

```
1. 资源（Resource）
   - 一切皆资源
   - 用 URL 标识资源
   - /users、/articles、/orders

2. 统一接口
   - 用 HTTP 方法表示操作
   - GET = 查询
   - POST = 创建
   - PUT/PATCH = 更新
   - DELETE = 删除

3. 无状态
   - 每个请求包含所有信息
   - 服务器不保存客户端状态
   - 用 Token 代替 Session

4. 分层系统
   - 客户端不知道连接的是真实服务器还是代理
   - 可以加负载均衡、缓存等中间层
```

### 1.3 RESTful vs 非 RESTful

```
非 RESTful（传统风格）：
GET  /getUser?id=1
POST /createUser
POST /updateUser
POST /deleteUser?id=1

RESTful 风格：
GET    /users/1      获取用户
POST   /users        创建用户
PUT    /users/1      更新用户
DELETE /users/1      删除用户

区别：
- RESTful 用 URL 表示资源
- 用 HTTP 方法表示操作
- 更简洁、更规范
```

---

## 2. URL 设计规范 ⭐⭐⭐

### 2.1 基本规则

```
1. 使用名词，不用动词
   ✅ /users
   ❌ /getUsers

2. 使用复数形式
   ✅ /users
   ❌ /user

3. 使用小写字母
   ✅ /users
   ❌ /Users

4. 使用连字符 - 而不是下划线 _
   ✅ /user-profiles
   ❌ /user_profiles

5. 不要有文件扩展名
   ✅ /users/1
   ❌ /users/1.json
```

### 2.2 资源层级

```
单个资源：
GET /users/1              获取 ID 为 1 的用户

嵌套资源：
GET /users/1/orders       获取用户 1 的所有订单
GET /users/1/orders/5     获取用户 1 的订单 5

关联资源：
GET /orders/5/items       获取订单 5 的所有商品

建议：嵌套不超过 2 层
❌ /users/1/orders/5/items/3/reviews
✅ /order-items/3/reviews
```

### 2.3 查询参数

```
分页：
GET /users?page=1&size=20

排序：
GET /users?sort=created_at&order=desc

过滤：
GET /users?status=active&role=admin

搜索：
GET /users?q=张三
GET /users?search=张三

字段选择：
GET /users?fields=id,name,email

组合使用：
GET /users?status=active&page=1&size=20&sort=created_at
```

### 2.4 URL 设计示例

```
用户相关：
GET    /users              获取用户列表
GET    /users/1            获取单个用户
POST   /users              创建用户
PUT    /users/1            完整更新用户
PATCH  /users/1            部分更新用户
DELETE /users/1            删除用户

文章相关：
GET    /articles           获取文章列表
GET    /articles/1         获取单篇文章
POST   /articles           创建文章
PUT    /articles/1         更新文章
DELETE /articles/1         删除文章
GET    /articles/1/comments 获取文章评论
POST   /articles/1/comments 添加评论

特殊操作：
POST   /users/1/follow     关注用户
DELETE /users/1/follow     取消关注
POST   /articles/1/like    点赞文章
DELETE /articles/1/like    取消点赞
```

---

## 3. HTTP 方法使用 ⭐⭐⭐

### 3.1 方法对应操作

| 方法 | 操作 | 幂等 | 安全 | 示例 |
|------|------|------|------|------|
| GET | 查询 | ✅ | ✅ | 获取用户信息 |
| POST | 创建 | ❌ | ❌ | 创建新用户 |
| PUT | 完整更新 | ✅ | ❌ | 更新用户所有字段 |
| PATCH | 部分更新 | ✅ | ❌ | 只更新用户名 |
| DELETE | 删除 | ✅ | ❌ | 删除用户 |

### 3.2 GET - 查询

```go
// 获取用户列表
// GET /users?page=1&size=20
func GetUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    
    users, total := userService.List(page, size)
    
    c.JSON(200, gin.H{
        "data":  users,
        "total": total,
        "page":  page,
        "size":  size,
    })
}

// 获取单个用户
// GET /users/:id
func GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := userService.GetByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }
    c.JSON(200, user)
}
```

### 3.3 POST - 创建

```go
// 创建用户
// POST /users
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := userService.Create(req)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create user"})
        return
    }
    
    // 返回 201 Created
    c.JSON(201, user)
}
```

### 3.4 PUT vs PATCH

```go
// PUT - 完整更新（需要提供所有字段）
// PUT /users/:id
func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 完整替换
    user, err := userService.Update(id, req)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to update user"})
        return
    }
    c.JSON(200, user)
}

// PATCH - 部分更新（只提供要修改的字段）
// PATCH /users/:id
func PatchUser(c *gin.Context) {
    id := c.Param("id")
    var req map[string]interface{}
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 只更新提供的字段
    user, err := userService.Patch(id, req)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to patch user"})
        return
    }
    c.JSON(200, user)
}
```

### 3.5 DELETE - 删除

```go
// 删除用户
// DELETE /users/:id
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    
    if err := userService.Delete(id); err != nil {
        c.JSON(500, gin.H{"error": "Failed to delete user"})
        return
    }
    
    // 返回 204 No Content
    c.Status(204)
}
```

---

## 4. 状态码使用 ⭐⭐⭐

### 4.1 常用状态码

```
成功：
200 OK              - 请求成功
201 Created         - 创建成功
204 No Content      - 删除成功（无返回内容）

客户端错误：
400 Bad Request     - 请求参数错误
401 Unauthorized    - 未认证（需要登录）
403 Forbidden       - 无权限
404 Not Found       - 资源不存在
409 Conflict        - 资源冲突（如用户名已存在）
422 Unprocessable   - 参数验证失败
429 Too Many Requests - 请求过多（限流）

服务器错误：
500 Internal Error  - 服务器内部错误
502 Bad Gateway     - 网关错误
503 Unavailable     - 服务不可用
```

### 4.2 状态码使用场景

```go
// 200 - 查询成功
c.JSON(200, user)

// 201 - 创建成功
c.JSON(201, newUser)

// 204 - 删除成功
c.Status(204)

// 400 - 参数错误
c.JSON(400, gin.H{"error": "Invalid email format"})

// 401 - 未登录
c.JSON(401, gin.H{"error": "Please login first"})

// 403 - 无权限
c.JSON(403, gin.H{"error": "Permission denied"})

// 404 - 不存在
c.JSON(404, gin.H{"error": "User not found"})

// 409 - 冲突
c.JSON(409, gin.H{"error": "Username already exists"})

// 500 - 服务器错误
c.JSON(500, gin.H{"error": "Internal server error"})
```

---

## 5. 响应格式设计 ⭐⭐⭐

### 5.1 统一响应格式

```go
// 成功响应
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "张三",
        "email": "zhangsan@example.com"
    }
}

// 列表响应
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [...],
        "total": 100,
        "page": 1,
        "size": 20
    }
}

// 错误响应
{
    "code": 10001,
    "message": "用户名已存在",
    "data": null
}
```

### 5.2 Go 实现统一响应

```go
// 响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 分页响应
type PageResponse struct {
    List  interface{} `json:"list"`
    Total int64       `json:"total"`
    Page  int         `json:"page"`
    Size  int         `json:"size"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

// 分页响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, size int) {
    c.JSON(200, Response{
        Code:    0,
        Message: "success",
        Data: PageResponse{
            List:  list,
            Total: total,
            Page:  page,
            Size:  size,
        },
    })
}

// 错误响应
func Error(c *gin.Context, httpCode int, code int, message string) {
    c.JSON(httpCode, Response{
        Code:    code,
        Message: message,
        Data:    nil,
    })
}

// 使用示例
func GetUser(c *gin.Context) {
    user, err := userService.GetByID(id)
    if err != nil {
        Error(c, 404, 10001, "用户不存在")
        return
    }
    Success(c, user)
}
```

### 5.3 错误码设计

```go
// 错误码规范
// 10xxx - 用户相关
// 20xxx - 文章相关
// 30xxx - 订单相关
// 40xxx - 系统错误

const (
    // 通用错误
    ErrSuccess       = 0
    ErrUnknown       = 99999
    ErrInvalidParam  = 10000
    
    // 用户错误 10xxx
    ErrUserNotFound     = 10001
    ErrUserExists       = 10002
    ErrPasswordWrong    = 10003
    ErrTokenInvalid     = 10004
    ErrTokenExpired     = 10005
    
    // 文章错误 20xxx
    ErrArticleNotFound  = 20001
    ErrArticleExists    = 20002
)

var errMessages = map[int]string{
    ErrSuccess:          "success",
    ErrUnknown:          "未知错误",
    ErrInvalidParam:     "参数错误",
    ErrUserNotFound:     "用户不存在",
    ErrUserExists:       "用户已存在",
    ErrPasswordWrong:    "密码错误",
    ErrTokenInvalid:     "Token 无效",
    ErrTokenExpired:     "Token 已过期",
    ErrArticleNotFound:  "文章不存在",
}
```

---

## 6. API 版本控制 ⭐⭐

### 6.1 版本控制方式

```
方式1：URL 路径（推荐）
/api/v1/users
/api/v2/users

方式2：请求头
GET /api/users
Header: Api-Version: 1

方式3：查询参数
/api/users?version=1
```

### 6.2 Go 实现版本控制

```go
func main() {
    r := gin.Default()
    
    // v1 版本
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", v1GetUsers)
        v1.POST("/users", v1CreateUser)
    }
    
    // v2 版本（新功能）
    v2 := r.Group("/api/v2")
    {
        v2.GET("/users", v2GetUsers)  // 返回格式可能不同
        v2.POST("/users", v2CreateUser)
    }
    
    r.Run(":8080")
}
```

### 6.3 版本迁移策略

```
1. 新版本发布时，旧版本继续可用
2. 给用户足够的迁移时间
3. 旧版本标记为 deprecated
4. 最终下线旧版本

时间线示例：
- v1 发布
- v2 发布，v1 标记 deprecated
- 6 个月后，v1 下线
```

---

## 7. 认证与授权 ⭐⭐⭐

### 7.1 常用认证方式

```
1. JWT Token（推荐）
   Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

2. API Key
   X-API-Key: your-api-key

3. Basic Auth
   Authorization: Basic base64(username:password)

4. OAuth 2.0
   用于第三方授权
```

### 7.2 JWT 认证示例

```go
// 登录获取 Token
// POST /auth/login
func Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)
    
    user, err := userService.Verify(req.Username, req.Password)
    if err != nil {
        Error(c, 401, ErrPasswordWrong, "用户名或密码错误")
        return
    }
    
    token, _ := jwt.GenerateToken(user.ID)
    Success(c, gin.H{"token": token})
}

// JWT 中间件
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            Error(c, 401, ErrTokenInvalid, "请先登录")
            c.Abort()
            return
        }
        
        // 去掉 "Bearer " 前缀
        token = strings.TrimPrefix(token, "Bearer ")
        
        claims, err := jwt.ParseToken(token)
        if err != nil {
            Error(c, 401, ErrTokenInvalid, "Token 无效")
            c.Abort()
            return
        }
        
        c.Set("userID", claims.UserID)
        c.Next()
    }
}

// 使用中间件
func main() {
    r := gin.Default()
    
    // 公开接口
    r.POST("/auth/login", Login)
    r.POST("/auth/register", Register)
    
    // 需要认证的接口
    auth := r.Group("/api")
    auth.Use(JWTAuth())
    {
        auth.GET("/users/me", GetCurrentUser)
        auth.PUT("/users/me", UpdateCurrentUser)
    }
    
    r.Run(":8080")
}
```

---

## 8. 最佳实践 ⭐⭐⭐

### 8.1 设计原则

```
1. 保持简单
   - URL 简洁明了
   - 参数命名清晰

2. 保持一致
   - 命名风格统一
   - 响应格式统一
   - 错误处理统一

3. 面向资源
   - 用名词不用动词
   - 用 HTTP 方法表示操作

4. 合理使用状态码
   - 不要所有都返回 200
   - 用状态码表示结果

5. 提供文档
   - Swagger/OpenAPI
   - 示例代码
```

### 8.2 安全建议

```
1. 使用 HTTPS
2. 验证所有输入
3. 限流防刷
4. 敏感数据加密
5. 日志记录
6. 权限控制
```

### 8.3 性能优化

```
1. 分页查询
2. 字段过滤
3. 缓存响应
4. 压缩传输
5. 连接池
```

---

## 9. 本章总结

### 必须记住的内容

```
REST 核心：
- 资源用 URL 表示
- 操作用 HTTP 方法表示
- 无状态

URL 设计：
- 用名词复数
- 小写字母
- 连字符分隔

HTTP 方法：
- GET = 查询
- POST = 创建
- PUT = 完整更新
- PATCH = 部分更新
- DELETE = 删除

状态码：
- 2xx = 成功
- 4xx = 客户端错误
- 5xx = 服务器错误
```

### 面试常见问题

```
Q: RESTful API 的特点？
A: 资源导向、统一接口、无状态、分层系统。

Q: PUT 和 PATCH 的区别？
A: PUT 完整更新，需要提供所有字段；PATCH 部分更新，只提供要修改的字段。

Q: 如何设计 API 版本控制？
A: 推荐 URL 路径方式，如 /api/v1/users。

Q: 401 和 403 的区别？
A: 401 未认证（没登录），403 无权限（登录了但没权限）。
```

---

下一章：[07-WebSocket实时通信](./07-WebSocket实时通信.md)
