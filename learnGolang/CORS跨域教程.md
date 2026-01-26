# CORS 跨域教程

---

## 1. 什么是跨域

**同源策略**：浏览器的安全机制，只允许网页请求同源的资源。

**同源 = 协议 + 域名 + 端口 都相同**

```
https://example.com:443/api

协议: https
域名: example.com
端口: 443
```

### 1.1 哪些是跨域

```
当前页面: https://example.com

✅ 同源: https://example.com/api
❌ 跨域: http://example.com/api      (协议不同)
❌ 跨域: https://api.example.com     (域名不同)
❌ 跨域: https://example.com:8080    (端口不同)
```

### 1.2 跨域会怎样

```
前端: http://localhost:3000
后端: http://localhost:8080

前端请求后端 → 浏览器拦截 → 报错：
Access to fetch at 'http://localhost:8080/api' from origin 
'http://localhost:3000' has been blocked by CORS policy
```

**注意**：跨域是浏览器的限制，后端其实收到请求了，响应也发了，但浏览器不让前端拿到。

---

## 2. CORS 是什么

CORS = Cross-Origin Resource Sharing（跨域资源共享）

**简单说：后端告诉浏览器"这个域名可以访问我"**

```
后端响应头:
Access-Control-Allow-Origin: http://localhost:3000

浏览器看到这个头: "哦，后端允许，那我放行"
```

---

## 3. 简单请求 vs 预检请求

### 3.1 简单请求

满足以下条件的是简单请求：
- 方法：GET、POST、HEAD
- 请求头：只有基本头（Accept、Content-Type 等）
- Content-Type：只能是这三种
  - text/plain
  - multipart/form-data
  - application/x-www-form-urlencoded

简单请求直接发，浏览器检查响应头。

### 3.2 预检请求（Preflight）

不满足简单请求条件的，浏览器会先发一个 OPTIONS 请求"问一下"：

```
1. 浏览器: OPTIONS /api/users (我能用 PUT 方法吗？)
2. 服务器: 可以，允许 PUT
3. 浏览器: PUT /api/users (真正的请求)
4. 服务器: 返回数据
```

**什么情况会触发预检**：
- 方法是 PUT、DELETE、PATCH
- Content-Type 是 application/json
- 有自定义请求头（如 Authorization）

---

## 4. CORS 响应头

### 4.1 必须的头

```
Access-Control-Allow-Origin: http://localhost:3000
```
- 允许哪个域名访问
- `*` 表示允许所有（不安全，不能带 Cookie）

### 4.2 预检请求相关

```
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```
- Allow-Methods: 允许的 HTTP 方法
- Allow-Headers: 允许的请求头
- Max-Age: 预检结果缓存时间（秒），避免每次都预检

### 4.3 携带 Cookie

```
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://localhost:3000  # 不能用 *
```
- 允许携带 Cookie
- 此时 Origin 不能是 `*`，必须指定具体域名

---

## 5. Go 实现 CORS

### 5.1 简单版

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 设置 CORS 头
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // 处理预检请求
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello"))
    })
    
    http.Handle("/", CORSMiddleware(handler))
    http.ListenAndServe(":8080", nil)
}
```

### 5.2 完整版（支持 Cookie）

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // 允许的域名列表
        allowedOrigins := map[string]bool{
            "http://localhost:3000": true,
            "https://example.com":   true,
        }
        
        // 检查是否允许
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Max-Age", "86400")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### 5.3 使用 Gin 框架

#### 方式一：手写中间件

```go
package main

import "github.com/gin-gonic/gin"

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Max-Age", "86400")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }
        
        c.Next()
    }
}

func main() {
    r := gin.Default()
    
    r.Use(CORSMiddleware())
    
    r.GET("/api/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"name": "秦阳"})
    })
    
    r.Run(":8080")
}
```

#### 方式二：使用 gin-contrib/cors 库（推荐）

安装：
```bash
go get github.com/gin-contrib/cors
```

**最简单用法**：
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    
    // 允许所有跨域
    r.Use(cors.Default())
    
    r.GET("/api/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"name": "秦阳"})
    })
    
    r.Run(":8080")
}
```

**自定义配置**：
```go
package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "https://example.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    r.GET("/api/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"name": "秦阳"})
    })
    
    r.Run(":8080")
}
```

**允许所有域名（开发环境）**：
```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

**只对某些路由启用 CORS**：
```go
func main() {
    r := gin.Default()
    
    corsConfig := cors.New(cors.Config{
        AllowAllOrigins: true,
        AllowMethods:    []string{"GET", "POST"},
    })
    
    // 只有 /api 路由组启用 CORS
    api := r.Group("/api")
    api.Use(corsConfig)
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
    }
    
    // 其他路由不启用
    r.GET("/health", healthCheck)
    
    r.Run(":8080")
}
```

#### cors.Config 配置说明

```go
cors.Config{
    // 允许的域名列表
    AllowOrigins: []string{"http://localhost:3000"},
    
    // 允许所有域名
    AllowAllOrigins: true,
    
    // 自定义判断函数
    AllowOriginFunc: func(origin string) bool {
        return strings.HasSuffix(origin, ".example.com")
    },
    
    // 允许的方法
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    
    // 允许的请求头
    AllowHeaders: []string{"Content-Type", "Authorization"},
    
    // 暴露给前端的响应头
    ExposeHeaders: []string{"X-Custom-Header"},
    
    // 允许携带 Cookie
    AllowCredentials: true,
    
    // 预检缓存时间
    MaxAge: 12 * time.Hour,
}
```

### 5.4 使用 rs/cors 库（推荐）

先安装：
```bash
go get github.com/rs/cors
```

#### 方式一：包装整个 Handler

```go
package main

import (
    "net/http"
    "github.com/rs/cors"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"name":"秦阳"}`))
    })
    
    // 创建 CORS 中间件
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           86400,
    })
    
    // 用 c.Handler() 包装
    http.ListenAndServe(":8080", c.Handler(mux))
}
```

#### 方式二：作为中间件链的一环

```go
package main

import (
    "log"
    "net/http"
    "time"
    "github.com/rs/cors"
)

// 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"name":"秦阳"}`))
    })
    
    // CORS 中间件
    corsMiddleware := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
    })
    
    // 中间件链：Logging -> CORS -> Handler
    handler := LoggingMiddleware(corsMiddleware.Handler(mux))
    
    http.ListenAndServe(":8080", handler)
}
```

#### 方式三：允许所有（开发环境用）

```go
// 最简单，允许所有跨域请求
c := cors.AllowAll()
http.ListenAndServe(":8080", c.Handler(mux))
```

#### 方式四：默认配置

```go
// 默认配置：允许所有 Origin，只允许 GET/POST
c := cors.Default()
http.ListenAndServe(":8080", c.Handler(mux))
```

#### cors.Options 配置说明

```go
cors.Options{
    // 允许的域名，* 表示所有
    AllowedOrigins: []string{"http://localhost:3000"},
    
    // 允许的 HTTP 方法
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    
    // 允许的请求头
    AllowedHeaders: []string{"Content-Type", "Authorization"},
    
    // 暴露给前端的响应头
    ExposedHeaders: []string{"X-Custom-Header"},
    
    // 是否允许携带 Cookie
    AllowCredentials: true,
    
    // 预检请求缓存时间（秒）
    MaxAge: 86400,
    
    // 是否允许所有 Origin（设为 true 时 AllowedOrigins 无效）
    AllowOriginFunc: func(origin string) bool {
        return true  // 自定义判断逻辑
    },
    
    // 调试模式，打印 CORS 日志
    Debug: true,
}
```

---

## 6. 前端配置

### 6.1 fetch 携带 Cookie

```javascript
fetch('http://localhost:8080/api/users', {
    method: 'GET',
    credentials: 'include'  // 携带 Cookie
})
```

### 6.2 axios 携带 Cookie

```javascript
axios.defaults.withCredentials = true

axios.get('http://localhost:8080/api/users')
```

---

## 7. 常见问题

### 7.1 预检请求失败

```
Access to XMLHttpRequest has been blocked by CORS policy: 
Response to preflight request doesn't pass access control check
```

**原因**：OPTIONS 请求没有正确处理

**解决**：
```go
if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK)
    return
}
```

### 7.2 Cookie 发不出去

**检查**：
1. 后端设置了 `Access-Control-Allow-Credentials: true`
2. 后端 `Access-Control-Allow-Origin` 不是 `*`
3. 前端设置了 `credentials: 'include'`

### 7.3 自定义头被拦截

```
Request header field Authorization is not allowed
```

**解决**：后端添加
```go
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

---

## 8. CORS vs 其他跨域方案

| 方案 | 说明 | 适用场景 |
|------|------|----------|
| CORS | 后端设置响应头 | 最常用，标准方案 |
| Nginx 反向代理 | 前后端同域 | 生产环境推荐 |
| JSONP | 利用 script 标签 | 老方案，只支持 GET |
| WebSocket | 不受同源限制 | 实时通信 |

---

## 9. 完整示例

### 后端（Go）

```go
package main

import (
    "encoding/json"
    "net/http"
)

func main() {
    http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        // CORS 头
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "name": "秦阳",
        })
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### 前端（HTML）

```html
<!DOCTYPE html>
<html>
<body>
    <button onclick="fetchData()">请求数据</button>
    <div id="result"></div>
    
    <script>
        async function fetchData() {
            const resp = await fetch('http://localhost:8080/api/users')
            const data = await resp.json()
            document.getElementById('result').innerText = JSON.stringify(data)
        }
    </script>
</body>
</html>
```

---

## 总结

| 响应头 | 作用 |
|--------|------|
| Access-Control-Allow-Origin | 允许的域名 |
| Access-Control-Allow-Methods | 允许的方法 |
| Access-Control-Allow-Headers | 允许的请求头 |
| Access-Control-Allow-Credentials | 允许携带 Cookie |
| Access-Control-Max-Age | 预检缓存时间 |

**核心记住**：
1. 跨域是浏览器的限制
2. CORS 是后端告诉浏览器"我允许这个域名访问"
3. 复杂请求会先发 OPTIONS 预检
4. 生产环境推荐用 Nginx 反向代理，根本不存在跨域
