# Goroutine 调度模型（GMP）

## 1. GMP 模型概述

```
┌─────────────────────────────────────────────────────────────┐
│                        Go Runtime                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                 │
│   │  G  │ │  G  │ │  G  │ │  G  │ │  G  │  ... Goroutines │
│   └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘                 │
│      │       │       │       │       │                      │
│   ┌──▼───────▼───────▼───────▼───────▼──┐                  │
│   │         Global Run Queue            │                  │
│   └─────────────────────────────────────┘                  │
│                      │                                      │
│      ┌───────────────┼───────────────┐                     │
│      │               │               │                      │
│   ┌──▼──┐         ┌──▼──┐         ┌──▼──┐                  │
│   │  P  │         │  P  │         │  P  │  Processors      │
│   │Local│         │Local│         │Local│  (GOMAXPROCS)    │
│   │Queue│         │Queue│         │Queue│                  │
│   └──┬──┘         └──┬──┘         └──┬──┘                  │
│      │               │               │                      │
│   ┌──▼──┐         ┌──▼──┐         ┌──▼──┐                  │
│   │  M  │         │  M  │         │  M  │  OS Threads      │
│   └──┬──┘         └──┬──┘         └──┬──┘                  │
│      │               │               │                      │
└──────┼───────────────┼───────────────┼──────────────────────┘
       │               │               │
    ┌──▼──┐         ┌──▼──┐         ┌──▼──┐
    │ CPU │         │ CPU │         │ CPU │  Hardware
    └─────┘         └─────┘         └─────┘
```

## 2. G、M、P 详解

### G (Goroutine)
```go
// runtime/runtime2.go
type g struct {
    stack       stack   // 栈内存范围 [stack.lo, stack.hi)
    stackguard0 uintptr // 栈溢出检查
    
    m         *m       // 当前关联的 M
    sched     gobuf    // 调度信息（SP, PC, BP 等）
    atomicstatus uint32 // 状态
    goid      int64    // goroutine ID
    
    // ... 更多字段
}

// Goroutine 状态
const (
    _Gidle     = iota // 刚创建，未初始化
    _Grunnable        // 可运行，在队列中等待
    _Grunning         // 正在运行
    _Gsyscall         // 系统调用中
    _Gwaiting         // 等待中（channel、锁等）
    _Gdead            // 已结束
)
```

### M (Machine/OS Thread)
```go
type m struct {
    g0      *g     // 调度栈，用于执行调度代码
    curg    *g     // 当前运行的 G
    p       *p     // 关联的 P
    
    spinning bool  // 是否在自旋寻找工作
    // ...
}
```

### P (Processor)
```go
type p struct {
    id          int32
    status      uint32 // Pidle, Prunning, Psyscall, Pgcstop, Pdead
    
    m           *m     // 关联的 M
    runqhead    uint32 // 本地队列头
    runqtail    uint32 // 本地队列尾
    runq        [256]guintptr // 本地运行队列（固定大小 256）
    
    runnext     guintptr // 下一个要运行的 G（优先级最高）
    // ...
}
```

## 3. 调度流程

### 创建 Goroutine
```go
go func() {
    // ...
}()

// 1. 创建 G 结构体
// 2. 优先放入当前 P 的 runnext
// 3. 如果 runnext 已有，放入本地队列
// 4. 如果本地队列满，放一半到全局队列
```

### 调度循环
```go
// runtime/proc.go
func schedule() {
    // 1. 每 61 次调度，从全局队列获取 G（防止饥饿）
    if gp == nil {
        if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
            gp = globrunqget(_g_.m.p.ptr(), 1)
        }
    }
    
    // 2. 从本地队列获取
    if gp == nil {
        gp, inheritTime = runqget(_g_.m.p.ptr())
    }
    
    // 3. 从其他地方获取（全局队列、网络轮询、其他 P）
    if gp == nil {
        gp, inheritTime = findrunnable()
    }
    
    // 4. 执行 G
    execute(gp, inheritTime)
}
```

### Work Stealing（工作窃取）
```go
// 当 P 的本地队列为空时
func findrunnable() (gp *g, inheritTime bool) {
    // 1. 检查本地队列
    // 2. 检查全局队列
    // 3. 检查网络轮询器
    // 4. 从其他 P 窃取一半的 G
    
    for i := 0; i < 4; i++ {
        // 随机选择一个 P
        // 窃取其一半的 G
    }
}
```

## 4. 抢占式调度

### 基于协作的抢占（Go 1.13 之前）
```go
// 在函数调用时检查是否需要抢占
func someFunction() {
    // 编译器插入的检查代码
    if stackguard0 == stackPreempt {
        // 让出 CPU
    }
    // 函数体...
}

// 问题：死循环无法被抢占
for {
    // 没有函数调用，永远不会被抢占
}
```

### 基于信号的抢占（Go 1.14+）
```go
// 使用 SIGURG 信号实现异步抢占
// 即使是死循环也能被抢占

// sysmon 监控线程定期检查
func sysmon() {
    for {
        // 检查运行时间过长的 G
        // 发送 SIGURG 信号
        // G 收到信号后保存状态，让出 CPU
    }
}
```

## 5. 系统调用处理

```go
// 同步系统调用（如文件 I/O）
func syscall() {
    // 1. M 进入系统调用
    // 2. P 与 M 解绑
    // 3. P 寻找其他 M 继续执行
    // 4. 系统调用返回后，M 尝试获取 P
    // 5. 如果没有空闲 P，M 进入休眠
}

// 异步系统调用（如网络 I/O）
// 使用 netpoller，不会阻塞 M
func netpoll() {
    // 使用 epoll/kqueue 等
    // G 等待时不占用 M
}
```

## 6. GOMAXPROCS

```go
import "runtime"

// 获取当前 P 的数量
n := runtime.GOMAXPROCS(0)

// 设置 P 的数量
runtime.GOMAXPROCS(4)

// 默认值 = CPU 核心数
// 通常不需要修改
```

## 7. 调度相关的运行时函数

```go
// 让出当前 G 的执行权
runtime.Gosched()

// 获取当前 goroutine 数量
runtime.NumGoroutine()

// 获取当前 G 的 ID（不推荐使用）
// Go 故意不暴露 goroutine ID

// 锁定当前 G 到当前 M
runtime.LockOSThread()
runtime.UnlockOSThread()
```

## 8. 调度器的优化

```go
// 1. 本地队列减少锁竞争
// 2. Work Stealing 负载均衡
// 3. runnext 优化局部性
// 4. 自旋 M 减少线程切换
// 5. 网络轮询器避免阻塞

// 查看调度器状态
GODEBUG=schedtrace=1000 ./myprogram
// 每 1000ms 打印一次调度器状态
```

## 9. 常见面试题

**Q1: Goroutine 和线程的区别？**
```
1. 内存占用：Goroutine 初始栈 2KB，线程通常 1-8MB
2. 创建开销：Goroutine 很小，线程需要系统调用
3. 切换开销：Goroutine 用户态切换，线程需要内核态切换
4. 调度：Goroutine 由 Go runtime 调度，线程由 OS 调度
```

**Q2: 什么情况下 Goroutine 会阻塞？**
```go
// 1. channel 操作
ch <- data  // 发送到满的 channel
<-ch        // 从空的 channel 接收

// 2. 锁操作
mu.Lock()   // 获取已被持有的锁

// 3. 系统调用
file.Read() // 同步 I/O

// 4. time.Sleep
time.Sleep(time.Second)

// 5. select 等待
select {
case <-ch1:
case <-ch2:
}
```

**Q3: 如何控制 Goroutine 的数量？**
```go
// 方法1: 使用带缓冲的 channel 作为信号量
sem := make(chan struct{}, 10) // 最多 10 个并发

for _, task := range tasks {
    sem <- struct{}{} // 获取信号量
    go func(t Task) {
        defer func() { <-sem }() // 释放信号量
        process(t)
    }(task)
}

// 方法2: 使用 worker pool
func workerPool(numWorkers int, tasks <-chan Task) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range tasks {
                process(task)
            }
        }()
    }
    wg.Wait()
}
```

## 10. 调试技巧

```bash
# 查看调度器跟踪
GODEBUG=schedtrace=1000 ./myprogram

# 输出示例
SCHED 0ms: gomaxprocs=4 idleprocs=2 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 [0 0 0 0]

# 更详细的跟踪
GODEBUG=schedtrace=1000,scheddetail=1 ./myprogram

# 使用 trace 工具
go test -trace=trace.out
go tool trace trace.out
```

## 练习

1. 写一个程序，观察 GOMAXPROCS 对性能的影响
2. 实现一个简单的 goroutine 池
3. 使用 `runtime.Gosched()` 观察调度行为
4. 分析一个死锁场景，理解调度器的行为
