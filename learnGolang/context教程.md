# Go Context 教程（通俗版）

---

## 1. Context 是什么？一句话解释

**Context 就是一个"遥控器"，用来远程关闭 goroutine。**

---

## 2. 为什么需要 Context？

看这个问题：

```go
func worker() {
    for {
        fmt.Println("我在干活...")
        time.Sleep(1 * time.Second)
    }
}

func main() {
    go worker()  // 启动一个 goroutine
    
    time.Sleep(3 * time.Second)
    
    // 问题：怎么让 worker 停下来？
    // 你喊 "停！" 它听不到啊
}
```

goroutine 一旦启动，就像放出去的风筝，你没法直接控制它。

**Context 就是那根风筝线，让你能把它收回来。**

---

## 3. 最简单的例子：手动取消

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():  // 检查：老板让我停了吗？
            fmt.Println("收到！我停了")
            return
        default:
            fmt.Println("干活中...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // 1. 创建一个遥控器（ctx）和关闭按钮（cancel）
    ctx, cancel := context.WithCancel(context.Background())
    
    // 2. 把遥控器给 worker
    go worker(ctx)
    
    // 3. 等 2 秒
    time.Sleep(2 * time.Second)
    
    // 4. 按下关闭按钮
    cancel()
    
    time.Sleep(1 * time.Second)
    fmt.Println("程序结束")
}
```

运行结果：
```
干活中...
干活中...
干活中...
干活中...
收到！我停了
程序结束
```

### 理解要点

```go
ctx, cancel := context.WithCancel(context.Background())
```
- `ctx` = 遥控器（传给 worker）
- `cancel` = 关闭按钮（你自己留着）
- `context.Background()` = 一个空的起点，别管它

```go
case <-ctx.Done():
```
- worker 不断检查 `ctx.Done()`
- 你按 `cancel()` 后，`ctx.Done()` 就会收到信号
- worker 收到信号就退出

---

## 4. 自动取消：超时控制

不想手动按按钮？可以设置倒计时，时间到了自动取消。

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("超时了，我停了")
            return
        default:
            fmt.Println("干活中...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // 2 秒后自动取消
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()  // 好习惯，先写上
    
    go worker(ctx)
    
    time.Sleep(3 * time.Second)
    fmt.Println("程序结束")
}
```

运行结果：
```
干活中...
干活中...
干活中...
干活中...
超时了，我停了
程序结束
```

**WithTimeout vs WithCancel**：
- `WithCancel` = 手动按按钮
- `WithTimeout` = 设置倒计时，时间到自动按

---

## 5. 实际场景：请求超时

你调用一个 API，如果 2 秒没响应就放弃：

```go
func callAPI(ctx context.Context) (string, error) {
    // 模拟 API 调用需要 3 秒
    select {
    case <-time.After(3 * time.Second):
        return "API 返回的数据", nil
    case <-ctx.Done():
        return "", fmt.Errorf("请求超时")
    }
}

func main() {
    // 只等 2 秒
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    result, err := callAPI(ctx)
    if err != nil {
        fmt.Println("失败:", err)  // 会打印这个
    } else {
        fmt.Println("成功:", result)
    }
}
```

API 要 3 秒，但我们只等 2 秒，所以会超时。

---

## 6. 一次取消多个 goroutine

一个遥控器可以控制多个 worker：

```go
func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("worker %d 停了\n", id)
            return
        default:
            fmt.Printf("worker %d 干活中\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // 启动 3 个 worker，用同一个 ctx
    go worker(ctx, 1)
    go worker(ctx, 2)
    go worker(ctx, 3)
    
    time.Sleep(2 * time.Second)
    
    cancel()  // 一键关闭所有
    
    time.Sleep(1 * time.Second)
}
```

**一个 cancel()，三个 worker 全停。**

---

## 7. Context 传值（了解就行）

可以在 context 里塞点数据：

```go
func worker(ctx context.Context) {
    userID := ctx.Value("userID")
    fmt.Println("用户ID:", userID)
}

func main() {
    ctx := context.Background()
    ctx = context.WithValue(ctx, "userID", 12345)
    
    worker(ctx)  // 用户ID: 12345
}
```

一般用来传请求 ID，方便追踪日志。

---

## 8. 记住这几点就够了

### 8.1 创建 context

```go
// 手动取消
ctx, cancel := context.WithCancel(context.Background())

// 超时取消
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

### 8.2 检查是否取消

```go
select {
case <-ctx.Done():
    // 被取消了，退出
    return
default:
    // 继续干活
}
```

### 8.3 好习惯

```go
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()  // 创建后立刻 defer，别忘了
```

---

## 9. 练习

### 练习 1：手动取消

```go
// 启动一个 goroutine 每秒打印 "hello"
// 5 秒后取消它
func main() {
    // 你来实现
}
```

### 练习 2：超时控制

```go
// 模拟一个耗时 3 秒的操作
// 设置 1 秒超时，打印 "超时了"
func slowWork(ctx context.Context) {
    // 你来实现
}

func main() {
    // 你来实现
}
```

---

## 总结

| 概念 | 比喻 |
|------|------|
| Context | 遥控器 |
| cancel() | 关闭按钮 |
| ctx.Done() | worker 检查"老板让我停了吗" |
| WithCancel | 手动按按钮 |
| WithTimeout | 设置倒计时自动按 |

**核心就一句话：Context 是用来远程关闭 goroutine 的。**
