# UDP 协议

> **重要程度：⭐⭐ 重要**
> UDP 虽然简单，但在特定场景下比 TCP 更合适！

## 📚 本章学习目标

学完本章，你将能够：

- 理解 UDP 的特点
- 掌握 TCP 和 UDP 的区别
- 知道什么场景用 UDP
- 用 Go 实现 UDP 通信

---

## 1. UDP 是什么？⭐⭐⭐

### 1.1 UDP 的全称

```
UDP = User Datagram Protocol
      用户数据报协议

特点：简单、快速、不可靠
```

### 1.2 UDP 的特点

```
1. 无连接
   - 不需要建立连接
   - 直接发送数据
   - 就像寄信，不用先打电话

2. 不可靠
   - 数据可能丢失
   - 数据可能乱序
   - 数据可能重复
   - 不保证送达

3. 面向数据报
   - 保留消息边界
   - 一次发送 = 一次接收
   - 不会像 TCP 那样拆分合并

4. 无拥塞控制
   - 想发多快就多快
   - 不管网络是否拥堵

5. 支持广播和多播
   - 可以一对多发送
   - TCP 只能一对一
```

### 1.3 UDP 数据报格式

```
┌─────────────────────────────────────────────────┐
│                  UDP 数据报                      │
├──────────────┬──────────────┬──────────────────┤
│   源端口      │   目标端口    │                  │
│   16 位       │   16 位      │                  │
├──────────────┼──────────────┤     UDP 头部     │
│   长度        │   校验和      │     8 字节       │
│   16 位       │   16 位      │                  │
├──────────────┴──────────────┴──────────────────┤
│                                                 │
│                   数据                          │
│                                                 │
└─────────────────────────────────────────────────┘

UDP 头部只有 8 字节，非常简洁！
TCP 头部至少 20 字节
```

---

## 2. TCP vs UDP 详细对比 ⭐⭐⭐

### 2.1 对比表

| 特性      | TCP            | UDP              |
| --------- | -------------- | ---------------- |
| 连接      | 需要三次握手   | 无需连接         |
| 可靠性    | 可靠，保证送达 | 不可靠，可能丢失 |
| 顺序      | 保证顺序       | 不保证顺序       |
| 速度      | 较慢           | 较快             |
| 头部大小  | 20+ 字节       | 8 字节           |
| 流量控制  | 有             | 无               |
| 拥塞控制  | 有             | 无               |
| 传输方式  | 字节流         | 数据报           |
| 广播/多播 | 不支持         | 支持             |
| 适用场景  | 可靠传输       | 实时性要求高     |

### 2.2 形象比喻

```
TCP 像打电话：
1. 先拨号建立连接
2. 确认对方接听
3. 开始通话
4. 说完挂断
- 可靠，但有延迟

UDP 像发短信：
1. 直接发送
2. 不知道对方收没收到
3. 可能丢失
- 快速，但不可靠

TCP 像快递：
- 签收确认
- 保证送达
- 有追踪

UDP 像平信：
- 投进邮筒就完事
- 可能丢失
- 没有追踪
```

### 2.3 什么时候用 TCP？

```
需要可靠传输的场景：

1. HTTP/HTTPS
   - 网页必须完整加载
   - 不能丢数据

2. 文件传输（FTP）
   - 文件不能损坏
   - 必须完整

3. 邮件（SMTP）
   - 邮件内容不能丢

4. 数据库连接
   - 数据必须准确

5. SSH 远程登录
   - 命令必须准确执行
```

### 2.4 什么时候用 UDP？

```
实时性要求高、允许丢包的场景：

1. 视频直播/视频通话
   - 丢几帧没关系
   - 延迟低更重要
   - 丢包了也不重传（重传的已经过时了）

2. 在线游戏
   - 位置更新要快
   - 丢一个位置包没关系
   - 下一个包马上就来

3. DNS 查询
   - 数据量小（通常 < 512 字节）
   - 查询简单
   - 丢了重发就行
   - 注意：超过 512 字节或区域传送时使用 TCP

4. 物联网/传感器数据
   - 数据量大
   - 允许丢失部分数据
   - 实时性重要

5. 广播/多播
   - 一对多发送
   - TCP 不支持
```

---

## 3. Go 中使用 UDP ⭐⭐⭐

### 3.1 UDP 服务器

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // 监听 UDP 端口
    addr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        panic(err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    fmt.Println("UDP 服务器启动，监听 :8080")

    buffer := make([]byte, 1024)
    for {
        // 接收数据
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("读取错误:", err)
            continue
        }

        message := string(buffer[:n])
        fmt.Printf("收到来自 %s 的消息: %s\n", clientAddr, message)

        // 回复数据
        response := []byte("服务器收到: " + message)
        conn.WriteToUDP(response, clientAddr)
    }
}
```

### 3.2 UDP 客户端

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    // 解析服务器地址
    serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        panic(err)
    }

    // 创建 UDP 连接
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    fmt.Println("UDP 客户端启动")

    reader := bufio.NewReader(os.Stdin)
    buffer := make([]byte, 1024)

    for {
        fmt.Print("请输入消息: ")
        input, _ := reader.ReadString('\n')

        // 发送数据
        conn.Write([]byte(input))

        // 接收响应
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("读取错误:", err)
            continue
        }

        fmt.Printf("服务器响应: %s\n", string(buffer[:n]))
    }
}
```

### 3.3 UDP 广播

```go
package main

import (
    "fmt"
    "net"
    "time"
)

// 广播发送端
func broadcast() {
    addr, _ := net.ResolveUDPAddr("udp", "255.255.255.255:8080")
    conn, _ := net.DialUDP("udp", nil, addr)
    defer conn.Close()

    for {
        message := fmt.Sprintf("广播消息: %s", time.Now().Format("15:04:05"))
        conn.Write([]byte(message))
        fmt.Println("发送广播:", message)
        time.Sleep(2 * time.Second)
    }
}

// 广播接收端
func receive() {
    addr, _ := net.ResolveUDPAddr("udp", ":8080")
    conn, _ := net.ListenUDP("udp", addr)
    defer conn.Close()

    buffer := make([]byte, 1024)
    for {
        n, remoteAddr, _ := conn.ReadFromUDP(buffer)
        fmt.Printf("收到来自 %s 的广播: %s\n", remoteAddr, string(buffer[:n]))
    }
}
```

---

## 4. UDP 的可靠性改进 ⭐⭐

### 4.1 应用层实现可靠性

```
UDP 本身不可靠，但可以在应用层实现：

1. 确认机制
   - 接收方收到后发送 ACK
   - 发送方没收到 ACK 就重发

2. 序列号
   - 给每个包编号
   - 接收方按序号排序

3. 超时重传
   - 设置超时时间
   - 超时未收到 ACK 就重发

4. 校验和
   - 检测数据是否损坏
```

### 4.2 简单的可靠 UDP 示例

```go
package main

import (
    "encoding/binary"
    "fmt"
    "net"
    "time"
)

// 带序列号的消息
type Message struct {
    SeqNum uint32
    Data   []byte
}

// 发送并等待确认
func sendWithAck(conn *net.UDPConn, addr *net.UDPAddr, msg Message) error {
    // 序列化消息
    buf := make([]byte, 4+len(msg.Data))
    binary.BigEndian.PutUint32(buf[:4], msg.SeqNum)
    copy(buf[4:], msg.Data)

    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        // 发送
        conn.WriteToUDP(buf, addr)

        // 等待 ACK
        conn.SetReadDeadline(time.Now().Add(time.Second))
        ackBuf := make([]byte, 4)
        _, _, err := conn.ReadFromUDP(ackBuf)
        if err == nil {
            ackSeq := binary.BigEndian.Uint32(ackBuf)
            if ackSeq == msg.SeqNum {
                return nil // 收到正确的 ACK
            }
        }
        fmt.Printf("重试 %d/%d\n", i+1, maxRetries)
    }
    return fmt.Errorf("发送失败，超过最大重试次数")
}
```

### 4.3 现有的可靠 UDP 协议

```
1. QUIC（Quick UDP Internet Connections）
   - Google 开发
   - HTTP/3 的基础
   - 在 UDP 上实现了类似 TCP 的可靠性
   - 比 TCP 更快（0-RTT 连接）

2. KCP
   - 快速可靠的 ARQ 协议
   - 比 TCP 延迟低 30%-40%
   - 游戏常用

3. UDT（UDP-based Data Transfer）
   - 高性能数据传输
   - 适合大文件传输

4. RUDP（Reliable UDP）
   - 在 UDP 上实现可靠传输
   - 各种实现版本
```

---

## 5. 练习题

### 练习1：选择协议

**问题**：以下场景应该用 TCP 还是 UDP？

1. 在线视频直播
2. 网页浏览
3. 在线游戏（FPS）
4. 文件下载
5. DNS 查询
6. 即时通讯（文字消息）

<details>
<summary>点击查看答案</summary>

```
1. 在线视频直播 → UDP
   - 实时性重要
   - 丢几帧可以接受

2. 网页浏览 → TCP
   - 内容必须完整
   - 不能丢数据

3. 在线游戏（FPS）→ UDP
   - 延迟要低
   - 位置更新要快
   - 丢一个包没关系

4. 文件下载 → TCP
   - 文件必须完整
   - 不能丢数据

5. DNS 查询 → UDP
   - 数据量小
   - 简单查询
   - 丢了重发

6. 即时通讯（文字消息）→ TCP
   - 消息不能丢
   - 顺序要正确
```

</details>

### 练习2：理解差异

**问题**：为什么视频通话用 UDP 而不是 TCP？

<details>
<summary>点击查看答案</summary>

```
原因：

1. 实时性要求
   - 视频通话要求低延迟
   - TCP 的重传机制会增加延迟
   - 等重传的数据到达时，已经过时了

2. 丢包处理
   - TCP：丢包后重传，等待数据到齐
   - UDP：丢包就丢了，继续播放新的帧
   - 用户体验：卡顿 vs 花屏
   - 花屏一瞬间，卡顿很难受

3. 拥塞控制
   - TCP 遇到拥塞会降速
   - 视频通话需要稳定的带宽
   - UDP 可以自己控制发送速率

4. 实际做法
   - 使用 UDP 传输
   - 应用层做简单的丢包补偿
   - 比如：前向纠错（FEC）
```

</details>

---

## 6. 本章总结

### 必须记住的内容

```
UDP 特点：
- 无连接
- 不可靠
- 快速
- 支持广播

TCP vs UDP：
- TCP：可靠、有序、慢
- UDP：不可靠、无序、快

UDP 使用场景：
- 视频/音频
- 游戏
- DNS
- 广播

Go UDP 编程：
- net.ListenUDP()
- conn.ReadFromUDP()
- conn.WriteToUDP()
```

### 面试常见问题

```
Q: TCP 和 UDP 的区别？
A: TCP 面向连接、可靠、有序；UDP 无连接、不可靠、快速。

Q: 什么场景用 UDP？
A: 实时性要求高、允许丢包的场景，如视频、游戏、DNS。

Q: UDP 如何实现可靠传输？
A: 应用层实现确认、重传、序列号机制，或使用 QUIC、KCP 等协议。

Q: 为什么 DNS 用 UDP？
A: 查询数据量小，一个包就够；简单快速；丢了重发即可。
```

---

下一章：[06-RESTful-API设计](./06-RESTful-API设计.md)
