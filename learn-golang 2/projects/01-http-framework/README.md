# 项目一：实现一个简单的 HTTP 框架

## 目标
通过实现一个类似 Gin 的 HTTP 框架，深入理解：
- HTTP 协议
- 路由树（Trie）
- 中间件模式
- Context 设计

## 功能要求

### 阶段1：基础路由
- [ ] 支持 GET、POST、PUT、DELETE 方法
- [ ] 支持静态路由 `/users`
- [ ] 支持参数路由 `/users/:id`
- [ ] 支持通配符路由 `/static/*filepath`

### 阶段2：中间件
- [ ] 实现中间件链
- [ ] 实现 Logger 中间件
- [ ] 实现 Recovery 中间件
- [ ] 支持路由组中间件

### 阶段3：Context
- [ ] 封装 Request 和 Response
- [ ] 实现 JSON/HTML/String 响应
- [ ] 实现参数获取（Path/Query/Body）
- [ ] 实现 Context 传值

## API 设计

```go
func main() {
    r := mini.New()
    
    // 全局中间件
    r.Use(mini.Logger(), mini.Recovery())
    
    // 静态路由
    r.GET("/", func(c *mini.Context) {
        c.String(200, "Hello, World!")
    })
    
    // 参数路由
    r.GET("/users/:id", func(c *mini.Context) {
        id := c.Param("id")
        c.JSON(200, mini.H{"id": id})
    })
    
    // 路由组
    api := r.Group("/api")
    api.Use(authMiddleware())
    {
        api.GET("/users", listUsers)
        api.POST("/users", createUser)
    }
    
    r.Run(":8080")
}
```

## 实现提示

### 路由树结构
```go
type node struct {
    pattern  string  // 完整路由，如 /users/:id
    part     string  // 当前节点部分，如 :id
    children []*node // 子节点
    isWild   bool    // 是否是通配符节点
}
```

### 中间件链
```go
type HandlerFunc func(*Context)

type Context struct {
    handlers []HandlerFunc
    index    int
}

func (c *Context) Next() {
    c.index++
    for c.index < len(c.handlers) {
        c.handlers[c.index](c)
        c.index++
    }
}
```

## 参考资源
- [Gin 源码](https://github.com/gin-gonic/gin)
- [7天用Go从零实现Web框架](https://geektutu.com/post/gee.html)

## 学习收获
- 理解 HTTP 服务器的工作原理
- 掌握 Trie 树在路由匹配中的应用
- 理解中间件的洋葱模型
- 学会设计优雅的 API
