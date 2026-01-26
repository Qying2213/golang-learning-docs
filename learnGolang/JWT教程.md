# JWT 教程

---

## 1. JWT 是什么

JWT = JSON Web Token

**简单说：一个加密的字符串，用来证明"你是谁"**

传统方式用 Session：
```
用户登录 → 服务器存 Session → 返回 SessionID → 浏览器存 Cookie
每次请求带 Cookie → 服务器查 Session → 确认身份
```

JWT 方式：
```
用户登录 → 服务器生成 Token → 返回给前端
每次请求带 Token → 服务器验证 Token → 确认身份
```

**区别**：Session 存在服务器，JWT 存在客户端。

---

## 2. JWT 长什么样

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoi56em6ZizIiwiZXhwIjoxNzAzNDkwMDAwfQ.abc123signature
```

三部分，用 `.` 分隔：

```
Header.Payload.Signature
头部.载荷.签名
```

### 2.1 Header（头部）

```json
{
  "alg": "HS256",  // 签名算法
  "typ": "JWT"     // 类型
}
```

Base64 编码后：`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`

### 2.2 Payload（载荷）

```json
{
  "user_id": 123,
  "username": "秦阳",
  "exp": 1703490000  // 过期时间
}
```

Base64 编码后：`eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoi56em6ZizIiwiZXhwIjoxNzAzNDkwMDAwfQ`

**注意**：Payload 只是 Base64 编码，不是加密！别放敏感信息（密码等）。

### 2.3 Signature（签名）

```
HMACSHA256(
  base64(header) + "." + base64(payload),
  secret  // 密钥，只有服务器知道
)
```

签名用来验证 Token 没被篡改。

---

## 3. JWT 工作流程

```
1. 用户登录
   POST /login {username, password}
   
2. 服务器验证密码，生成 JWT
   返回: {token: "eyJhbG..."}
   
3. 前端保存 Token（localStorage 或内存）

4. 之后每次请求带上 Token
   Authorization: Bearer eyJhbG...
   
5. 服务器验证 Token
   - 验证签名（没被篡改）
   - 验证过期时间
   - 从 Payload 取出用户信息
```

---

## 4. Go 实现 JWT

### 4.1 安装库

```bash
go get github.com/golang-jwt/jwt/v5
```

### 4.2 生成 Token

```go
package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// 密钥（实际项目从配置读取，不要硬编码）
var jwtSecret = []byte("your-secret-key")

// 自定义 Claims
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// 生成 Token
func GenerateToken(userID int, username string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "my-app",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func main() {
    token, _ := GenerateToken(123, "秦阳")
    fmt.Println(token)
}
```

### 4.3 验证 Token

```go
// 验证并解析 Token
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}

func main() {
    tokenString := "eyJhbG..."
    
    claims, err := ParseToken(tokenString)
    if err != nil {
        fmt.Println("Token 无效:", err)
        return
    }
    
    fmt.Println("用户ID:", claims.UserID)
    fmt.Println("用户名:", claims.Username)
}
```

---

## 5. Gin 框架中使用 JWT

### 5.1 完整示例

```go
package main

import (
    "net/http"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

// Claims
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// 生成 Token
func GenerateToken(userID int, username string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, fmt.Errorf("invalid token")
}

// JWT 中间件
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Header 获取 Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 Token"})
            c.Abort()
            return
        }
        
        // 格式：Bearer <token>
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误"})
            c.Abort()
            return
        }
        
        // 验证 Token
        claims, err := ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效"})
            c.Abort()
            return
        }
        
        // 把用户信息存到 Context，后续 handler 可以用
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        
        c.Next()
    }
}

func main() {
    r := gin.Default()
    
    // 登录接口（不需要 Token）
    r.POST("/login", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        c.BindJSON(&req)
        
        // 验证用户名密码（这里简化了）
        if req.Username == "秦阳" && req.Password == "123456" {
            token, _ := GenerateToken(1, req.Username)
            c.JSON(200, gin.H{"token": token})
        } else {
            c.JSON(401, gin.H{"error": "用户名或密码错误"})
        }
    })
    
    // 需要登录的接口
    auth := r.Group("/api")
    auth.Use(JWTMiddleware())
    {
        auth.GET("/profile", func(c *gin.Context) {
            userID := c.GetInt("userID")
            username := c.GetString("username")
            c.JSON(200, gin.H{
                "user_id":  userID,
                "username": username,
            })
        })
        
        auth.GET("/orders", func(c *gin.Context) {
            userID := c.GetInt("userID")
            c.JSON(200, gin.H{
                "user_id": userID,
                "orders":  []string{"订单1", "订单2"},
            })
        })
    }
    
    r.Run(":8080")
}
```

### 5.2 测试

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"秦阳","password":"123456"}'

# 返回: {"token":"eyJhbG..."}

# 用 Token 访问接口
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer eyJhbG..."

# 返回: {"user_id":1,"username":"秦阳"}
```

---

## 6. Token 刷新

Token 过期了怎么办？两种方案：

### 6.1 方案一：双 Token

```
Access Token:  有效期短（如 15 分钟）
Refresh Token: 有效期长（如 7 天）
```

```go
// 登录时返回两个 Token
func Login() {
    accessToken := GenerateToken(userID, 15*time.Minute)
    refreshToken := GenerateToken(userID, 7*24*time.Hour)
    
    return gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    }
}

// 刷新接口
r.POST("/refresh", func(c *gin.Context) {
    refreshToken := c.PostForm("refresh_token")
    
    claims, err := ParseToken(refreshToken)
    if err != nil {
        c.JSON(401, gin.H{"error": "Refresh Token 无效"})
        return
    }
    
    // 生成新的 Access Token
    newAccessToken := GenerateToken(claims.UserID, 15*time.Minute)
    c.JSON(200, gin.H{"access_token": newAccessToken})
})
```

### 6.2 方案二：自动续期

```go
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... 验证 Token ...
        
        // 如果快过期了（比如还剩 1 小时），自动续期
        if claims.ExpiresAt.Time.Sub(time.Now()) < time.Hour {
            newToken, _ := GenerateToken(claims.UserID, claims.Username)
            c.Header("X-New-Token", newToken)
        }
        
        c.Next()
    }
}
```

---

## 7. 安全注意事项

### 7.1 密钥要保密

```go
// ❌ 不要硬编码
var jwtSecret = []byte("123456")

// ✅ 从环境变量读取
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

### 7.2 使用 HTTPS

Token 在网络传输中可能被截获，必须用 HTTPS。

### 7.3 设置合理的过期时间

```go
// Access Token 不要太长
ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
```

### 7.4 不要在 Payload 放敏感信息

```go
// ❌ 不要放密码
claims := Claims{
    Password: "123456",  // 危险！
}

// ✅ 只放必要信息
claims := Claims{
    UserID:   123,
    Username: "秦阳",
}
```

### 7.5 Token 存储

前端存储 Token 的位置：

| 位置 | 优点 | 缺点 |
|------|------|------|
| localStorage | 方便 | XSS 攻击可以读取 |
| Cookie (HttpOnly) | XSS 无法读取 | CSRF 攻击风险 |
| 内存 | 最安全 | 刷新页面丢失 |

---

## 8. JWT vs Session

| 对比 | JWT | Session |
|------|-----|---------|
| 存储位置 | 客户端 | 服务器 |
| 扩展性 | 好（无状态） | 差（需要共享 Session） |
| 安全性 | 无法主动失效 | 可以随时删除 |
| 性能 | 不需要查数据库 | 需要查 Session 存储 |
| 适用场景 | 分布式系统、移动端 | 传统 Web 应用 |

---

## 9. JWT 和 CORS 的关系

JWT 和 CORS 是两个不同的东西，但经常一起用：

- **CORS**：解决跨域问题（浏览器允不允许请求）
- **JWT**：解决身份认证问题（你是谁）

### 9.1 为什么要一起用

```
前端: http://localhost:3000
后端: http://localhost:8080

1. 前端请求后端 → 跨域了 → 需要 CORS
2. 前端带 Token 请求 → Authorization 是自定义头 → 触发预检请求
3. 后端必须允许 Authorization 头 → CORS 配置要加上
```

### 9.2 CORS 配置要点

```go
// CORS 必须允许 Authorization 头，否则 JWT 发不过去
cors.Config{
    AllowHeaders: []string{"Content-Type", "Authorization"},  // 重要！
}
```

### 9.3 中间件顺序

**CORS 要在 JWT 前面**：

```go
r := gin.Default()

// 1. 先 CORS（处理预检请求）
r.Use(cors.Default())

// 2. 再 JWT（验证身份）
r.Use(JWTMiddleware())
```

为什么？
- 预检请求（OPTIONS）不带 Token
- 如果 JWT 在前面，预检请求会被拦截返回 401
- CORS 在前面可以先处理 OPTIONS，直接返回 200

### 9.4 完整示例：CORS + JWT

```go
package main

import (
    "net/http"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, fmt.Errorf("invalid token")
}

// JWT 中间件
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 Token"})
            c.Abort()
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误"})
            c.Abort()
            return
        }
        
        claims, err := ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效"})
            c.Abort()
            return
        }
        
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}

func main() {
    r := gin.Default()
    
    // ========== 1. CORS 中间件（全局，最先执行）==========
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},  // 必须包含 Authorization
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    // ========== 2. 公开接口（不需要 Token）==========
    r.POST("/login", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        c.BindJSON(&req)
        
        if req.Username == "秦阳" && req.Password == "123456" {
            token, _ := GenerateToken(1, req.Username)
            c.JSON(200, gin.H{"token": token})
        } else {
            c.JSON(401, gin.H{"error": "用户名或密码错误"})
        }
    })
    
    r.POST("/register", func(c *gin.Context) {
        // 注册逻辑...
        c.JSON(200, gin.H{"message": "注册成功"})
    })
    
    // ========== 3. 需要认证的接口（需要 Token）==========
    auth := r.Group("/api")
    auth.Use(JWTMiddleware())  // 只对这个组启用 JWT
    {
        auth.GET("/profile", func(c *gin.Context) {
            userID := c.GetInt("userID")
            username := c.GetString("username")
            c.JSON(200, gin.H{
                "user_id":  userID,
                "username": username,
            })
        })
        
        auth.GET("/orders", func(c *gin.Context) {
            c.JSON(200, gin.H{"orders": []string{"订单1", "订单2"}})
        })
    }
    
    r.Run(":8080")
}
```

### 9.5 前端请求示例

```javascript
// 登录
const res = await fetch('http://localhost:8080/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: '秦阳', password: '123456' })
})
const { token } = await res.json()

// 保存 Token
localStorage.setItem('token', token)

// 带 Token 请求
const profile = await fetch('http://localhost:8080/api/profile', {
    headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
})
```

### 9.6 常见问题

**问题 1：预检请求返回 401**
```
原因：JWT 中间件拦截了 OPTIONS 请求
解决：CORS 中间件放在 JWT 前面
```

**问题 2：Authorization 头被拦截**
```
原因：CORS 没有允许 Authorization 头
解决：AllowHeaders 加上 "Authorization"
```

**问题 3：Token 发不出去**
```
原因：跨域请求，浏览器不让发自定义头
解决：后端配置 CORS，允许 Authorization 头
```

---

## 10. 练习

### 练习 1：基本使用

```go
// 实现一个简单的登录系统
// 1. POST /register 注册用户
// 2. POST /login 登录返回 Token
// 3. GET /api/me 获取当前用户信息（需要 Token）
```

### 练习 2：Token 刷新

```go
// 在练习 1 基础上
// 1. 登录返回 access_token 和 refresh_token
// 2. POST /refresh 用 refresh_token 换新的 access_token
```

---

## 总结

| 概念 | 说明 |
|------|------|
| JWT | 三部分：Header.Payload.Signature |
| Header | 算法和类型 |
| Payload | 用户数据（不加密，别放敏感信息） |
| Signature | 签名，防篡改 |
| 生成 | `jwt.NewWithClaims()` + `SignedString()` |
| 验证 | `jwt.ParseWithClaims()` |

**核心流程**：
1. 登录 → 生成 Token → 返回给前端
2. 前端请求带 `Authorization: Bearer <token>`
3. 后端验证 Token → 从 Payload 取用户信息
