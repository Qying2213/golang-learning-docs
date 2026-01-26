# gRPC 教程 1：入门

> 边学边练，跟着做

---

## 🎯 本节目标

学完这节你能：
- 理解 gRPC 是什么
- 搭建环境
- 写出第一个 gRPC 服务

---

## 第 1 课：gRPC 是什么

### 1.1 一句话解释

**gRPC 让你像调用本地函数一样调用远程服务器的函数。**

```go
// 本地调用
result := Add(1, 2)

// gRPC 调用（看起来一样，但实际调用的是远程服务器）
result := client.Add(ctx, &pb.AddRequest{A: 1, B: 2})
```

### 1.2 gRPC vs HTTP

| | gRPC | HTTP/REST |
|---|---|---|
| 数据格式 | 二进制（小、快） | JSON（大、慢） |
| 适用场景 | 微服务内部通信 | 对外 API |

**什么时候用 gRPC**：微服务之间调用，追求性能。
**什么时候用 HTTP**：给前端/外部用的 API。

### ✅ 理解检查

回答这个问题：gRPC 和 HTTP 最大的区别是什么？

<details>
<summary>点击查看答案</summary>
gRPC 用二进制传输，更快更小；HTTP 用 JSON 文本，更通用。
</details>

---

## 第 2 课：搭建环境

### 2.1 安装 protoc

```bash
# macOS
brew install protobuf

# 验证
protoc --version
```

### 2.2 安装 Go 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2.3 配置 PATH

```bash
# 添加到 ~/.zshrc
export PATH="$PATH:$(go env GOPATH)/bin"

# 生效
source ~/.zshrc
```

### ✅ 动手检查

运行这两个命令，都能输出版本号就对了：
```bash
protoc --version
which protoc-gen-go
```

---

## 第 3 课：创建项目

### 3.1 创建文件夹

```bash
mkdir grpc-learn
cd grpc-learn
go mod init grpc-learn
```

### 3.2 安装依赖

```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

### 3.3 项目结构

```
grpc-learn/
├── hello.proto      # 定义服务
├── pb/              # 生成的代码
├── server/
│   └── main.go      # 服务端
├── client/
│   └── main.go      # 客户端
└── go.mod
```

### ✅ 动手检查

`go.mod` 文件存在，里面有 `grpc` 依赖。

---

## 第 4 课：写 Proto 文件

### 4.1 什么是 Proto

Proto 文件定义：
- **消息**（数据结构）
- **服务**（有哪些方法）

### 4.2 创建 hello.proto

```protobuf
syntax = "proto3";

package hello;

option go_package = "./pb";

// 请求消息
message HelloRequest {
    string name = 1;
}

// 响应消息
message HelloResponse {
    string message = 1;
}

// 服务定义
service Greeter {
    rpc SayHello(HelloRequest) returns (HelloResponse);
}
```

### 4.3 理解语法

```protobuf
message HelloRequest {
    string name = 1;  // 1 是字段编号，不是默认值
}
```

- `message` = Go 的 struct
- `string name = 1` = 字段名是 name，类型是 string，编号是 1
- `rpc SayHello(...)` = 定义一个远程方法

### ✅ 练习

在 `hello.proto` 里添加一个 `age` 字段（int32 类型）到 HelloRequest：

<details>
<summary>点击查看答案</summary>

```protobuf
message HelloRequest {
    string name = 1;
    int32 age = 2;
}
```
</details>

---

## 第 5 课：生成 Go 代码

### 5.1 运行命令

```bash
protoc --go_out=. --go-grpc_out=. hello.proto
```

### 5.2 检查生成的文件

```
pb/
├── hello.pb.go       # 消息定义
└── hello_grpc.pb.go  # 服务接口
```

### 5.3 看看生成了什么

打开 `pb/hello.pb.go`，你会看到：

```go
type HelloRequest struct {
    Name string
}

type HelloResponse struct {
    Message string
}
```

这就是 proto 里的 message 变成了 Go struct。

### ✅ 动手检查

`pb/` 文件夹里有两个 `.go` 文件。

---

## 第 6 课：写服务端

### 6.1 创建 server/main.go

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "grpc-learn/pb"
    "google.golang.org/grpc"
)

// 服务结构体
type server struct {
    pb.UnimplementedGreeterServer
}

// 实现 SayHello 方法
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    fmt.Println("收到请求:", req.Name)
    return &pb.HelloResponse{
        Message: "你好, " + req.Name + "!",
    }, nil
}

func main() {
    // 监听端口
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }

    // 创建 gRPC 服务器
    s := grpc.NewServer()

    // 注册服务
    pb.RegisterGreeterServer(s, &server{})

    fmt.Println("服务端启动，监听 :50051")
    s.Serve(lis)
}
```

### 6.2 理解代码

```go
// 这是关键：实现 proto 里定义的方法
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    // req.Name 是客户端传来的
    // 返回 HelloResponse 给客户端
    return &pb.HelloResponse{Message: "你好"}, nil
}
```

### ✅ 练习

修改 SayHello，让它返回 "Hello, xxx! 你今年 xx 岁"（用上 age 字段）：

<details>
<summary>点击查看答案</summary>

```go
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    msg := fmt.Sprintf("Hello, %s! 你今年 %d 岁", req.Name, req.Age)
    return &pb.HelloResponse{Message: msg}, nil
}
```
</details>

---

## 第 7 课：写客户端

### 7.1 创建 client/main.go

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "grpc-learn/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // 连接服务器
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // 创建客户端
    client := pb.NewGreeterClient(conn)

    // 设置超时
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // 调用方法
    resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("收到响应:", resp.Message)
}
```

### 7.2 理解代码

```go
// 这行就是调用远程方法
resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "秦阳"})
// 看起来像本地调用，实际上是网络请求
```

---

## 第 8 课：运行测试

### 8.1 启动服务端

```bash
# 终端 1
go run server/main.go
```

输出：`服务端启动，监听 :50051`

### 8.2 运行客户端

```bash
# 终端 2
go run client/main.go
```

输出：`收到响应: 你好, 秦阳!`

### ✅ 成功检查

客户端能收到响应就成功了！

---

## 📝 本节练习

### 练习 1：添加 SayGoodbye 方法

1. 在 `hello.proto` 添加：
```protobuf
rpc SayGoodbye(HelloRequest) returns (HelloResponse);
```

2. 重新生成代码：
```bash
protoc --go_out=. --go-grpc_out=. hello.proto
```

3. 在服务端实现 SayGoodbye

4. 在客户端调用它

### 练习 2：计算器服务

创建一个新的 `calculator.proto`，实现加法：

```protobuf
message AddRequest {
    int32 a = 1;
    int32 b = 2;
}

message AddResponse {
    int32 result = 1;
}

service Calculator {
    rpc Add(AddRequest) returns (AddResponse);
}
```

然后实现服务端和客户端，测试 `1 + 2 = 3`。

---

## 🐛 常见错误

**错误 1**：`protoc: command not found`
```
解决：brew install protobuf
```

**错误 2**：`pb/hello.pb.go: no such file`
```
解决：先运行 protoc 生成代码
```

**错误 3**：`connection refused`
```
解决：确保服务端先启动
```

**错误 4**：`undefined: pb.xxx`
```
解决：go mod tidy
```

---

## 下一节

完成练习后，继续 [gRPC教程2-进阶](./gRPC教程2-进阶.md)
