# HTTP 教程（一）：基础原理

---

## 1. HTTP 是什么

HTTP = HyperText Transfer Protocol（超文本传输协议）

**简单说：浏览器和服务器之间"说话"的规则。**

```
浏览器                          服务器
  |                               |
  |  ---- HTTP 请求 ---->         |
  |       "我要首页"              |
  |                               |
  |  <---- HTTP 响应 ----         |
  |       "给你首页内容"          |
```

---

## 2. HTTP 在网络中的位置

```
应用层    HTTP、HTTPS、FTP、DNS    ← 你写代码在这层
传输层    TCP、UDP
网络层    IP
链路层    以太网、WiFi
```

HTTP 基于 TCP，所以：
1. 先建立 TCP 连接（三次握手）
2. 然后发送 HTTP 请求
3. 服务器返回 HTTP 响应
4. 关闭连接（或保持连接复用）

---

## 3. HTTP 请求格式

一个 HTTP 请求长这样：

```
GET /api/users HTTP/1.1          ← 请求行
Host: example.com                ← 请求头
Content-Type: application/json
Authorization: Bearer xxx

{"name": "秦阳"}                 ← 请求体（可选）
```

### 3.1 请求行

```
GET /api/users HTTP/1.1
 |      |         |
方法   路径     HTTP版本
```

### 3.2 常用请求方法

| 方法 | 用途 | 有请求体 |
|------|------|----------|
| GET | 获取资源 | 否 |
| POST | 创建资源 | 是 |
| PUT | 更新资源（全量） | 是 |
| PATCH | 更新资源（部分） | 是 |
| DELETE | 删除资源 | 否 |

### 3.3 常用请求头

```
Host: example.com                    # 目标主机
Content-Type: application/json       # 请求体格式
Content-Length: 123                  # 请求体长度
Authorization: Bearer token123       # 认证信息
User-Agent: Mozilla/5.0              # 客户端信息
Accept: application/json             # 期望的响应格式
Cookie: session=abc123               # Cookie
```

---

## 4. HTTP 响应格式

```
HTTP/1.1 200 OK                      ← 状态行
Content-Type: application/json       ← 响应头
Content-Length: 27

{"id": 1, "name": "秦阳"}            ← 响应体
```

### 4.1 状态码

| 范围 | 含义 | 常见 |
|------|------|------|
| 1xx | 信息 | 100 Continue |
| 2xx | 成功 | 200 OK, 201 Created, 204 No Content |
| 3xx | 重定向 | 301 永久重定向, 302 临时重定向, 304 未修改 |
| 4xx | 客户端错误 | 400 Bad Request, 401 未授权, 403 禁止, 404 Not Found |
| 5xx | 服务器错误 | 500 Internal Error, 502 Bad Gateway, 503 Service Unavailable |

**常用状态码**：

```
200 OK           - 成功
201 Created      - 创建成功
204 No Content   - 成功但无返回内容
400 Bad Request  - 请求参数错误
401 Unauthorized - 未登录
403 Forbidden    - 没权限
404 Not Found    - 资源不存在
500 Internal Server Error - 服务器内部错误
```

---

## 5. HTTP 底层：TCP 连接

### 5.1 一次完整的 HTTP 请求过程

```
1. DNS 解析：example.com → 93.184.216.34
2. TCP 三次握手：建立连接
3. 发送 HTTP 请求
4. 服务器处理请求
5. 返回 HTTP 响应
6. TCP 四次挥手：关闭连接（或保持连接）
```

### 5.2 TCP 三次握手

```
客户端                    服务器
  |                         |
  | ---- SYN seq=x ---->    |  第一次：客户端发起
  |                         |
  | <-- SYN+ACK seq=y -->   |  第二次：服务器确认
  |      ack=x+1            |
  |                         |
  | ---- ACK ack=y+1 --->   |  第三次：客户端确认
  |                         |
  |    连接建立完成          |
```

**为什么要三次？**
- 确认双方都能收发数据
- 两次不够：服务器不知道客户端能不能收到
- 四次多余：三次已经够了

### 5.3 TCP 四次挥手

```
客户端                    服务器
  |                         |
  | ---- FIN ---->          |  第一次：客户端说"我发完了"
  |                         |
  | <---- ACK ----          |  第二次：服务器说"知道了"
  |                         |
  | <---- FIN ----          |  第三次：服务器说"我也发完了"
  |                         |
  | ---- ACK ---->          |  第四次：客户端说"好的，拜拜"
  |                         |
```

**为什么要四次？**
- 关闭是双向的，每个方向都要关闭
- 服务器收到 FIN 后可能还有数据要发，所以 ACK 和 FIN 分开

---

## 6. HTTP 版本演进

### 6.1 HTTP/1.0

- 每次请求都要新建 TCP 连接
- 请求完就断开
- 效率低

### 6.2 HTTP/1.1（目前最常用）

- **持久连接**：一个 TCP 连接可以发多个请求
- **管道化**：可以连续发请求，不用等响应
- **Host 头**：支持虚拟主机

```
Connection: keep-alive  # 保持连接
```

### 6.3 HTTP/2

- **多路复用**：一个连接同时处理多个请求
- **头部压缩**：减少传输量
- **服务器推送**：服务器主动推送资源
- **二进制格式**：不再是文本

### 6.4 HTTP/3

- 基于 **QUIC**（UDP），不是 TCP
- 更快的连接建立
- 更好的丢包处理

---

## 7. HTTPS

HTTPS = HTTP + TLS/SSL（加密）

```
HTTP:  明文传输，不安全
HTTPS: 加密传输，安全
```

### 7.1 HTTPS 握手过程（简化版）

```
1. 客户端发送支持的加密算法列表
2. 服务器选择算法，发送证书（包含公钥）
3. 客户端验证证书，生成随机密钥，用公钥加密发送
4. 服务器用私钥解密，得到密钥
5. 双方用这个密钥加密通信
```

### 7.2 为什么 HTTPS 更安全

- **加密**：数据被加密，中间人看不懂
- **身份验证**：证书证明服务器身份
- **完整性**：数据不会被篡改

---

## 8. 用 Go 看 HTTP 请求

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    // 发送 GET 请求
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // 打印状态码
    fmt.Println("状态码:", resp.StatusCode)
    
    // 打印响应头
    fmt.Println("响应头:")
    for k, v := range resp.Header {
        fmt.Printf("  %s: %s\n", k, v)
    }
    
    // 打印响应体
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("响应体:", string(body))
}
```

---

## 下一篇

[HTTP教程-2-RESTful API](./HTTP教程-2-RESTful-API.md)
