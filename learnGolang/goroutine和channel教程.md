# Go Goroutine 和 Channel 完全教程

> 从零开始理解并发编程，包含底层原理

---

## 一、什么是 Goroutine

**Goroutine 是 Go 的轻量级线程，由 Go 运行时管理。**

```go
// 启动一个 goroutine，只需要在函数前加 go
go func() {
    fmt.Println("我在另一个 goroutine 里")
}()
```

**Goroutine vs 线程：**

| 对比 | Goroutine | 线程 |
|------|-----------|------|
| 创建成本 | 约 2KB 栈内存 | 约 1MB 栈内存 |
| 切换成本 | 几十纳秒 | 几微秒 |
| 数量 | 可以创建几十万个 | 几千个就很多了 |
| 调度 | Go 运行时调度 | 操作系统调度 |

---

## 二、Goroutine 基本使用

### 2.1 启动 Goroutine

```go
func sayHello(name string) {
    fmt.Println("Hello", name)
}

func main() {
    go sayHello("秦阳")  // 启动一个 goroutine
    go sayHello("小明")  // 再启动一个
    
    time.Sleep(time.Second)  // 等待 goroutine 执行
    // 如果不等待，main 退出，所有 goroutine 都会被杀死
}
```

### 2.2 匿名函数

```go
func main() {
    go func() {
        fmt.Println("匿名函数 goroutine")
    }()
    
    // 带参数
    go func(msg string) {
        fmt.Println(msg)
    }("hello")
    
    time.Sleep(time.Second)
}
```

### 2.3 常见错误：闭包陷阱

```go
// ❌ 错误写法
func main() {
    for i := 0; i < 3; i++ {
        go func() {
            fmt.Println(i)  // 可能打印 3, 3, 3
        }()
    }
    time.Sleep(time.Second)
}

// ✅ 正确写法：传参
func main() {
    for i := 0; i < 3; i++ {
        go func(n int) {
            fmt.Println(n)  // 打印 0, 1, 2（顺序不定）
        }(i)
    }
    time.Sleep(time.Second)
}
```

---

## 三、GMP 调度模型（底层原理）

### 3.1 GMP 是什么

```
G = Goroutine（协程）
M = Machine（操作系统线程）
P = Processor（处理器，逻辑CPU）

┌─────────────────────────────────────────────────────────────┐
│                        GMP 模型                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│    ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐                          │
│    │ G │ │ G │ │ G │ │ G │ │ G │  ... 很多 Goroutine       │
│    └─┬─┘ └─┬─┘ └─┬─┘ └─┬─┘ └─┬─┘                          │
│      │     │     │     │     │                             │
│      └──┬──┴──┬──┘     └──┬──┘                             │
│         │     │           │                                │
│         ▼     ▼           ▼                                │
│       ┌───┐ ┌───┐       ┌───┐                              │
│       │ P │ │ P │       │ P │    P 的数量 = CPU 核心数      │
│       └─┬─┘ └─┬─┘       └─┬─┘                              │
│         │     │           │                                │
│         ▼     ▼           ▼                                │
│       ┌───┐ ┌───┐       ┌───┐                              │
│       │ M │ │ M │       │ M │    M = 操作系统线程           │
│       └───┘ └───┘       └───┘                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 各部分职责

**G（Goroutine）：**
```go
// runtime/runtime2.go
type g struct {
    stack       stack   // 栈内存
    sched       gobuf   // 保存执行上下文（PC、SP等）
    goid        int64   // goroutine ID
    status      uint32  // 状态（运行、等待、死亡等）
    // ...
}
```

**P（Processor）：**
```go
type p struct {
    id          int32
    status      uint32
    runqhead    uint32     // 本地队列头
    runqtail    uint32     // 本地队列尾
    runq        [256]guintptr  // 本地运行队列，最多256个G
    runnext     guintptr   // 下一个要运行的G
    // ...
}
```

**M（Machine）：**
```go
type m struct {
    g0      *g     // 调度用的特殊 goroutine
    curg    *g     // 当前运行的 goroutine
    p       puintptr  // 绑定的 P
    // ...
}
```

### 3.3 调度流程

```
1. 创建 G
   go func() {...}
        │
        ▼
2. 放入 P 的本地队列
   ┌─────────────────┐
   │ P 的本地队列     │
   │ [G1][G2][G3]... │
   └─────────────────┘
        │
        ▼
3. M 从 P 获取 G 执行
   M ◄─── P ◄─── G
        │
        ▼
4. G 执行完或被阻塞
   - 执行完：G 销毁
   - 阻塞：G 让出，M 执行下一个 G
```

### 3.4 P 的本地队列和全局队列

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│                    全局队列（Global Queue）                  │
│                    [G][G][G][G][G]...                       │
│                           ▲                                 │
│                           │ 本地队列满了，放全局             │
│         ┌─────────────────┼─────────────────┐              │
│         │                 │                 │              │
│         ▼                 ▼                 ▼              │
│    ┌─────────┐       ┌─────────┐       ┌─────────┐        │
│    │   P0    │       │   P1    │       │   P2    │        │
│    │ [G][G]  │       │ [G][G]  │       │ [G]     │        │
│    │ 本地队列 │       │ 本地队列 │       │ 本地队列 │        │
│    └─────────┘       └─────────┘       └─────────┘        │
│         │                 │                 │              │
│         ▼                 ▼                 ▼              │
│        M0                M1                M2              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3.5 Work Stealing（工作窃取）

```
当 P 的本地队列空了：

1. 先从全局队列拿
2. 全局队列也空了，从其他 P 偷一半

P0 队列空了          P1 队列有很多 G
┌─────────┐         ┌─────────┐
│   P0    │ ◄────── │   P1    │
│ [空]    │  偷一半  │ [G][G][G][G] │
└─────────┘         └─────────┘
```

### 3.6 什么时候会切换 G

```
1. 主动让出
   - runtime.Gosched()
   
2. 系统调用
   - 文件IO、网络IO等
   
3. 阻塞操作
   - channel 操作
   - 锁操作
   - time.Sleep
   
4. 抢占
   - 运行超过 10ms，会被抢占
```

---

## 四、Channel 基本使用

### 4.1 什么是 Channel

**Channel 是 goroutine 之间通信的管道。**

```
Goroutine A ──── 数据 ────► Channel ──── 数据 ────► Goroutine B
```

### 4.2 创建 Channel

```go
// 无缓冲 channel
ch := make(chan int)

// 有缓冲 channel（可以存 3 个元素）
ch := make(chan int, 3)

// 指定类型
ch1 := make(chan string)
ch2 := make(chan bool)
ch3 := make(chan []int)
```

### 4.3 发送和接收

```go
ch := make(chan int)

// 发送
ch <- 42

// 接收
v := <-ch

// 接收并丢弃
<-ch
```

### 4.4 无缓冲 vs 有缓冲

**无缓冲（同步）：**
```go
ch := make(chan int)  // 无缓冲

// 发送会阻塞，直到有人接收
go func() {
    ch <- 42  // 阻塞，等待接收
}()

v := <-ch  // 接收，发送方才能继续
```

**有缓冲（异步）：**
```go
ch := make(chan int, 3)  // 缓冲区大小 3

ch <- 1  // 不阻塞
ch <- 2  // 不阻塞
ch <- 3  // 不阻塞
ch <- 4  // 阻塞！缓冲区满了
```

```
无缓冲：
发送方 ──► 必须等接收方准备好

有缓冲：
发送方 ──► [缓冲区] ──► 接收方
           缓冲区没满就不阻塞
```

### 4.5 关闭 Channel

```go
ch := make(chan int)

close(ch)  // 关闭

// 关闭后：
ch <- 1     // ❌ panic: send on closed channel
v := <-ch   // ✅ 返回零值
v, ok := <-ch  // ✅ ok=false 表示已关闭
```

### 4.6 遍历 Channel

```go
ch := make(chan int)

go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // 必须关闭，否则 range 会一直等待
}()

// range 遍历，直到 channel 关闭
for v := range ch {
    fmt.Println(v)
}
```

### 4.7 单向 Channel

```go
// 只能发送
func send(ch chan<- int) {
    ch <- 42
    // <-ch  // ❌ 编译错误
}

// 只能接收
func receive(ch <-chan int) {
    v := <-ch
    // ch <- 1  // ❌ 编译错误
}

func main() {
    ch := make(chan int)
    go send(ch)
    go receive(ch)
}
```

---

## 五、Channel 底层原理

### 5.1 hchan 结构

```go
// runtime/chan.go
type hchan struct {
    qcount   uint           // 当前元素个数
    dataqsiz uint           // 缓冲区大小
    buf      unsafe.Pointer // 缓冲区指针（环形队列）
    elemsize uint16         // 元素大小
    closed   uint32         // 是否关闭
    sendx    uint           // 发送索引
    recvx    uint           // 接收索引
    recvq    waitq          // 等待接收的 goroutine 队列
    sendq    waitq          // 等待发送的 goroutine 队列
    lock     mutex          // 锁
}
```

### 5.2 结构图

```
┌─────────────────────────────────────────────────────────────┐
│                         hchan                               │
├─────────────────────────────────────────────────────────────┤
│  qcount   = 2          // 当前有 2 个元素                    │
│  dataqsiz = 4          // 缓冲区大小 4                       │
│  closed   = 0          // 未关闭                            │
│  sendx    = 2          // 下次发送到位置 2                   │
│  recvx    = 0          // 下次从位置 0 接收                  │
├─────────────────────────────────────────────────────────────┤
│  buf（环形缓冲区）                                           │
│  ┌─────┬─────┬─────┬─────┐                                 │
│  │  A  │  B  │     │     │                                 │
│  └─────┴─────┴─────┴─────┘                                 │
│    ↑           ↑                                            │
│  recvx       sendx                                          │
├─────────────────────────────────────────────────────────────┤
│  sendq: [G1] -> [G2] -> nil   // 等待发送的 goroutine       │
│  recvq: [G3] -> nil           // 等待接收的 goroutine       │
└─────────────────────────────────────────────────────────────┘
```

### 5.3 发送流程

```
ch <- data

1. 加锁

2. 检查 recvq 有没有等待的接收者
   ├── 有：直接把数据给它，唤醒它
   └── 没有：继续

3. 检查缓冲区有没有空位
   ├── 有：数据放入缓冲区
   └── 没有：继续

4. 缓冲区满了
   ├── 把自己放入 sendq
   └── 阻塞等待

5. 解锁
```

### 5.4 接收流程

```
data := <-ch

1. 加锁

2. 检查 sendq 有没有等待的发送者
   ├── 有：从它那拿数据，唤醒它
   └── 没有：继续

3. 检查缓冲区有没有数据
   ├── 有：从缓冲区取数据
   └── 没有：继续

4. 缓冲区空了
   ├── 把自己放入 recvq
   └── 阻塞等待

5. 解锁
```

---

## 六、Select 多路复用

### 6.1 基本用法

```go
select {
case v := <-ch1:
    fmt.Println("从 ch1 收到:", v)
case ch2 <- 42:
    fmt.Println("发送到 ch2")
case <-time.After(time.Second):
    fmt.Println("超时")
default:
    fmt.Println("没有 channel 准备好")
}
```

### 6.2 特点

```
1. 随机选择一个准备好的 case 执行
2. 如果没有 case 准备好：
   - 有 default：执行 default
   - 没有 default：阻塞等待
3. 空 select{} 永远阻塞
```

### 6.3 常见用法

**超时控制：**
```go
select {
case result := <-ch:
    fmt.Println(result)
case <-time.After(3 * time.Second):
    fmt.Println("超时了")
}
```

**非阻塞操作：**
```go
select {
case v := <-ch:
    fmt.Println(v)
default:
    fmt.Println("channel 没数据")
}
```

---

## 七、sync 包常用工具

### 7.1 WaitGroup（等待一组 goroutine）

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)  // 计数 +1
    go func(n int) {
        defer wg.Done()  // 计数 -1
        fmt.Println(n)
    }(i)
}

wg.Wait()  // 等待计数变为 0
fmt.Println("全部完成")
```

### 7.2 Mutex（互斥锁）

```go
var (
    mu    sync.Mutex
    count int
)

func increment() {
    mu.Lock()         // 加锁
    defer mu.Unlock() // 解锁
    count++
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            increment()
        }()
    }
    wg.Wait()
    fmt.Println(count)  // 1000
}
```

### 7.3 RWMutex（读写锁）

```go
var (
    mu    sync.RWMutex
    data  int
)

func read() int {
    mu.RLock()         // 读锁（多个可以同时读）
    defer mu.RUnlock()
    return data
}

func write(v int) {
    mu.Lock()          // 写锁（独占）
    defer mu.Unlock()
    data = v
}
```

### 7.4 Once（只执行一次）

```go
var once sync.Once

func initialize() {
    fmt.Println("初始化")
}

func main() {
    for i := 0; i < 10; i++ {
        go func() {
            once.Do(initialize)  // 只会执行一次
        }()
    }
    time.Sleep(time.Second)
}
// 只打印一次 "初始化"
```

---

## 八、并发模式实战

### 8.1 生产者-消费者

```go
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
        fmt.Println("生产:", i)
    }
    close(ch)
}

func consumer(ch <-chan int) {
    for v := range ch {
        fmt.Println("消费:", v)
    }
}

func main() {
    ch := make(chan int, 3)
    
    go producer(ch)
    consumer(ch)
}
```

### 8.2 Worker Pool（工作池）

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d 处理任务 %d\n", id, job)
        time.Sleep(time.Second)  // 模拟耗时
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // 启动 3 个 worker
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // 发送 9 个任务
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)
    
    // 收集结果
    for r := 1; r <= 9; r++ {
        fmt.Println("结果:", <-results)
    }
}
```

### 8.3 超时控制

```go
func doWork() <-chan int {
    ch := make(chan int)
    go func() {
        time.Sleep(2 * time.Second)  // 模拟耗时操作
        ch <- 42
    }()
    return ch
}

func main() {
    select {
    case result := <-doWork():
        fmt.Println("结果:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("超时了！")
    }
}
```

### 8.4 优雅退出

```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("收到退出信号")
            return
        default:
            fmt.Println("工作中...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan struct{})
    
    go worker(done)
    
    time.Sleep(2 * time.Second)
    close(done)  // 发送退出信号
    
    time.Sleep(time.Second)
    fmt.Println("程序退出")
}
```

### 8.5 并发安全的计数器（三种方式）

```go
// 方式1：Mutex
type Counter1 struct {
    mu    sync.Mutex
    count int
}

func (c *Counter1) Add() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

// 方式2：atomic
type Counter2 struct {
    count int64
}

func (c *Counter2) Add() {
    atomic.AddInt64(&c.count, 1)
}

// 方式3：channel
type Counter3 struct {
    ch chan struct{}
}

func NewCounter3() *Counter3 {
    c := &Counter3{ch: make(chan struct{})}
    go func() {
        count := 0
        for range c.ch {
            count++
        }
    }()
    return c
}

func (c *Counter3) Add() {
    c.ch <- struct{}{}
}
```

---

## 九、常见错误和陷阱

### 9.1 Goroutine 泄漏

```go
// ❌ 泄漏：channel 没人读，goroutine 永远阻塞
func leak() {
    ch := make(chan int)
    go func() {
        ch <- 42  // 永远阻塞，没人接收
    }()
    // 函数返回，但 goroutine 还在
}

// ✅ 修复：用 select + done
func noLeak() {
    ch := make(chan int)
    done := make(chan struct{})
    
    go func() {
        select {
        case ch <- 42:
        case <-done:
            return
        }
    }()
    
    close(done)  // 通知退出
}
```

### 9.2 死锁

```go
// ❌ 死锁：无缓冲 channel，同一个 goroutine 发送和接收
func main() {
    ch := make(chan int)
    ch <- 42    // 阻塞，等待接收
    v := <-ch   // 永远执行不到
}

// ✅ 修复：用另一个 goroutine
func main() {
    ch := make(chan int)
    go func() {
        ch <- 42
    }()
    v := <-ch
}
```

### 9.3 关闭已关闭的 Channel

```go
ch := make(chan int)
close(ch)
close(ch)  // ❌ panic: close of closed channel

// 只能关闭一次，通常由发送方关闭
```

### 9.4 向已关闭的 Channel 发送

```go
ch := make(chan int)
close(ch)
ch <- 42  // ❌ panic: send on closed channel
```

---

## 十、练习题

### 练习1：并发打印 ⭐

```go
// 启动 3 个 goroutine，分别打印 A、B、C，各打印 5 次
// 用 WaitGroup 等待完成

func main() {
    // 你来实现
}
```

### 练习2：用 channel 求和 ⭐⭐

```go
// 把切片分成两半，用两个 goroutine 分别求和
// 最后把结果加起来

func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    // 你来实现
    // 结果应该是 55
}
```

### 练习3：实现 Worker Pool ⭐⭐

```go
// 创建 3 个 worker，处理 10 个任务
// 每个任务就是打印任务编号

func main() {
    // 你来实现
}
```

### 练习4：实现超时控制 ⭐⭐

```go
// 模拟一个可能很慢的操作
// 如果超过 2 秒还没完成，就放弃

func slowOperation() int {
    time.Sleep(3 * time.Second)
    return 42
}

func main() {
    // 你来实现
    // 应该打印 "超时了"
}
```

### 练习5：并发安全的计数器 ⭐⭐

```go
// 用 Mutex 实现一个并发安全的计数器
// 启动 1000 个 goroutine，每个加 1
// 最后结果应该是 1000

func main() {
    // 你来实现
}
```

---

## 十一、总结

```
┌─────────────────────────────────────────────────────────────┐
│                    Goroutine 要点                           │
├─────────────────────────────────────────────────────────────┤
│  1. go func() 启动 goroutine                                │
│  2. 比线程轻量（2KB vs 1MB）                                 │
│  3. GMP 模型：G=协程，M=线程，P=处理器                       │
│  4. 注意闭包陷阱，循环变量要传参                             │
│  5. main 退出，所有 goroutine 都会被杀死                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     Channel 要点                            │
├─────────────────────────────────────────────────────────────┤
│  1. make(chan T) 无缓冲，make(chan T, n) 有缓冲             │
│  2. 无缓冲是同步的，有缓冲是异步的                           │
│  3. close() 关闭，只能关闭一次                              │
│  4. 向关闭的 channel 发送会 panic                           │
│  5. 从关闭的 channel 接收返回零值                           │
│  6. 用 v, ok := <-ch 判断是否关闭                           │
│  7. 用 range 遍历直到关闭                                   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                      并发工具                               │
├─────────────────────────────────────────────────────────────┤
│  WaitGroup  - 等待一组 goroutine 完成                       │
│  Mutex      - 互斥锁，保护共享资源                          │
│  RWMutex    - 读写锁，读多写少场景                          │
│  Once       - 只执行一次                                    │
│  select     - 多路复用，超时控制                            │
└─────────────────────────────────────────────────────────────┘
```
