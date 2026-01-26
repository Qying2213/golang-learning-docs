# HTTP 教程（三）：实战进阶

---

## 1. HTTP 客户端进阶

### 1.1 自定义请求

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
    // 创建请求体
    data := map[string]interface{}{
        "name": "秦阳",
        "age":  22,
    }
    jsonData, _ := json.Marshal(data)
    
    // 创建请求
    req, _ := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonData))
    
    // 设置请求头
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer your-token")
    req.Header.Set("User-Agent", "MyApp/1.0")
    
    // 发送请求
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // 读取响应
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 1.2 设置超时

```go
client := &http.Client{
    Timeout: 10 * time.Second,  // 整体超时
}

// 或者更细粒度的控制
transport := &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   5 * time.Second,  // 连接超时
        KeepAlive: 30 * time.Second,
    }).DialContext,
    ResponseHeaderTimeout: 5 * time.Second,  // 响应头超时
    IdleConnTimeout:       90 * time.Second, // 空闲连接超时
}

client := &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}
```

### 1.3 连接池复用

```go
// 默认的 http.Client 会复用连接
// 但要注意读完响应体并关闭

resp, _ := http.Get(url)
io.Copy(io.Discard, resp.Body)  // 读完响应体
resp.Body.Close()                // 关闭，连接才能复用
```

---

## 2. HTTP 服务器进阶

### 2.1 中间件模式

```go
// 中间件：在处理请求前后做一些事情
type Middleware func(http.Handler) http.Handler

// 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 调用下一个处理器
        next.ServeHTTP(w, r)
        
        // 记录日志
        fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
    })
}

// 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // 验证 token...
        
        next.ServeHTTP(w, r)
    })
}

// 使用中间件
func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello"))
    })
    
    // 套娃：Auth -> Logging -> Handler
    wrapped := LoggingMiddleware(AuthMiddleware(handler))
    
    http.Handle("/", wrapped)
    http.ListenAndServe(":8080", nil)
}
```

### 2.2 路由分组

```go
// 简单的路由器
type Router struct {
    mux *http.ServeMux
}

func NewRouter() *Router {
    return &Router{mux: http.NewServeMux()}
}

func (r *Router) Group(prefix string) *RouteGroup {
    return &RouteGroup{router: r, prefix: prefix}
}

type RouteGroup struct {
    router *Router
    prefix string
}

func (g *RouteGroup) GET(path string, handler http.HandlerFunc) {
    g.router.mux.HandleFunc(g.prefix+path, func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            http.Error(w, "Method not allowed", 405)
            return
        }
        handler(w, r)
    })
}

func (g *RouteGroup) POST(path string, handler http.HandlerFunc) {
    g.router.mux.HandleFunc(g.prefix+path, func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "Method not allowed", 405)
            return
        }
        handler(w, r)
    })
}

// 使用
func main() {
    router := NewRouter()
    
    api := router.Group("/api")
    api.GET("/users", listUsers)
    api.POST("/users", createUser)
    
    http.ListenAndServe(":8080", router.mux)
}
```

---

## 3. 常见问题处理

### 3.1 跨域 CORS

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
```

### 3.2 请求体大小限制

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 限制请求体最大 1MB
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
    
    var data map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
        return
    }
}
```

### 3.3 优雅关闭

```go
func main() {
    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
    
    // 启动服务器
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // 优雅关闭，等待 30 秒
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited")
}
```

---

## 4. 性能优化

### 4.1 连接池配置

```go
transport := &http.Transport{
    MaxIdleConns:        100,              // 最大空闲连接数
    MaxIdleConnsPerHost: 10,               // 每个主机最大空闲连接
    MaxConnsPerHost:     100,              // 每个主机最大连接数
    IdleConnTimeout:     90 * time.Second, // 空闲连接超时
}

client := &http.Client{
    Transport: transport,
}
```

### 4.2 响应压缩

```go
import "compress/gzip"

func GzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检查客户端是否支持 gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        // 设置响应头
        w.Header().Set("Content-Encoding", "gzip")
        
        // 创建 gzip writer
        gz := gzip.NewWriter(w)
        defer gz.Close()
        
        // 包装 ResponseWriter
        gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
        next.ServeHTTP(gzw, r)
    })
}

type gzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}
```

### 4.3 JSON 性能优化

```go
import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 使用方式和标准库一样
json.Marshal(data)
json.Unmarshal(data, &result)
```

---

## 5. 安全相关

### 5.1 防止 SQL 注入

```go
// ❌ 危险：直接拼接
query := "SELECT * FROM users WHERE name = '" + name + "'"

// ✅ 安全：使用参数化查询
query := "SELECT * FROM users WHERE name = ?"
db.Query(query, name)
```

### 5.2 防止 XSS

```go
import "html"

// 转义 HTML
safe := html.EscapeString(userInput)
```

### 5.3 HTTPS

```go
// 使用 HTTPS
http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handler)
```

### 5.4 安全响应头

```go
func SecurityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 6. 常用 HTTP 框架

### 6.1 Gin（最流行）

```go
import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    
    r.GET("/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"users": users})
    })
    
    r.POST("/users", func(c *gin.Context) {
        var user User
        c.BindJSON(&user)
        c.JSON(201, user)
    })
    
    r.Run(":8080")
}
```

### 6.2 Echo

```go
import "github.com/labstack/echo/v4"

func main() {
    e := echo.New()
    
    e.GET("/users", func(c echo.Context) error {
        return c.JSON(200, users)
    })
    
    e.Start(":8080")
}
```

### 6.3 选择建议

- **标准库**：简单项目，学习用
- **Gin**：生产项目首选，性能好，生态丰富
- **Echo**：和 Gin 类似，API 更简洁

---

## 7. 练习题

### 练习 1：实现一个简单的 API

```go
// 实现一个待办事项 API
// GET    /api/todos      获取所有待办
// POST   /api/todos      创建待办
// PUT    /api/todos/:id  更新待办
// DELETE /api/todos/:id  删除待办

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

### 练习 2：添加中间件

```go
// 给上面的 API 添加：
// 1. 日志中间件：记录每个请求的方法、路径、耗时
// 2. 认证中间件：检查 Authorization 头
```

### 练习 3：HTTP 客户端

```go
// 写一个函数，并发请求多个 URL，返回所有结果
// 要求：有超时控制，任何一个失败不影响其他
func FetchAll(urls []string, timeout time.Duration) []Result {
    // 你来实现
}
```

---

## 总结

| 概念 | 说明 |
|------|------|
| HTTP | 浏览器和服务器通信的协议 |
| REST | 用 URL 定位资源，用 HTTP 方法操作 |
| 状态码 | 2xx 成功，4xx 客户端错误，5xx 服务器错误 |
| 中间件 | 在请求处理前后做一些事情 |
| CORS | 跨域资源共享 |
| HTTPS | HTTP + 加密 |

**核心记住**：
1. RESTful = 资源（名词）+ HTTP 方法（动词）
2. 状态码要用对
3. 响应格式要统一
4. 注意安全（HTTPS、参数校验、防注入）
