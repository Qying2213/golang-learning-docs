# gRPC 拦截器教程（边学边练）

## 什么是拦截器？

**拦截器 = gRPC 的中间件**

就像你在 Gin 里用的中间件一样：

```
Gin 中间件：
请求 → CORS中间件 → JWT中间件 → Handler → 响应

gRPC 拦截器：
请求 → 日志拦截器 → 认证拦截器 → Handler → 响应
```

**拦截器能做什么？**
- 记录日志（每个请求记录耗时、参数）
- 认证鉴权（检查 Token）
- 错误处理（统一处理 panic）
- 监控统计（请求次数、响应时间）

---

## 拦截器的两种类型

| 类型 | 说明 | 使用场景 |
|------|------|---------|
| **Unary 拦截器** | 一元调用（一问一答） | 普通 RPC 调用 |
| **Stream 拦截器** | 流式调用 | 流式 RPC |

我们先学 **Unary 拦截器**，这是最常用的。

---

## 第一步：理解拦截器函数签名

### 服务端拦截器

```go
// 服务端 Unary 拦截器的函数签名
type UnaryServerInterceptor func(
    ctx context.Context,                    // 上下文
    req interface{},                        // 请求参数
    info *grpc.UnaryServerInfo,            // 方法信息（方法名等）
    handler grpc.UnaryHandler,             // 真正的处理函数
) (resp interface{}, err error)
```

**参数解释：**
- `ctx`：上下文，可以传递数据
- `req`：客户端发来的请求
- `info`：包含方法名等信息
- `handler`：真正的业务处理函数，你需要调用它

### 客户端拦截器

```go
// 客户端 Unary 拦截器的函数签名
type UnaryClientInterceptor func(
    ctx context.Context,                    // 上下文
    method string,                          // 调用的方法名
    req, reply interface{},                 // 请求和响应
    cc *grpc.ClientConn,                   // 连接
    invoker grpc.UnaryInvoker,             // 真正的调用函数
    opts ...grpc.CallOption,               // 调用选项
) error
```

---

## 第二步：写一个日志拦截器（服务端）

### 2.1 最简单的拦截器

```go
package main

import (
    "context"
    "fmt"
    "time"

    "google.golang.org/grpc"
)

// 日志拦截器
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // ========== 请求前 ==========
    start := time.Now()
    fmt.Printf("[请求开始] 方法: %s\n", info.FullMethod)

    // ========== 调用真正的处理函数 ==========
    resp, err := handler(ctx, req)

    // ========== 请求后 ==========
    duration := time.Since(start)
    fmt.Printf("[请求结束] 方法: %s, 耗时: %v, 错误: %v\n", 
        info.FullMethod, duration, err)

    return resp, err
}
```

**关键点：**
1. `handler(ctx, req)` 是调用真正的业务逻辑
2. 在 `handler` 前后可以做任何事情
3. 必须返回 `handler` 的结果

### 2.2 注册拦截器

```go
func main() {
    // 创建 gRPC 服务器时注册拦截器
    s := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),  // 注册拦截器
    )
    
    // ... 注册服务、启动等
}
```

---

## 第三步：完整示例（边学边练）

### 3.1 项目结构

```
grpc-interceptor/
├── go.mod
├── hello.proto
├── pb/
│   ├── hello.pb.go
│   └── hello_grpc.pb.go
├── server/
│   └── main.go
└── client/
    └── main.go
```

### 3.2 proto 文件（hello.proto）

```protobuf
syntax = "proto3";

package hello;

option go_package = "./pb";

message HelloRequest {
    string name = 1;
}

message HelloResponse {
    string message = 1;
}

service Greeter {
    rpc SayHello(HelloRequest) returns (HelloResponse);
}
```

生成代码：
```bash
protoc --go_out=. --go-grpc_out=. hello.proto
```

### 3.3 服务端代码（server/main.go）

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"

    "grpc-interceptor/pb"

    "google.golang.org/grpc"
)

// ========== 服务实现 ==========
type server struct {
    pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    // 模拟业务处理
    time.Sleep(100 * time.Millisecond)
    return &pb.HelloResponse{
        Message: "你好 " + req.Name,
    }, nil
}

// ========== 拦截器1：日志 ==========
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    
    // 打印请求信息
    fmt.Printf("📥 [请求] 方法: %s, 参数: %+v\n", info.FullMethod, req)
    
    // 调用真正的处理函数
    resp, err := handler(ctx, req)
    
    // 打印响应信息
    duration := time.Since(start)
    if err != nil {
        fmt.Printf("📤 [响应] 方法: %s, 耗时: %v, 错误: %v\n", 
            info.FullMethod, duration, err)
    } else {
        fmt.Printf("📤 [响应] 方法: %s, 耗时: %v, 结果: %+v\n", 
            info.FullMethod, duration, resp)
    }
    
    return resp, err
}

// ========== 拦截器2：耗时统计 ==========
func timingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    
    resp, err := handler(ctx, req)
    
    duration := time.Since(start)
    if duration > 200*time.Millisecond {
        fmt.Printf("⚠️ [慢请求] 方法: %s, 耗时: %v\n", info.FullMethod, duration)
    }
    
    return resp, err
}

// ========== 拦截器3：panic 恢复 ==========
func recoveryInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    // 捕获 panic
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("🔥 [Panic] 方法: %s, 错误: %v\n", info.FullMethod, r)
            err = fmt.Errorf("内部错误")
        }
    }()
    
    return handler(ctx, req)
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }

    // ========== 注册多个拦截器 ==========
    s := grpc.NewServer(
        grpc.ChainUnaryInterceptor(  // 链式注册多个拦截器
            recoveryInterceptor,     // 第1个执行（最外层）
            loggingInterceptor,      // 第2个执行
            timingInterceptor,       // 第3个执行（最内层）
        ),
    )

    pb.RegisterGreeterServer(s, &server{})

    fmt.Println("🚀 服务启动在 :50051")
    s.Serve(lis)
}
```

### 3.4 客户端代码（client/main.go）

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "grpc-interceptor/pb"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// ========== 客户端拦截器：日志 ==========
func clientLoggingInterceptor(
    ctx context.Context,
    method string,
    req, reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption,
) error {
    start := time.Now()
    
    fmt.Printf("📤 [发送] 方法: %s, 参数: %+v\n", method, req)
    
    // 调用真正的 RPC
    err := invoker(ctx, method, req, reply, cc, opts...)
    
    duration := time.Since(start)
    if err != nil {
        fmt.Printf("📥 [接收] 方法: %s, 耗时: %v, 错误: %v\n", method, duration, err)
    } else {
        fmt.Printf("📥 [接收] 方法: %s, 耗时: %v, 结果: %+v\n", method, duration, reply)
    }
    
    return err
}

func main() {
    // 连接时注册客户端拦截器
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(clientLoggingInterceptor),  // 客户端拦截器
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewGreeterClient(conn)

    // 调用
    resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "秦阳"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("最终结果:", resp.Message)
}
```

---

## 第四步：运行测试

**1. 启动服务端：**
```bash
go run server/main.go
```

**2. 运行客户端：**
```bash
go run client/main.go
```

**服务端输出：**
```
🚀 服务启动在 :50051
📥 [请求] 方法: /hello.Greeter/SayHello, 参数: name:"秦阳"
📤 [响应] 方法: /hello.Greeter/SayHello, 耗时: 100.123ms, 结果: message:"你好 秦阳"
```

**客户端输出：**
```
📤 [发送] 方法: /hello.Greeter/SayHello, 参数: name:"秦阳"
📥 [接收] 方法: /hello.Greeter/SayHello, 耗时: 102.456ms, 结果: message:"你好 秦阳"
最终结果: 你好 秦阳
```

---

## 第五步：认证拦截器（实战常用）

### 5.1 服务端认证拦截器

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
)

// 认证拦截器
func authInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // 跳过不需要认证的方法
    if info.FullMethod == "/hello.Greeter/Login" {
        return handler(ctx, req)
    }

    // 从 metadata 获取 token
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "缺少认证信息")
    }

    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return nil, status.Error(codes.Unauthenticated, "缺少 Token")
    }

    token := tokens[0]
    
    // 验证 token（这里简化处理）
    if token != "Bearer valid-token" {
        return nil, status.Error(codes.Unauthenticated, "Token 无效")
    }

    // 可以把用户信息存到 context
    // newCtx := context.WithValue(ctx, "userID", 123)
    // return handler(newCtx, req)

    return handler(ctx, req)
}
```

### 5.2 客户端传递 Token

```go
import "google.golang.org/grpc/metadata"

func main() {
    // ... 建立连接

    // 创建带 Token 的 context
    md := metadata.Pairs("authorization", "Bearer valid-token")
    ctx := metadata.NewOutgoingContext(context.Background(), md)

    // 调用时传入带 Token 的 ctx
    resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
}
```

---

## 拦截器执行顺序

```go
grpc.ChainUnaryInterceptor(
    interceptor1,  // 第1个
    interceptor2,  // 第2个
    interceptor3,  // 第3个
)
```

**执行顺序（洋葱模型）：**

```
请求进来
    ↓
interceptor1 前半部分
    ↓
interceptor2 前半部分
    ↓
interceptor3 前半部分
    ↓
handler（真正的业务逻辑）
    ↓
interceptor3 后半部分
    ↓
interceptor2 后半部分
    ↓
interceptor1 后半部分
    ↓
响应返回
```

**和 Gin 中间件一样的洋葱模型！**

---

## 练习

### 练习1：写一个请求计数拦截器

统计服务收到了多少个请求：

```go
var requestCount int64

func countInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // TODO: 实现
    // 提示：用 atomic.AddInt64 原子操作
}
```

<details>
<summary>点击查看答案</summary>

```go
import "sync/atomic"

var requestCount int64

func countInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // 原子加1
    count := atomic.AddInt64(&requestCount, 1)
    fmt.Printf("📊 第 %d 个请求\n", count)
    
    return handler(ctx, req)
}
```

</details>

### 练习2：写一个限流拦截器

每秒最多处理 10 个请求：

```go
func rateLimitInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // TODO: 实现
    // 提示：可以用 channel 或 time.Ticker
}
```

<details>
<summary>点击查看答案</summary>

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// 简单的令牌桶
var limiter = make(chan struct{}, 10)

func init() {
    // 每100ms放一个令牌
    go func() {
        for {
            time.Sleep(100 * time.Millisecond)
            select {
            case limiter <- struct{}{}:
            default:
            }
        }
    }()
}

func rateLimitInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    select {
    case <-limiter:
        // 获取到令牌，继续处理
        return handler(ctx, req)
    default:
        // 没有令牌，拒绝请求
        return nil, status.Error(codes.ResourceExhausted, "请求太频繁")
    }
}
```

</details>

---

## 总结

| 概念 | 说明 |
|------|------|
| 拦截器 | gRPC 的中间件 |
| `grpc.UnaryInterceptor` | 注册单个拦截器 |
| `grpc.ChainUnaryInterceptor` | 注册多个拦截器（链式） |
| `handler(ctx, req)` | 调用真正的业务逻辑 |
| 执行顺序 | 洋葱模型（和 Gin 一样） |

**常用拦截器：**
1. 日志拦截器 - 记录请求响应
2. 认证拦截器 - 验证 Token
3. 恢复拦截器 - 捕获 panic
4. 限流拦截器 - 防止过载
5. 监控拦截器 - 统计指标
