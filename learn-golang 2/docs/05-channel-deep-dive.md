# Channel 深入理解

## 1. Channel 底层结构

```go
// runtime/chan.go
type hchan struct {
    qcount   uint           // 当前队列中的元素数量
    dataqsiz uint           // 环形队列大小（缓冲区容量）
    buf      unsafe.Pointer // 环形队列指针
    elemsize uint16         // 元素大小
    closed   uint32         // 是否已关闭
    elemtype *_type         // 元素类型
    sendx    uint           // 发送索引
    recvx    uint           // 接收索引
    recvq    waitq          // 等待接收的 goroutine 队列
    sendq    waitq          // 等待发送的 goroutine 队列
    lock     mutex          // 互斥锁
}

type waitq struct {
    first *sudog  // 等待队列头
    last  *sudog  // 等待队列尾
}
```

## 2. Channel 的创建

```go
// 无缓冲 channel（同步）
ch1 := make(chan int)

// 有缓冲 channel（异步）
ch2 := make(chan int, 10)

// 只读 channel
var ch3 <-chan int

// 只写 channel
var ch4 chan<- int

// nil channel
var ch5 chan int  // ch5 == nil
```

## 3. 发送和接收操作

### 无缓冲 Channel
```go
ch := make(chan int)

// 发送会阻塞，直到有接收者
go func() {
    ch <- 42  // 阻塞，直到有人接收
}()

// 接收会阻塞，直到有发送者
val := <-ch  // 阻塞，直到有人发送
```

### 有缓冲 Channel
```go
ch := make(chan int, 3)

// 缓冲区未满时，发送不阻塞
ch <- 1  // 不阻塞
ch <- 2  // 不阻塞
ch <- 3  // 不阻塞
ch <- 4  // 阻塞！缓冲区已满

// 缓冲区非空时，接收不阻塞
<-ch  // 不阻塞
```

## 4. Channel 的状态和操作结果

| 操作 | nil channel | 已关闭 channel | 正常 channel |
|------|-------------|----------------|--------------|
| 发送 | 永久阻塞 | panic | 阻塞或成功 |
| 接收 | 永久阻塞 | 返回零值 | 阻塞或成功 |
| 关闭 | panic | panic | 成功 |

```go
// 检查 channel 是否关闭
val, ok := <-ch
if !ok {
    // channel 已关闭
}

// 使用 range 遍历（自动检测关闭）
for val := range ch {
    fmt.Println(val)
}
// channel 关闭后自动退出循环
```

## 5. Select 多路复用

```go
select {
case val := <-ch1:
    fmt.Println("received from ch1:", val)
case ch2 <- 42:
    fmt.Println("sent to ch2")
case <-time.After(time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("no communication")
}
```

### Select 的特性
```go
// 1. 随机选择就绪的 case
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 2

select {
case <-ch1:
    fmt.Println("ch1")  // 可能
case <-ch2:
    fmt.Println("ch2")  // 可能
}

// 2. 没有 default 时会阻塞
select {
case <-ch1:
    // ...
case <-ch2:
    // ...
}
// 如果 ch1 和 ch2 都没有数据，会阻塞

// 3. 空 select 永久阻塞
select {}  // 永久阻塞，常用于阻止 main 退出
```

## 6. 常见模式

### 模式1: 信号通知
```go
done := make(chan struct{})

go func() {
    // 做一些工作...
    close(done)  // 通知完成
}()

<-done  // 等待完成
```

### 模式2: 超时控制
```go
select {
case result := <-ch:
    fmt.Println(result)
case <-time.After(3 * time.Second):
    fmt.Println("timeout")
}
```

### 模式3: 取消操作
```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return  // 收到取消信号
        default:
            // 做工作...
        }
    }
}

ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)
// ...
cancel()  // 取消
```

### 模式4: 扇出/扇入
```go
// 扇出：一个输入，多个输出
func fanOut(in <-chan int, n int) []<-chan int {
    outs := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        out := make(chan int)
        outs[i] = out
        go func() {
            for val := range in {
                out <- val * 2
            }
            close(out)
        }()
    }
    return outs
}

// 扇入：多个输入，一个输出
func fanIn(ins ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, in := range ins {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for val := range ch {
                out <- val
            }
        }(in)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### 模式5: 信号量
```go
// 使用带缓冲的 channel 作为信号量
sem := make(chan struct{}, 3)  // 最多 3 个并发

for i := 0; i < 10; i++ {
    sem <- struct{}{}  // 获取信号量
    go func(id int) {
        defer func() { <-sem }()  // 释放信号量
        // 做工作...
    }(i)
}
```

### 模式6: 生产者-消费者
```go
func producer(ch chan<- int) {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)
}

func consumer(ch <-chan int, done chan<- struct{}) {
    for val := range ch {
        fmt.Println(val)
    }
    done <- struct{}{}
}

func main() {
    ch := make(chan int, 5)
    done := make(chan struct{})
    
    go producer(ch)
    go consumer(ch, done)
    
    <-done
}
```

## 7. Channel 的关闭

```go
// 只有发送方应该关闭 channel
// 接收方不应该关闭

// 关闭后：
// - 发送会 panic
// - 接收会返回零值和 false
// - 再次关闭会 panic

ch := make(chan int, 3)
ch <- 1
ch <- 2
close(ch)

// 可以继续接收已有数据
fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
fmt.Println(<-ch)  // 0 (零值)

// 检查是否关闭
val, ok := <-ch
if !ok {
    fmt.Println("channel closed")
}
```

## 8. 常见陷阱

```go
// 陷阱1: 向已关闭的 channel 发送
ch := make(chan int)
close(ch)
ch <- 1  // panic!

// 陷阱2: 重复关闭
ch := make(chan int)
close(ch)
close(ch)  // panic!

// 陷阱3: nil channel 操作
var ch chan int
ch <- 1   // 永久阻塞
<-ch      // 永久阻塞
close(ch) // panic!

// 陷阱4: 在 select 中忘记 default 导致阻塞
select {
case <-ch:
    // 如果 ch 没有数据，会永久阻塞
}

// 陷阱5: goroutine 泄漏
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // 永远阻塞，goroutine 泄漏
        fmt.Println(val)
    }()
    // 函数返回，但 goroutine 还在等待
}
```

## 9. 性能考虑

```go
// 1. 有缓冲 vs 无缓冲
// 有缓冲可以减少阻塞，提高吞吐量
// 但会增加内存使用和延迟

// 2. channel 大小选择
// 太小：频繁阻塞
// 太大：内存浪费，延迟增加

// 3. 避免在热路径上使用 channel
// channel 操作有锁开销
// 高性能场景考虑使用 atomic 或 sync.Pool

// 4. 批量操作
// 不好：每个元素一次发送
for _, item := range items {
    ch <- item
}

// 好：批量发送
ch <- items  // 发送整个切片
```

## 10. 面试题

**Q1: 无缓冲和有缓冲 channel 的区别？**
```
无缓冲：同步通信，发送和接收必须同时就绪
有缓冲：异步通信，缓冲区未满时发送不阻塞
```

**Q2: 如何优雅地关闭 channel？**
```go
// 方法1: 使用 sync.Once
type SafeChannel struct {
    ch   chan int
    once sync.Once
}

func (sc *SafeChannel) Close() {
    sc.once.Do(func() {
        close(sc.ch)
    })
}

// 方法2: 使用 context
ctx, cancel := context.WithCancel(context.Background())
// 用 ctx.Done() 代替关闭 channel
```

**Q3: 如何实现一个带超时的发送？**
```go
func sendWithTimeout(ch chan<- int, val int, timeout time.Duration) bool {
    select {
    case ch <- val:
        return true
    case <-time.After(timeout):
        return false
    }
}
```

## 练习

1. 实现一个 pipeline：生成数字 -> 平方 -> 过滤偶数 -> 打印
2. 实现一个带超时的 worker pool
3. 实现一个 rate limiter（限流器）
4. 实现一个 pub/sub 系统
