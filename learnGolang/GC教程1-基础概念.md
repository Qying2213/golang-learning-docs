# Go GC 教程1 - 基础概念

## 什么是 GC？

**GC = Garbage Collection = 垃圾回收**

程序运行时会不断申请内存，用完后需要释放。GC 就是**自动帮你释放不再使用的内存**。

```
没有 GC（C 语言）：
malloc() 申请内存
free() 手动释放    ← 忘了释放 = 内存泄漏，释放错了 = 程序崩溃

有 GC（Go 语言）：
make/new 申请内存
自动释放           ← 不用管，GC 帮你搞定
```

---

## 为什么需要 GC？
·
**没有 GC 的痛苦（C/C++）：**

```c
// C 语言
char* str = malloc(100);  // 申请内存
// ... 使用 str
free(str);                // 必须手动释放

// 问题1：忘记 free → 内存泄漏
// 问题2：free 两次 → 程序崩溃
// 问题3：free 后继续使用 → 野指针
```

**有 GC 的幸福（Go）：**

```go
// Go 语言
func foo() {
    str := make([]byte, 100)  // 申请内存
    // ... 使用 str
}   // 函数结束，str 不再使用，GC 自动回收

// 不用管内存释放，专心写业务逻辑
```

---

## GC 的核心问题

GC 要解决两个问题：

### 问题1：哪些内存可以回收？

**判断标准：这块内存还有没有人在用？**

```go
func main() {
    a := new(int)    // a 指向一块内存
    b := a           // b 也指向这块内存
    a = nil          // a 不再指向了
    // 这块内存能回收吗？不能！因为 b 还在用
    
    b = nil          // b 也不指向了
    // 现在可以回收了，没人用了
}
```

### 问题2：什么时候回收？

- 回收太频繁 → 程序卡顿（STW）
- 回收太少 → 内存占用高

Go 的 GC 会自动平衡这两点。

---

## 什么是"垃圾"？

**垃圾 = 不可达的对象 = 没有任何引用指向的内存**

```go
func createGarbage() {
    // 情况1：局部变量，函数结束后变成垃圾
    data := make([]int, 1000)
    _ = data
}   // data 变成垃圾

func main() {
    createGarbage()
    // 此时 data 那块内存没人引用了，是垃圾
    
    // 情况2：覆盖引用
    a := new(int)
    *a = 100
    a = new(int)  // 原来的内存没人引用了，是垃圾
    
    // 情况3：置为 nil
    b := make([]int, 100)
    b = nil       // 原来的切片内存是垃圾
}
```

---

## 可达性分析

**Go 用"可达性分析"判断对象是否是垃圾**

```
从"根对象"出发，能访问到的对象都是"活的"，访问不到的就是"垃圾"

根对象（GC Roots）：
├── 全局变量
├── 栈上的局部变量
├── 寄存器中的值
└── ...

        根对象
           │
           ↓
        对象A ──→ 对象B ──→ 对象C
           │
           ↓
        对象D
        
        对象E（没有引用指向它）← 垃圾！
```

**例子：**

```go
var global *int  // 全局变量，是根对象

func main() {
    a := new(int)   // a 是栈上变量，是根对象
    b := new(int)   // b 是栈上变量，是根对象
    
    global = a      // global 指向 a 指向的内存
    
    c := new(int)   // c 是栈上变量
    c = nil         // c 置空，原来的内存变成垃圾
    
    // 可达的对象：a 指向的内存（通过 a 和 global 可达）
    //            b 指向的内存（通过 b 可达）
    // 垃圾：原来 c 指向的内存
}
```

---

## 动手实验1：观察 GC

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    // 打印当前内存状态
    printMemStats("初始状态")
    
    // 分配一些内存
    data := make([]byte, 100*1024*1024)  // 100MB
    _ = data
    printMemStats("分配 100MB 后")
    
    // 让 data 变成垃圾
    data = nil
    printMemStats("data=nil 后（GC 前）")
    
    // 手动触发 GC
    runtime.GC()
    printMemStats("手动 GC 后")
}

func printMemStats(msg string) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("%s:\n", msg)
    fmt.Printf("  Alloc = %d MB\n", m.Alloc/1024/1024)
    fmt.Printf("  TotalAlloc = %d MB\n", m.TotalAlloc/1024/1024)
    fmt.Printf("  NumGC = %d\n\n", m.NumGC)
}
```

**运行：**
```bash
go run main.go
```

**预期输出：**
```
初始状态:
  Alloc = 0 MB
  TotalAlloc = 0 MB
  NumGC = 0

分配 100MB 后:
  Alloc = 100 MB
  TotalAlloc = 100 MB
  NumGC = 0

data=nil 后（GC 前）:
  Alloc = 100 MB      ← 内存还没释放
  TotalAlloc = 100 MB
  NumGC = 0

手动 GC 后:
  Alloc = 0 MB        ← 内存被回收了！
  TotalAlloc = 100 MB
  NumGC = 1           ← GC 执行了 1 次
```

---

## 动手实验2：观察自动 GC

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 开启 GC 日志
    // 运行时加环境变量：GODEBUG=gctrace=1 go run main.go
    
    for i := 0; i < 10; i++ {
        // 每次循环分配 10MB，然后丢弃
        data := make([]byte, 10*1024*1024)
        _ = data
        
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        fmt.Printf("第 %d 次: Alloc=%dMB, NumGC=%d\n", 
            i+1, m.Alloc/1024/1024, m.NumGC)
        
        time.Sleep(100 * time.Millisecond)
    }
}
```

**运行（带 GC 日志）：**
```bash
GODEBUG=gctrace=1 go run main.go
```

**你会看到类似这样的 GC 日志：**
```
gc 1 @0.012s 2%: 0.018+1.2+0.003 ms clock, 0.14+0/1.1/0+0.024 ms cpu, 4->4->0 MB, 5 MB goal, 8 P
```

---

## GC 日志解读

```
gc 1 @0.012s 2%: 0.018+1.2+0.003 ms clock, 0.14+0/1.1/0+0.024 ms cpu, 4->4->0 MB, 5 MB goal, 8 P
│   │       │   │                        │                          │           │         │
│   │       │   │                        │                          │           │         └─ P 数量
│   │       │   │                        │                          │           └─ 目标堆大小
│   │       │   │                        │                          └─ 堆大小变化：GC前->GC中->GC后
│   │       │   │                        └─ CPU 时间
│   │       │   └─ 墙钟时间（STW1 + 并发标记 + STW2）
│   │       └─ GC 占用 CPU 百分比
│   └─ 程序启动后的时间
└─ 第几次 GC
```

**重点关注：**
- `4->4->0 MB`：GC 前 4MB，GC 后 0MB，回收了 4MB
- `2%`：GC 只占用了 2% 的 CPU，影响很小

---

## 练习

### 练习1：预测 GC 行为

```go
func main() {
    a := make([]int, 1000)
    b := a
    a = nil
    runtime.GC()
    // 问：a 原来指向的内存会被回收吗？
}
```

<details>
<summary>点击查看答案</summary>

**不会被回收**

因为 `b` 还指向那块内存，所以它是可达的，不是垃圾。

</details>

### 练习2：制造垃圾

```go
// 写一个函数，每次调用都会产生垃圾
func createGarbage() {
    // TODO: 你来实现
}

func main() {
    for i := 0; i < 100; i++ {
        createGarbage()
    }
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("GC 次数: %d\n", m.NumGC)
}
```

<details>
<summary>点击查看答案</summary>

```go
func createGarbage() {
    // 分配内存，函数结束后变成垃圾
    data := make([]byte, 1024*1024)  // 1MB
    _ = data
}
```

</details>

---

## 本章总结

| 概念 | 说明 |
|------|------|
| GC | 自动回收不再使用的内存 |
| 垃圾 | 没有任何引用指向的对象 |
| 可达性分析 | 从根对象出发，能访问到的是活的 |
| 根对象 | 全局变量、栈上变量、寄存器 |
| runtime.GC() | 手动触发 GC |
| runtime.MemStats | 查看内存统计 |
| GODEBUG=gctrace=1 | 打印 GC 日志 |

---

## 下一篇

[GC教程2-三色标记法](./GC教程2-三色标记法.md) - Go GC 的核心算法
