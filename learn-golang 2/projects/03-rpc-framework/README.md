# 项目三：实现一个简单的 RPC 框架

## 目标
通过实现 RPC 框架，深入理解：
- 网络编程
- 序列化/反序列化
- 反射
- 服务注册与发现

## 功能要求

### 阶段1：基础 RPC
- [ ] 服务端注册服务
- [ ] 客户端调用服务
- [ ] 支持多种序列化（JSON/Gob）
- [ ] 错误处理

### 阶段2：高级功能
- [ ] 连接复用
- [ ] 超时控制
- [ ] 异步调用
- [ ] 负载均衡

### 阶段3：服务治理
- [ ] 服务注册
- [ ] 服务发现
- [ ] 健康检查
- [ ] 熔断器

## API 设计

```go
// 服务端
type Arith struct{}

func (a *Arith) Add(args *Args, reply *int) error {
    *reply = args.A + args.B
    return nil
}

func main() {
    server := rpc.NewServer()
    server.Register(&Arith{})
    server.Serve(":9000")
}

// 客户端
func main() {
    client, _ := rpc.Dial("tcp", "localhost:9000")
    defer client.Close()
    
    args := &Args{A: 1, B: 2}
    var reply int
    
    // 同步调用
    err := client.Call("Arith.Add", args, &reply)
    
    // 异步调用
    call := client.Go("Arith.Add", args, &reply, nil)
    <-call.Done
}
```

## 协议设计

```
+--------+----------+----------+----------+
| Header |  Method  |   Args   |  Checksum|
+--------+----------+----------+----------+
| 4 bytes| variable | variable | 4 bytes  |
+--------+----------+----------+----------+

Header:
- Magic Number (2 bytes)
- Version (1 byte)
- Serialization Type (1 byte)
```

## 实现提示

### 服务注册
```go
type service struct {
    name    string
    typ     reflect.Type
    val     reflect.Value
    methods map[string]*methodType
}

func (s *Server) Register(rcvr interface{}) error {
    // 使用反射获取方法信息
    // 验证方法签名
    // 注册到 map
}
```

### 客户端调用
```go
func (c *Client) Call(method string, args, reply interface{}) error {
    // 1. 序列化请求
    // 2. 发送请求
    // 3. 接收响应
    // 4. 反序列化响应
}
```

## 测试场景
1. 基本调用
2. 并发调用
3. 超时处理
4. 服务不存在
5. 网络断开重连

## 参考资源
- [net/rpc 标准库](https://golang.org/pkg/net/rpc/)
- [gRPC-Go](https://github.com/grpc/grpc-go)
- [7天用Go从零实现RPC框架](https://geektutu.com/post/geerpc.html)

## 学习收获
- 理解 RPC 的工作原理
- 掌握网络编程和协议设计
- 学会使用反射实现通用框架
- 理解服务治理的概念
