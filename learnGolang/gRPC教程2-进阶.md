# gRPC 教程 2：进阶

> 边学边练，跟着做

---

## 🎯 本节目标

学完这节你能：
- 处理错误
- 使用拦截器（中间件）
- 传递元数据（类似 HTTP Header）
- 设置超时

---

## 第 1 课：错误处理

### 1.1 gRPC 状态码

gRPC 有自己的错误码，常用的：

```go
codes.OK              // 成功
codes.NotFound        // 未找到
codes.InvalidArgument // 参数错误
codes.Unauthenticated // 未认证
codes.Internal        // 内部错误
```

### 1.2 服务端返回错误

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    if req.Id <= 0 {
        // 返回参数错误
        return nil, status.Error(codes.InvalidArgument, "id 必须大于 0")
    }
    
    user := findUser(req.Id)
    if user == nil {
        // 返回未找到
        return nil, status.Errorf(codes.NotFound, "用户 %d 不存在", req.Id)
    }
    
    return user, nil
}
```

### 1.3 客户端处理错误

```go
resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 999})
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.NotFound:
            fmt.Println("用户不存在")
        case codes.InvalidArgument:
            fmt.Println("参数错误:", st.Message())
        default:
            fmt.Println("其他错误:", st.Message())
        }
    }
    return
}
```

### ✅ 练习

修改你的 SayHello 方法，当 name 为空时返回 `InvalidArgument` 错误：

<details>
<summary>点击查看答案</summary>

```go
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    if req.Name == "" {
        return nil, status.Error(codes.InvalidArgument, "name 不能为空")
    }
    return &pb.HelloResponse{Message: "你好, " + req.Name}, nil
}
```
</details>

---

## 第 2 课：拦截器（中间件）

### 2.1 什么是拦截器

拦截器 = gRPC 的中间件，在请求前后做一些事情。

```
请求 → 拦截器1 → 拦截器2 → 实际方法 → 拦截器2 → 拦截器1 → 响应
```

### 2.2 日志拦截器

```go
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    
    // 调用实际方法
    resp, err := handler(ctx, req)
    
    // 打印日志
    fmt.Printf("[%s] %s %v\n", info.FullMethod, time.Since(start), err)
    
    return resp, err
}
```

### 2.3 注册拦截器

```go
func main() {
    s := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
    )
    // ...
}
```

### 2.4 多个拦截器

```go
s := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        loggingInterceptor,
        authInterceptor,
    ),
)
```

### ✅ 练习

给你的服务端添加一个日志拦截器，打印每个请求的方法名和耗时：

<details>
<summary>点击查看答案</summary>

```go
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    fmt.Printf("方法: %s, 耗时: %v\n", info.FullMethod, time.Since(start))
    return resp, err
}

func main() {
    lis, _ := net.Listen("tcp", ":50051")
    
    s := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
    )
    
    pb.RegisterGreeterServer(s, &server{})
    s.Serve(lis)
}
```
</details>

---

## 第 3 课：Metadata（元数据）

### 3.1 什么是 Metadata

Metadata 类似 HTTP Header，用来传递额外信息（比如 Token）。

### 3.2 客户端发送 Metadata

```go
import "google.golang.org/grpc/metadata"

func main() {
    // 创建 metadata
    md := metadata.Pairs(
        "authorization", "Bearer my-token",
        "request-id", "123456",
    )
    
    // 添加到 context
    ctx := metadata.NewOutgoingContext(context.Background(), md)
    
    // 发送请求
    resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
}
```

### 3.3 服务端接收 Metadata

```go
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    // 获取 metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if ok {
        // 获取 token
        if tokens := md.Get("authorization"); len(tokens) > 0 {
            fmt.Println("Token:", tokens[0])
        }
    }
    
    return &pb.HelloResponse{Message: "你好"}, nil
}
```

### 3.4 认证拦截器

```go
func authInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // 获取 metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "没有 metadata")
    }
    
    // 获取 token
    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return nil, status.Error(codes.Unauthenticated, "没有 token")
    }
    
    // 验证 token（这里简化了）
    if tokens[0] != "Bearer valid-token" {
        return nil, status.Error(codes.Unauthenticated, "token 无效")
    }
    
    // 验证通过，继续处理
    return handler(ctx, req)
}
```

### ✅ 练习

1. 客户端发送一个 `user-id` 的 metadata
2. 服务端接收并打印出来

<details>
<summary>点击查看答案</summary>

客户端：
```go
md := metadata.Pairs("user-id", "123")
ctx := metadata.NewOutgoingContext(context.Background(), md)
resp, _ := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
```

服务端：
```go
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    md, _ := metadata.FromIncomingContext(ctx)
    if ids := md.Get("user-id"); len(ids) > 0 {
        fmt.Println("User ID:", ids[0])
    }
    return &pb.HelloResponse{Message: "你好"}, nil
}
```
</details>

---

## 第 4 课：超时控制

### 4.1 客户端设置超时

```go
// 5 秒超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
if err != nil {
    st, _ := status.FromError(err)
    if st.Code() == codes.DeadlineExceeded {
        fmt.Println("请求超时了")
    }
}
```

### 4.2 服务端检查超时

```go
func (s *server) SlowMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    for i := 0; i < 10; i++ {
        // 检查是否超时
        select {
        case <-ctx.Done():
            return nil, status.Error(codes.Canceled, "客户端取消了")
        default:
        }
        
        time.Sleep(time.Second)
    }
    return &pb.Response{}, nil
}
```

### ✅ 练习

1. 让 SayHello 方法 sleep 3 秒
2. 客户端设置 1 秒超时
3. 观察超时错误

<details>
<summary>点击查看答案</summary>

服务端：
```go
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    time.Sleep(3 * time.Second)  // 模拟慢操作
    return &pb.HelloResponse{Message: "你好"}, nil
}
```

客户端：
```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
if err != nil {
    st, _ := status.FromError(err)
    if st.Code() == codes.DeadlineExceeded {
        fmt.Println("超时了！")
    }
}
```
</details>

---

## 📝 本节练习

### 综合练习：带认证的用户服务

1. 创建 `user.proto`：
```protobuf
message User {
    int64 id = 1;
    string name = 2;
}

message GetUserRequest {
    int64 id = 1;
}

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
}
```

2. 实现服务端：
   - 添加日志拦截器
   - 添加认证拦截器（检查 authorization header）
   - GetUser 方法：id 不存在返回 NotFound 错误

3. 实现客户端：
   - 发送 authorization metadata
   - 设置 3 秒超时
   - 处理各种错误

---

## 下一节

完成练习后，继续 [gRPC教程3-流式与实战](./gRPC教程3-流式与实战.md)
