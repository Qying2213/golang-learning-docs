# 计算机网络教程

> 专为 Go 后端开发者设计的计算机网络教程
> 
> 从零开始，结合实际工作场景，让你真正理解网络原理

## 📚 教程目录

### 基础篇
1. [网络基础概念](./01-网络基础概念.md) ⭐⭐⭐
   - 什么是计算机网络
   - 网络分层模型（OSI 和 TCP/IP）
   - IP 地址和子网掩码
   - 域名和 DNS

2. [HTTP 协议详解](./02-HTTP协议详解.md) ⭐⭐⭐
   - HTTP 请求和响应
   - HTTP 方法（GET、POST 等）
   - HTTP 状态码
   - HTTP 头部详解
   - Cookie 和 Session

3. [HTTPS 和加密](./03-HTTPS和加密.md) ⭐⭐⭐
   - 为什么需要 HTTPS
   - 对称加密和非对称加密
   - SSL/TLS 握手过程
   - 证书和 CA

### 传输层
4. [TCP 协议详解](./04-TCP协议详解.md) ⭐⭐⭐
   - TCP 三次握手
   - TCP 四次挥手
   - TCP 可靠传输
   - 流量控制和拥塞控制

5. [UDP 协议](./05-UDP协议.md) ⭐⭐
   - UDP 特点
   - TCP vs UDP
   - 使用场景

### 应用层
6. [RESTful API 设计](./06-RESTful-API设计.md) ⭐⭐⭐
   - REST 原则
   - API 设计最佳实践
   - 版本控制
   - 错误处理

7. [WebSocket 实时通信](./07-WebSocket实时通信.md) ⭐⭐
   - WebSocket 原理
   - 与 HTTP 的区别
   - Go 实现 WebSocket

8. [RPC 和 gRPC](./08-RPC和gRPC.md) ⭐⭐⭐
   - RPC 原理
   - gRPC 入门
   - Protobuf

### 网络层
9. [IP 协议和路由](./09-IP协议和路由.md) ⭐⭐
   - IP 地址分类
   - 子网划分
   - 路由原理
   - NAT 和端口映射

### 实战篇
10. [网络调试工具](./10-网络调试工具.md) ⭐⭐⭐
    - curl 命令详解
    - Postman 使用
    - Wireshark 抓包
    - tcpdump 使用

11. [网络性能优化](./11-网络性能优化.md) ⭐⭐⭐
    - 连接池
    - Keep-Alive
    - HTTP/2 和 HTTP/3
    - CDN 加速

12. [网络安全基础](./12-网络安全基础.md) ⭐⭐⭐
    - 常见攻击（XSS、CSRF、SQL 注入）
    - 防御措施
    - JWT 认证
    - OAuth 2.0

13. [负载均衡和反向代理](./13-负载均衡和反向代理.md) ⭐⭐⭐
    - 负载均衡算法
    - Nginx 反向代理
    - 健康检查
    - 会话保持

## 🎯 学习路线

### 第一周：HTTP 基础
- 01-网络基础概念
- 02-HTTP协议详解
- 10-网络调试工具（curl 部分）

### 第二周：传输层
- 04-TCP协议详解
- 05-UDP协议
- 实践：用 Go 写 TCP 服务器

### 第三周：应用层
- 06-RESTful-API设计
- 07-WebSocket实时通信
- 实践：设计和实现 RESTful API

### 第四周：安全和优化
- 03-HTTPS和加密
- 11-网络性能优化
- 12-网络安全基础

## 💡 学习建议

1. **理论结合实践**：每学完一章，用 Go 写代码验证
2. **抓包分析**：用 Wireshark 看实际的网络包
3. **动手实验**：自己搭建服务器，测试各种场景
4. **结合工作**：把学到的知识应用到实际项目中

## 🔧 准备工作

### 安装工具
```bash
# Mac
brew install curl
brew install wireshark
brew install postman

# Ubuntu
sudo apt install curl
sudo apt install wireshark
sudo snap install postman
```

### Go 网络编程库
```bash
go get github.com/gin-gonic/gin
go get github.com/gorilla/websocket
go get google.golang.org/grpc
```

## 📖 参考资料

- 《计算机网络：自顶向下方法》
- 《HTTP 权威指南》
- 《TCP/IP 详解》
- Go 官方文档：https://golang.org/pkg/net/

## ⭐ 重要程度说明

- ⭐⭐⭐ 必须掌握（面试必考，工作必用）
- ⭐⭐ 重要（工作中常用）
- ⭐ 了解即可（特定场景使用）

---

开始学习吧！从第一章开始，每章都有详细的讲解、图示和练习题！💪
