# Go GC 教程3 - GC 流程详解

## Go GC 的演进历史

| 版本 | GC 特点 | STW 时间 |
|------|---------|---------|
| Go 1.0 | 标记-清除，完全 STW | 几百毫秒 |
| Go 1.3 | 并发清除 | 几十毫秒 |
| Go 1.5 | 三色标记 + 并发标记 | 几毫秒 |
| Go 1.8 | 混合写屏障 | 亚毫秒 |
| Go 1.12+ | 持续优化 | < 1ms |

**现在的 Go GC 非常快，STW 通常在 1ms 以内！**

---

## GC 的完整流程

```
┌─────────────────────────────────────────────────────────────────┐
│                        GC 完整流程                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐  ┌──────────────────────┐  ┌──────────┐          │
│  │  STW 1   │  │     并发标记          │  │  STW 2   │          │
│  │ Mark     │  │   Concurrent Mark    │  │ Mark     │          │
│  │ Setup    │  │   (和用户程序并发)     │  │Termination│         │
│  └──────────┘  └──────────────────────┘  └──────────┘          │
│      ↓                   ↓                    ↓                 │
│   开启写屏障          扫描所有对象           关闭写屏障            │
│   准备工作            标记存活对象           最终确认              │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    并发清除                               │   │
│  │              Concurrent Sweep                            │   │
│  │              (和用户程序并发)                              │   │
│  └──────────────────────────────────────────────────────────┘   │
│                          ↓                                      │
│                    回收白色对象                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 阶段1：Mark Setup（STW）

**这是第一次 STW，非常短暂（通常 < 0.1ms）**

```
做的事情：
1. 开启写屏障
2. 扫描栈，找到根对象
3. 把根对象引用的对象标记为灰色

┌─────────────────────────────────────┐
│           Mark Setup                │
│                                     │
│  1. 暂停所有 goroutine（STW）        │
│  2. 开启写屏障                       │
│  3. 扫描每个 goroutine 的栈          │
│  4. 标记根对象                       │
│  5. 恢复所有 goroutine              │
│                                     │
└─────────────────────────────────────┘
```

---

## 阶段2：Concurrent Mark（并发）

**这是最耗时的阶段，但和用户程序并发执行，不会卡住程序**

```
做的事情：
1. 从灰色对象开始，遍历所有可达对象
2. 灰色变黑色，引用变灰色
3. 直到没有灰色对象

┌─────────────────────────────────────┐
│        Concurrent Mark              │
│                                     │
│  用户程序          GC 标记           │
│     │                │              │
│     ↓                ↓              │
│  正常执行         扫描灰色对象        │
│  写屏障保护       标记存活对象        │
│     │                │              │
│     └───── 并发执行 ─────┘           │
│                                     │
└─────────────────────────────────────┘
```

**GC 使用 25% 的 CPU 进行标记（可配置）**

---

## 阶段3：Mark Termination（STW）

**这是第二次 STW，也很短暂（通常 < 0.5ms）**

```
做的事情：
1. 暂停所有 goroutine
2. 处理写屏障期间新增的灰色对象
3. 关闭写屏障
4. 计算下次 GC 的触发阈值

┌─────────────────────────────────────┐
│       Mark Termination              │
│                                     │
│  1. 暂停所有 goroutine（STW）        │
│  2. 完成剩余的标记工作               │
│  3. 关闭写屏障                       │
│  4. 恢复所有 goroutine              │
│                                     │
└─────────────────────────────────────┘
```

---

## 阶段4：Concurrent Sweep（并发）

**回收白色对象，和用户程序并发执行**

```
做的事情：
1. 遍历所有内存页
2. 回收白色对象占用的内存
3. 把内存归还给内存池

┌─────────────────────────────────────┐
│        Concurrent Sweep             │
│                                     │
│  用户程序          GC 清除           │
│     │                │              │
│     ↓                ↓              │
│  正常执行         回收白色对象        │
│  分配内存         归还内存池         │
│     │                │              │
│     └───── 并发执行 ─────┘           │
│                                     │
└─────────────────────────────────────┘
```

---

## 动手实验：观察 GC 各阶段

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

func main() {
    // 关闭自动 GC，手动控制
    debug.SetGCPercent(-1)
    
    // 分配一些内存
    var data [][]byte
    for i := 0; i < 100; i++ {
        data = append(data, make([]byte, 1024*1024))  // 100MB
    }
    
    fmt.Println("分配完成，准备 GC")
    
    // 让一部分变成垃圾
    data = data[:50]  // 只保留前 50MB
    
    // 记录 GC 前的状态
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    fmt.Printf("GC 前: Alloc=%dMB\n", m1.Alloc/1024/1024)
    
    // 手动触发 GC 并计时
    start := time.Now()
    runtime.GC()
    duration := time.Since(start)
    
    // 记录 GC 后的状态
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("GC 后: Alloc=%dMB\n", m2.Alloc/1024/1024)
    fmt.Printf("GC 耗时: %v\n", duration)
    fmt.Printf("GC 次数: %d\n", m2.NumGC)
    fmt.Printf("GC 总暂停时间: %v\n", time.Duration(m2.PauseTotalNs))
    fmt.Printf("上次 GC 暂停时间: %v\n", time.Duration(m2.PauseNs[(m2.NumGC+255)%256]))
}
```

**运行（带详细 GC 日志）：**
```bash
GODEBUG=gctrace=1 go run main.go
```

---

## GC 触发时机

### 1. 堆内存达到阈值

```
当前堆大小 >= 上次 GC 后堆大小 * (1 + GOGC/100)

默认 GOGC=100，意味着堆翻倍时触发 GC

例如：
- 上次 GC 后堆大小：100MB
- GOGC=100
- 触发阈值：100MB * 2 = 200MB
- 当堆达到 200MB 时触发 GC
```

### 2. 定时触发

```
如果 2 分钟没有 GC，强制触发一次
（防止内存一直不释放）
```

### 3. 手动触发

```go
runtime.GC()  // 手动触发
```

---

## GOGC 参数

**GOGC 控制 GC 的触发频率**

```bash
# 默认值 100，堆翻倍时触发
GOGC=100 ./myapp

# 设为 200，堆增长到 3 倍时触发（GC 更少，内存占用更高）
GOGC=200 ./myapp

# 设为 50，堆增长到 1.5 倍时触发（GC 更频繁，内存占用更低）
GOGC=50 ./myapp

# 关闭 GC（危险！内存会一直增长）
GOGC=off ./myapp
```

**代码中设置：**
```go
import "runtime/debug"

// 设置 GOGC
debug.SetGCPercent(200)

// 关闭 GC
debug.SetGCPercent(-1)
```

---

## 动手实验：GOGC 的影响

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

func allocateMemory() {
    var data [][]byte
    for i := 0; i < 1000; i++ {
        data = append(data, make([]byte, 1024*1024))  // 1MB
        if i%100 == 0 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            fmt.Printf("i=%d, Alloc=%dMB, NumGC=%d\n", 
                i, m.Alloc/1024/1024, m.NumGC)
        }
        time.Sleep(time.Millisecond)
    }
}

func main() {
    // 试试不同的 GOGC 值
    // debug.SetGCPercent(50)   // GC 更频繁
    // debug.SetGCPercent(100)  // 默认
    // debug.SetGCPercent(200)  // GC 更少
    
    allocateMemory()
}
```

---

## GC Pacer（调度器）

**Go 的 GC 有一个智能调度器，自动调整 GC 节奏**

```
目标：
1. GC 占用 CPU 不超过 25%
2. 堆大小不超过目标值
3. STW 时间尽可能短

调整策略：
- 如果堆增长太快 → 加快标记速度
- 如果 CPU 占用太高 → 减慢标记速度
- 动态平衡
```

---

## GC 日志详解

```bash
GODEBUG=gctrace=1 go run main.go
```

```
gc 1 @0.012s 2%: 0.018+1.2+0.003 ms clock, 0.14+0/1.1/0+0.024 ms cpu, 4->4->0 MB, 5 MB goal, 8 P
```

**拆解：**

| 字段 | 含义 |
|------|------|
| `gc 1` | 第 1 次 GC |
| `@0.012s` | 程序启动后 0.012 秒 |
| `2%` | GC 占用 CPU 2% |
| `0.018+1.2+0.003 ms clock` | STW1 + 并发标记 + STW2 的墙钟时间 |
| `0.14+0/1.1/0+0.024 ms cpu` | CPU 时间 |
| `4->4->0 MB` | GC 前堆 → 标记时堆 → GC 后堆 |
| `5 MB goal` | 目标堆大小 |
| `8 P` | 使用 8 个 P |

**重点关注：**
- STW 时间（第一个和第三个数字）
- 堆大小变化
- GC 频率

---

## 练习

### 练习1：分析 GC 日志

```
gc 5 @1.234s 3%: 0.021+2.5+0.004 ms clock, 0.16+0/2.3/0.1+0.032 ms cpu, 8->9->4 MB, 10 MB goal, 8 P
```

问题：
1. 这是第几次 GC？
2. STW 总时间是多少？
3. GC 后堆大小是多少？
4. 回收了多少内存？

<details>
<summary>点击查看答案</summary>

1. 第 5 次 GC
2. STW 总时间 = 0.021 + 0.004 = 0.025ms
3. GC 后堆大小 = 4MB
4. 回收了 8 - 4 = 4MB（从 GC 前 8MB 到 GC 后 4MB）

</details>

### 练习2：调整 GOGC

```go
// 场景：你的程序内存占用太高，想让 GC 更频繁
// 应该怎么设置 GOGC？

// A. debug.SetGCPercent(200)
// B. debug.SetGCPercent(100)
// C. debug.SetGCPercent(50)
// D. debug.SetGCPercent(-1)
```

<details>
<summary>点击查看答案</summary>

**C. debug.SetGCPercent(50)**

GOGC=50 意味着堆增长 50% 就触发 GC，比默认的 100% 更频繁，内存占用更低。

- A (200)：GC 更少，内存更高
- B (100)：默认值
- D (-1)：关闭 GC，内存会一直增长

</details>

---

## 本章总结

| 阶段 | 是否 STW | 做什么 |
|------|---------|--------|
| Mark Setup | ✓ (< 0.1ms) | 开启写屏障，扫描根对象 |
| Concurrent Mark | ✗ | 并发标记所有可达对象 |
| Mark Termination | ✓ (< 0.5ms) | 关闭写屏障，完成标记 |
| Concurrent Sweep | ✗ | 并发回收白色对象 |

| 参数 | 说明 |
|------|------|
| GOGC | 控制 GC 触发频率，默认 100 |
| GODEBUG=gctrace=1 | 打印 GC 日志 |
| runtime.GC() | 手动触发 GC |
| debug.SetGCPercent() | 代码中设置 GOGC |

---

## 下一篇

[GC教程4-性能优化](./GC教程4-性能优化.md) - 如何减少 GC 压力，优化程序性能
