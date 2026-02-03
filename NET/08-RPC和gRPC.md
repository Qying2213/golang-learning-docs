# RPC 和 gRPC

> **重要程度：⭐⭐⭐ 必须掌握**
> 微服务架构下，gRPC 是服务间通信的首选！

## 📚 本章学习目标

学完本章，你将能够：

- 理解 RPC 的原理
- 掌握 gRPC 的使用
- 学会 Protobuf 定义接口
- 实现 gRPC 服务端和客户端

---

## 1. RPC 是什么？⭐⭐⭐

### 1.1 RPC 的概念

```
RPC = Remote Procedure Call
      远程过程调用

简单说：像调用本地函数一样调用远程服务

本地调用：
result := add(1, 2)

远程调用（看起来一样）：
result := remoteService.Add(1, 2)

但实际上：
1. 参数被序列化
2. 通过网络发送到远程服务器
3. 远程服务器执行函数
4. 结果序列化返回
5. 客户端反序列化得到结果
```

### 1.2 RPC vs HTTP

```
HTTP RESTful：
- 基于 HTTP 协议
- 文本格式（JSON）
- 可读性好
- 性能一般
- 适合对外 API

RPC：
- 可以基于 TCP/HTTP2
- 二进制格式（Protobuf）
- 性能高
- 强类型
- 适合内部服务通信

选择建议：
- 对外 API → RESTful
- 内部微服务 → gRPC
```

### 1.3 RPC 调用流程

```
┌─────────────┐                           ┌─────────────┐
│   客户端     │                           │   服务端     │
├─────────────┤                           ├─────────────┤
│  业务代码    │                           │  业务代码    │
│     ↓       │                           │     ↑       │
│  客户端存根  │ ──── 网络传输 ────────→   │  服务端存根  │
│  (Stub)     │ ←─── 网络传输 ────────    │  (Skeleton) │
└─────────────┘                           └─────────────┘

流程：
1. 客户端调用本地 Stub
2. Stub 将参数序列化
3. 通过网络发送请求
4. 服务端 Skeleton 接收请求
5. 反序列化参数，调用真实函数
6. 将结果序列化返回
7. 客户端 Stub 反序列化结果
8. 返回给调用者
```

---

## 2. gRPC 简介 ⭐⭐⭐

### 2.1 什么是 gRPC

```
gRPC = Google Remote Procedure Call

特点：
1. Google 开源
2. 基于 HTTP/2
3. 使用 Protobuf 序列化
4. 支持多种语言
5. 支持流式传输
6. 高性能

HTTP/2 优势：
- 多路复用（一个连接多个请求）
- 头部压缩
- 服务端推送
- 二进制传输
```

### 2.2 gRPC 四种通信模式

```
1. 一元 RPC（Unary）
   客户端发一个请求，服务端返回一个响应
   最常用的模式

   rpc GetUser(GetUserRequest) returns (User);

2. 服务端流（Server Streaming）
   客户端发一个请求，服务端返回多个响应
   适合：下载大文件、实时数据推送

   rpc ListUsers(ListRequest) returns (stream User);

3. 客户端流（Client Streaming）
   客户端发多个请求，服务端返回一个响应
   适合：上传大文件

   rpc UploadFile(stream Chunk) returns (UploadResponse);

4. 双向流（Bidirectional Streaming）
   客户端和服务端都可以发送多个消息
   适合：聊天、实时游戏

   rpc Chat(stream Message) returns (stream Message);
```

---

## 3. Protobuf 入门 ⭐⭐⭐

### 3.1 什么是 Protobuf

```
Protobuf = Protocol Buffers

Google 开发的序列化格式：
- 二进制格式，体积小
- 序列化/反序列化快
- 强类型，有 Schema
- 跨语言支持

对比 JSON：
JSON: {"name":"张三","age":20}  → 约 30 字节
Protobuf: 二进制编码           → 约 10 字节
```

### 3.2 Protobuf 语法

```protobuf
// user.proto
syntax = "proto3";  // 使用 proto3 语法

package user;  // 包名

option go_package = "pb/user";  // Go 包路径

// 定义消息（类似 struct）
message User {
    int64 id = 1;           // 字段编号，不是默认值
    string name = 2;
    string email = 3;
    int32 age = 4;
    repeated string tags = 5;  // 数组
    Address address = 6;       // 嵌套消息
}

message Address {
    string city = 1;
    string street = 2;
}

// 定义服务
service UserService {
    // 一元 RPC
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (Empty);

    // 服务端流
    rpc ListUsers(ListUsersRequest) returns (stream User);
}

// 请求消息
message GetUserRequest {
    int64 id = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
    int32 age = 3;
}

message UpdateUserRequest {
    int64 id = 1;
    string name = 2;
    string email = 3;
}

message DeleteUserRequest {
    int64 id = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 size = 2;
}

message Empty {}
```

### 3.3 字段类型

```protobuf
// 标量类型
int32, int64      // 整数
uint32, uint64    // 无符号整数
float, double     // 浮点数
bool              // 布尔
string            // 字符串
bytes             // 字节数组

// 复合类型
message           // 消息（类似 struct）
enum              // 枚举
repeated          // 数组
map               // 映射

// 示例
message Example {
    int64 id = 1;
    string name = 2;
    bool active = 3;
    repeated string tags = 4;           // []string
    map<string, int32> scores = 5;      // map[string]int32
    Status status = 6;                  // 枚举
}

enum Status {
    UNKNOWN = 0;
    ACTIVE = 1;
    INACTIVE = 2;
}
```

### 3.4 安装 Protobuf 工具

```bash
# 安装 protoc 编译器
# Mac
brew install protobuf

# Ubuntu
sudo apt install protobuf-compiler

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 验证安装
protoc --version
```

### 3.5 生成 Go 代码

```bash
# 目录结构
project/
├── proto/
│   └── user.proto
├── pb/
│   └── user/
│       ├── user.pb.go        # 消息定义
│       └── user_grpc.pb.go   # gRPC 服务定义
└── main.go

# 生成命令
protoc --go_out=. --go-grpc_out=. proto/user.proto

# 或者指定路径
protoc \
    --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    proto/user.proto
```

---

## 4. gRPC 实战 ⭐⭐⭐

### 4.1 完整项目结构

```
grpc-demo/
├── proto/
│   └── user.proto
├── pb/
│   └── user/
│       ├── user.pb.go
│       └── user_grpc.pb.go
├── server/
│   └── main.go
├── client/
│   └── main.go
├── go.mod
└── go.sum
```

### 4.2 Proto 文件

```protobuf
// proto/user.proto
syntax = "proto3";

package user;

option go_package = "grpc-demo/pb/user";

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
    int32 age = 4;
}

message GetUserRequest {
    int64 id = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
    int32 age = 3;
}

message ListUsersRequest {
    int32 page = 1;
    int32 size = 2;
}

message DeleteUserRequest {
    int64 id = 1;
}

message Empty {}

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
    rpc DeleteUser(DeleteUserRequest) returns (Empty);
}
```

### 4.3 服务端实现

```go
// server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"

    "google.golang.org/grpc"
    pb "grpc-demo/pb/user"
)

// 用户服务实现
type userServer struct {
    pb.UnimplementedUserServiceServer
    mu    sync.RWMutex
    users map[int64]*pb.User
    nextID int64
}

func newUserServer() *userServer {
    return &userServer{
        users:  make(map[int64]*pb.User),
        nextID: 1,
    }
}

// 获取用户
func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, ok := s.users[req.Id]
    if !ok {
        return nil, fmt.Errorf("user not found: %d", req.Id)
    }
    return user, nil
}

// 创建用户
func (s *userServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    user := &pb.User{
        Id:    s.nextID,
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }
    s.users[user.Id] = user
    s.nextID++

    log.Printf("创建用户: %+v", user)
    return user, nil
}

// 列出用户（服务端流）
func (s *userServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, user := range s.users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    return nil
}

// 删除用户
func (s *userServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, ok := s.users[req.Id]; !ok {
        return nil, fmt.Errorf("user not found: %d", req.Id)
    }
    delete(s.users, req.Id)
    log.Printf("删除用户: %d", req.Id)
    return &pb.Empty{}, nil
}

func main() {
    // 监听端口
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("监听失败: %v", err)
    }

    // 创建 gRPC 服务器
    grpcServer := grpc.NewServer()

    // 注册服务
    pb.RegisterUserServiceServer(grpcServer, newUserServer())

    log.Println("gRPC 服务器启动: :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("服务失败: %v", err)
    }
}
```

### 4.4 客户端实现

```go
// client/main.go
package main

import (
    "context"
    "io"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "grpc-demo/pb/user"
)

func main() {
    // 连接服务器
    // 注意：grpc.Dial 在新版本中已 deprecated，推荐使用 grpc.NewClient
    // Go 1.21+ 可使用：conn, err := grpc.NewClient("localhost:50051", ...)
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer conn.Close()

    // 创建客户端
    client := pb.NewUserServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 创建用户
    user1, err := client.CreateUser(ctx, &pb.CreateUserRequest{
        Name:  "张三",
        Email: "zhangsan@example.com",
        Age:   25,
    })
    if err != nil {
        log.Fatalf("创建用户失败: %v", err)
    }
    log.Printf("创建用户成功: %+v", user1)

    // 创建另一个用户
    user2, _ := client.CreateUser(ctx, &pb.CreateUserRequest{
        Name:  "李四",
        Email: "lisi@example.com",
        Age:   30,
    })
    log.Printf("创建用户成功: %+v", user2)

    // 获取用户
    user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
    if err != nil {
        log.Fatalf("获取用户失败: %v", err)
    }
    log.Printf("获取用户: %+v", user)

    // 列出所有用户（服务端流）
    stream, err := client.ListUsers(ctx, &pb.ListUsersRequest{})
    if err != nil {
        log.Fatalf("列出用户失败: %v", err)
    }

    log.Println("所有用户:")
    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("接收失败: %v", err)
        }
        log.Printf("  - %+v", user)
    }

    // 删除用户
    _, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: 1})
    if err != nil {
        log.Fatalf("删除用户失败: %v", err)
    }
    log.Println("删除用户成功")
}
```

### 4.5 运行测试

```bash
# 安装依赖
go mod init grpc-demo
go get google.golang.org/grpc
go get google.golang.org/protobuf

# 生成代码
protoc --go_out=. --go-grpc_out=. proto/user.proto

# 启动服务端
go run server/main.go

# 另一个终端，运行客户端
go run client/main.go
```

---

## 5. gRPC 流式传输 ⭐⭐

### 5.1 服务端流示例

```protobuf
// 下载文件
service FileService {
    rpc Download(DownloadRequest) returns (stream Chunk);
}

message DownloadRequest {
    string filename = 1;
}

message Chunk {
    bytes data = 1;
}
```

```go
// 服务端
func (s *fileServer) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
    file, _ := os.Open(req.Filename)
    defer file.Close()

    buf := make([]byte, 1024)
    for {
        n, err := file.Read(buf)
        if err == io.EOF {
            break
        }
        stream.Send(&pb.Chunk{Data: buf[:n]})
    }
    return nil
}

// 客户端
stream, _ := client.Download(ctx, &pb.DownloadRequest{Filename: "test.txt"})
for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // 处理 chunk.Data
}
```

### 5.2 客户端流示例

```protobuf
// 上传文件
service FileService {
    rpc Upload(stream Chunk) returns (UploadResponse);
}

message UploadResponse {
    int64 size = 1;
}
```

```go
// 客户端
stream, _ := client.Upload(ctx)
file, _ := os.Open("test.txt")
buf := make([]byte, 1024)

for {
    n, err := file.Read(buf)
    if err == io.EOF {
        break
    }
    stream.Send(&pb.Chunk{Data: buf[:n]})
}

response, _ := stream.CloseAndRecv()
log.Printf("上传完成，大小: %d", response.Size)
```

### 5.3 双向流示例

```protobuf
// 聊天
service ChatService {
    rpc Chat(stream Message) returns (stream Message);
}

message Message {
    string user = 1;
    string content = 2;
}
```

```go
// 客户端
stream, _ := client.Chat(ctx)

// 发送消息
go func() {
    for {
        var input string
        fmt.Scanln(&input)
        stream.Send(&pb.Message{User: "张三", Content: input})
    }
}()

// 接收消息
for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    fmt.Printf("%s: %s\n", msg.User, msg.Content)
}
```

---

## 6. gRPC 拦截器 ⭐⭐

### 6.1 什么是拦截器

```
拦截器 = 中间件

可以在 RPC 调用前后执行逻辑：
- 日志记录
- 认证授权
- 错误处理
- 性能监控
```

### 6.2 一元拦截器

```go
// 服务端拦截器
func unaryServerInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // 前置处理
    start := time.Now()
    log.Printf("开始调用: %s", info.FullMethod)

    // 调用实际方法
    resp, err := handler(ctx, req)

    // 后置处理
    log.Printf("调用完成: %s, 耗时: %v, 错误: %v",
        info.FullMethod, time.Since(start), err)

    return resp, err
}

// 注册拦截器
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(unaryServerInterceptor),
)
```

### 6.3 认证拦截器

```go
// 服务端认证拦截器
func authInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // 从 metadata 获取 token
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }

    // 验证 token
    if !validateToken(tokens[0]) {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }

    return handler(ctx, req)
}

// 客户端添加 token
ctx := metadata.AppendToOutgoingContext(ctx, "authorization", "your-token")
client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
```

---

## 7. gRPC vs REST 对比 ⭐⭐⭐

| 特性       | gRPC               | REST         |
| ---------- | ------------------ | ------------ |
| 协议       | HTTP/2             | HTTP/1.1     |
| 数据格式   | Protobuf（二进制） | JSON（文本） |
| 性能       | 高                 | 一般         |
| 类型安全   | 强类型             | 弱类型       |
| 流式传输   | 支持               | 不支持       |
| 浏览器支持 | 需要 gRPC-Web      | 原生支持     |
| 可读性     | 差（二进制）       | 好（JSON）   |
| 适用场景   | 内部微服务         | 对外 API     |

---

## 8. 本章总结

### 必须记住的内容

```
RPC 概念：
- 远程过程调用
- 像调用本地函数一样调用远程服务

gRPC 特点：
- Google 开源
- 基于 HTTP/2
- 使用 Protobuf
- 高性能

四种模式：
- 一元 RPC
- 服务端流
- 客户端流
- 双向流

Protobuf：
- 二进制序列化
- 强类型
- 需要定义 .proto 文件
```

### 面试常见问题

```
Q: gRPC 和 REST 的区别？
A: gRPC 基于 HTTP/2 和 Protobuf，性能高，适合内部服务；
   REST 基于 HTTP/1.1 和 JSON，可读性好，适合对外 API。

Q: Protobuf 的优点？
A: 二进制格式体积小，序列化快，强类型，跨语言支持。

Q: gRPC 的四种通信模式？
A: 一元 RPC、服务端流、客户端流、双向流。

Q: 什么场景用 gRPC？
A: 微服务内部通信、需要高性能的场景、需要流式传输的场景。
```

---

下一章：[09-IP协议和路由](./09-IP协议和路由.md)
