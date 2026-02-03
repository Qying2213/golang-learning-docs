# WebSocket 实时通信

> **重要程度：⭐⭐ 重要**
> 实时通信是现代应用的标配，聊天、通知、实时数据都需要它！

## 📚 本章学习目标

学完本章，你将能够：
- 理解 WebSocket 的原理
- 掌握 WebSocket 与 HTTP 的区别
- 用 Go 实现 WebSocket 服务
- 实现简单的聊天室功能

---

## 1. 为什么需要 WebSocket？⭐⭐⭐

### 1.1 HTTP 的局限性

```
HTTP 是请求-响应模式：

客户端 → 请求 → 服务器
客户端 ← 响应 ← 服务器

问题：服务器不能主动推送消息给客户端！

场景：聊天应用
- 张三发消息给李四
- 服务器收到消息
- 但服务器无法主动通知李四！
- 李四只能不断轮询：有新消息吗？有新消息吗？
```

### 1.2 轮询的问题

```
短轮询：
客户端每隔几秒请求一次
GET /messages → 没有新消息
GET /messages → 没有新消息
GET /messages → 有新消息！

问题：
- 浪费带宽（大量无效请求）
- 浪费服务器资源
- 延迟高（最多等一个轮询周期）

长轮询：
客户端发请求，服务器 hold 住
有消息时才返回，然后客户端再发请求

问题：
- 连接频繁建立断开
- 服务器资源占用
- 实现复杂
```

### 1.3 WebSocket 的解决方案

```
WebSocket = 全双工通信

建立连接后，双方可以随时发送消息：

客户端 ←→ 服务器

特点：
1. 持久连接（不用反复建立）
2. 全双工（双向通信）
3. 低延迟（实时推送）
4. 低开销（头部只有 2-10 字节）
```

---

## 2. WebSocket 原理 ⭐⭐⭐

### 2.1 握手过程

```
WebSocket 通过 HTTP 升级建立连接：

1. 客户端发送升级请求
GET /chat HTTP/1.1
Host: example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13

2. 服务器同意升级
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=

3. 连接建立，开始 WebSocket 通信
```

### 2.2 数据帧格式

```
WebSocket 数据帧：

 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+

opcode 类型：
0x1 = 文本帧
0x2 = 二进制帧
0x8 = 关闭连接
0x9 = Ping
0xA = Pong
```

### 2.3 WebSocket vs HTTP

| 特性 | HTTP | WebSocket |
|------|------|-----------|
| 通信模式 | 请求-响应 | 全双工 |
| 连接 | 短连接/长连接 | 持久连接 |
| 服务器推送 | 不支持 | 支持 |
| 头部开销 | 大（几百字节） | 小（2-10字节） |
| 协议 | http:// | ws:// |
| 加密协议 | https:// | wss:// |

---

## 3. Go 实现 WebSocket ⭐⭐⭐

### 3.1 安装依赖

```bash
go get github.com/gorilla/websocket
```

### 3.2 简单的 WebSocket 服务器

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // 允许跨域（生产环境要限制）
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    // 升级 HTTP 连接为 WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("升级失败:", err)
        return
    }
    defer conn.Close()

    log.Println("客户端连接:", conn.RemoteAddr())

    for {
        // 读取消息
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("读取错误:", err)
            break
        }

        log.Printf("收到消息: %s", message)

        // 回复消息
        response := fmt.Sprintf("服务器收到: %s", message)
        err = conn.WriteMessage(messageType, []byte(response))
        if err != nil {
            log.Println("发送错误:", err)
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", wsHandler)
    
    log.Println("WebSocket 服务器启动: ws://localhost:8080/ws")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 3.3 HTML 客户端测试

```html
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket 测试</title>
</head>
<body>
    <h1>WebSocket 测试</h1>
    <input type="text" id="message" placeholder="输入消息">
    <button onclick="sendMessage()">发送</button>
    <div id="output"></div>

    <script>
        const ws = new WebSocket('ws://localhost:8080/ws');
        const output = document.getElementById('output');

        ws.onopen = function() {
            output.innerHTML += '<p>连接已建立</p>';
        };

        ws.onmessage = function(event) {
            output.innerHTML += '<p>收到: ' + event.data + '</p>';
        };

        ws.onclose = function() {
            output.innerHTML += '<p>连接已关闭</p>';
        };

        ws.onerror = function(error) {
            output.innerHTML += '<p>错误: ' + error + '</p>';
        };

        function sendMessage() {
            const message = document.getElementById('message').value;
            ws.send(message);
            output.innerHTML += '<p>发送: ' + message + '</p>';
        }
    </script>
</body>
</html>
```

---

## 4. 实现聊天室 ⭐⭐⭐

### 4.1 聊天室架构

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│ 客户端A  │     │ 客户端B  │     │ 客户端C  │
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
              ┌──────┴──────┐
              │   服务器     │
              │  (Hub)      │
              │             │
              │ 管理所有连接 │
              │ 广播消息     │
              └─────────────┘
```

### 4.2 完整聊天室代码

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

// 消息结构
type Message struct {
    Type     string `json:"type"`     // message, join, leave
    Username string `json:"username"`
    Content  string `json:"content"`
    Time     string `json:"time"`
}

// 客户端连接
type Client struct {
    conn     *websocket.Conn
    username string
    send     chan []byte
}

// 聊天室 Hub
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

var hub = &Hub{
    clients:    make(map[*Client]bool),
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// Hub 运行
func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            
            // 广播加入消息
            msg := Message{
                Type:     "join",
                Username: client.username,
                Content:  client.username + " 加入了聊天室",
                Time:     time.Now().Format("15:04:05"),
            }
            data, _ := json.Marshal(msg)
            h.broadcast <- data

        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mutex.Unlock()
            
            // 广播离开消息
            msg := Message{
                Type:     "leave",
                Username: client.username,
                Content:  client.username + " 离开了聊天室",
                Time:     time.Now().Format("15:04:05"),
            }
            data, _ := json.Marshal(msg)
            h.broadcast <- data

        case message := <-h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
        }
    }
}

// 读取消息
func (c *Client) readPump() {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        msg := Message{
            Type:     "message",
            Username: c.username,
            Content:  string(message),
            Time:     time.Now().Format("15:04:05"),
        }
        data, _ := json.Marshal(msg)
        hub.broadcast <- data
    }
}

// 发送消息
func (c *Client) writePump() {
    defer c.conn.Close()

    for message := range c.send {
        err := c.conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}

// WebSocket 处理
func wsHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        username = "匿名用户"
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("升级失败:", err)
        return
    }

    client := &Client{
        conn:     conn,
        username: username,
        send:     make(chan []byte, 256),
    }

    hub.register <- client

    go client.writePump()
    go client.readPump()
}

func main() {
    go hub.run()

    http.HandleFunc("/ws", wsHandler)
    http.Handle("/", http.FileServer(http.Dir("./static")))

    log.Println("聊天室启动: http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 4.3 聊天室前端

```html
<!-- static/index.html -->
<!DOCTYPE html>
<html>
<head>
    <title>聊天室</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        #messages { height: 400px; border: 1px solid #ccc; overflow-y: scroll; padding: 10px; margin-bottom: 10px; }
        .message { margin: 5px 0; padding: 5px; background: #f0f0f0; border-radius: 5px; }
        .join { color: green; }
        .leave { color: red; }
        .time { color: #999; font-size: 12px; }
        input[type="text"] { width: 70%; padding: 10px; }
        button { padding: 10px 20px; }
    </style>
</head>
<body>
    <h1>聊天室</h1>
    
    <div id="login">
        <input type="text" id="username" placeholder="输入用户名">
        <button onclick="connect()">进入聊天室</button>
    </div>
    
    <div id="chat" style="display:none;">
        <div id="messages"></div>
        <input type="text" id="message" placeholder="输入消息" onkeypress="if(event.keyCode==13)sendMessage()">
        <button onclick="sendMessage()">发送</button>
    </div>

    <script>
        let ws;
        const messages = document.getElementById('messages');

        function connect() {
            const username = document.getElementById('username').value || '匿名';
            ws = new WebSocket('ws://localhost:8080/ws?username=' + encodeURIComponent(username));

            ws.onopen = function() {
                document.getElementById('login').style.display = 'none';
                document.getElementById('chat').style.display = 'block';
            };

            ws.onmessage = function(event) {
                const msg = JSON.parse(event.data);
                let className = 'message';
                if (msg.type === 'join') className += ' join';
                if (msg.type === 'leave') className += ' leave';
                
                messages.innerHTML += `
                    <div class="${className}">
                        <span class="time">[${msg.time}]</span>
                        <strong>${msg.username}:</strong> ${msg.content}
                    </div>
                `;
                messages.scrollTop = messages.scrollHeight;
            };

            ws.onclose = function() {
                messages.innerHTML += '<div class="message leave">连接已断开</div>';
            };
        }

        function sendMessage() {
            const input = document.getElementById('message');
            if (input.value && ws) {
                ws.send(input.value);
                input.value = '';
            }
        }
    </script>
</body>
</html>
```

---

## 5. 心跳检测 ⭐⭐

### 5.1 为什么需要心跳？

```
问题：
- 网络断开时，连接可能不会立即关闭
- 服务器不知道客户端是否还在
- 资源无法及时释放

解决：心跳检测
- 定期发送 Ping/Pong
- 超时未响应则断开连接
```

### 5.2 实现心跳

```go
const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

func (c *Client) readPump() {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        // 处理消息...
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.conn.WriteMessage(websocket.TextMessage, message)

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

---

## 6. 本章总结

### 必须记住的内容

```
WebSocket 特点：
- 全双工通信
- 持久连接
- 低延迟
- 服务器可主动推送

与 HTTP 区别：
- HTTP：请求-响应
- WebSocket：双向通信

使用场景：
- 聊天应用
- 实时通知
- 在线游戏
- 实时数据展示

Go 实现：
- gorilla/websocket 库
- Upgrader 升级连接
- ReadMessage/WriteMessage
```

### 面试常见问题

```
Q: WebSocket 和 HTTP 的区别？
A: HTTP 是请求-响应模式，WebSocket 是全双工通信，服务器可以主动推送。

Q: WebSocket 如何建立连接？
A: 通过 HTTP 升级请求，服务器返回 101 状态码后切换到 WebSocket 协议。

Q: 如何保持 WebSocket 连接？
A: 使用心跳机制，定期发送 Ping/Pong 消息。

Q: WebSocket 适用于什么场景？
A: 需要实时双向通信的场景，如聊天、通知、实时数据。
```

---

下一章：[08-RPC和gRPC](./08-RPC和gRPC.md)
