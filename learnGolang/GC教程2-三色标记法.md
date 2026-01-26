# Go GC 教程2 - 三色标记法

## 回顾：GC 要做什么？

1. **找出垃圾**：哪些对象没人用了
2. **回收内存**：释放垃圾占用的内存

Go 用**三色标记法**来找出垃圾。

---

## 什么是三色标记法？

**把所有对象分成三种颜色：**

| 颜色 | 含义 | 状态 |
|------|------|------|
| ⚪ 白色 | 可能是垃圾 | 还没扫描到 |
| ⚫ 灰色 | 正在处理 | 扫描到了，但还没扫描它的引用 |
| ⬛ 黑色 | 确定存活 | 扫描完了，它和它的引用都处理过了 |

**最终：白色的就是垃圾，黑色的就是存活对象**

---

## 三色标记的过程

### 初始状态

```
所有对象都是白色

根对象
   │
   ↓
  A(⚪) ──→ B(⚪) ──→ C(⚪)
   │
   ↓
  D(⚪)
  
  E(⚪)  ← 没有引用指向它
```

### 第1步：从根对象开始

```
把根对象能直接访问的对象标记为灰色

根对象
   │
   ↓
  A(⚫) ──→ B(⚪) ──→ C(⚪)    A 变成灰色
   │
   ↓
  D(⚪)
  
  E(⚪)
```

### 第2步：处理灰色对象

```
取出一个灰色对象 A：
1. 把 A 标记为黑色
2. 把 A 引用的对象（B、D）标记为灰色

根对象
   │
   ↓
  A(⬛) ──→ B(⚫) ──→ C(⚪)    A 变黑，B、D 变灰
   │
   ↓
  D(⚫)
  
  E(⚪)
```

### 第3步：继续处理灰色对象

```
取出灰色对象 B：
1. 把 B 标记为黑色
2. 把 B 引用的对象（C）标记为灰色

根对象
   │
   ↓
  A(⬛) ──→ B(⬛) ──→ C(⚫)    B 变黑，C 变灰
   │
   ↓
  D(⚫)
  
  E(⚪)
```

### 第4步：继续处理

```
取出灰色对象 D：D 没有引用其他对象，直接变黑
取出灰色对象 C：C 没有引用其他对象，直接变黑

根对象
   │
   ↓
  A(⬛) ──→ B(⬛) ──→ C(⬛)    全部变黑
   │
   ↓
  D(⬛)
  
  E(⚪)  ← 还是白色！
```

### 第5步：回收白色对象

```
没有灰色对象了，标记结束
白色对象 E 就是垃圾，回收它！
```

---

## 用代码模拟三色标记

```go
package main

import "fmt"

// 对象
type Object struct {
    name  string
    color string      // "white", "gray", "black"
    refs  []*Object   // 引用的其他对象
}

// 三色标记
func triColorMark(roots []*Object) []*Object {
    // 1. 初始化：所有对象都是白色（假设已经是白色）
    
    // 2. 把根对象引用的对象标记为灰色
    grayList := make([]*Object, 0)
    for _, root := range roots {
        for _, ref := range root.refs {
            if ref.color == "white" {
                ref.color = "gray"
                grayList = append(grayList, ref)
                fmt.Printf("  %s: 白色 → 灰色\n", ref.name)
            }
        }
    }
    
    // 3. 处理灰色对象
    for len(grayList) > 0 {
        // 取出一个灰色对象
        obj := grayList[0]
        grayList = grayList[1:]
        
        // 标记为黑色
        obj.color = "black"
        fmt.Printf("  %s: 灰色 → 黑色\n", obj.name)
        
        // 把它引用的白色对象标记为灰色
        for _, ref := range obj.refs {
            if ref.color == "white" {
                ref.color = "gray"
                grayList = append(grayList, ref)
                fmt.Printf("  %s: 白色 → 灰色\n", ref.name)
            }
        }
    }
    
    // 4. 返回所有白色对象（垃圾）
    return nil  // 简化，实际要遍历所有对象找白色的
}

func main() {
    // 创建对象
    a := &Object{name: "A", color: "white"}
    b := &Object{name: "B", color: "white"}
    c := &Object{name: "C", color: "white"}
    d := &Object{name: "D", color: "white"}
    e := &Object{name: "E", color: "white"}  // 垃圾
    
    // 建立引用关系
    a.refs = []*Object{b, d}
    b.refs = []*Object{c}
    // c, d, e 没有引用其他对象
    
    // 根对象引用 a
    root := &Object{name: "root", refs: []*Object{a}}
    
    fmt.Println("开始三色标记：")
    triColorMark([]*Object{root})
    
    fmt.Println("\n最终状态：")
    for _, obj := range []*Object{a, b, c, d, e} {
        status := "存活"
        if obj.color == "white" {
            status = "垃圾"
        }
        fmt.Printf("  %s: %s (%s)\n", obj.name, obj.color, status)
    }
}
```

**运行结果：**
```
开始三色标记：
  A: 白色 → 灰色
  A: 灰色 → 黑色
  B: 白色 → 灰色
  D: 白色 → 灰色
  B: 灰色 → 黑色
  C: 白色 → 灰色
  D: 灰色 → 黑色
  C: 灰色 → 黑色

最终状态：
  A: black (存活)
  B: black (存活)
  C: black (存活)
  D: black (存活)
  E: white (垃圾)    ← E 是垃圾！
```

---

## 为什么用三色而不是两色？

**两色（黑白）的问题：**

如果只有黑白两色，GC 必须一次性完成，不能中断。因为中断后无法知道哪些对象处理过了。

**三色的好处：**

灰色表示"正在处理"，GC 可以随时暂停，下次从灰色对象继续。这就是**增量 GC** 的基础。

```
两色：必须一口气扫完，程序卡住
三色：可以分多次扫，每次扫一点，程序不会卡太久
```

---

## 并发标记的问题

Go 的 GC 是**并发**的，GC 和用户程序同时运行。这会带来问题：

### 问题：漏标

```
GC 正在标记，用户程序同时在修改引用

初始：
  A(⬛) ──→ B(⚫)
            │
            ↓
           C(⚪)

GC 已经扫描完 A（黑色），正准备扫描 B

此时用户程序执行：
  A.ref = C    // A 直接引用 C
  B.ref = nil  // B 不再引用 C

变成：
  A(⬛) ──→ C(⚪)    A 引用 C
  B(⚫)              B 不再引用 C

问题：
- A 是黑色，不会再扫描
- B 不再引用 C
- C 还是白色，会被当成垃圾回收！
- 但 C 实际上被 A 引用着，不应该回收！
```

**这就是"漏标"问题，会导致程序崩溃！**

---

## 写屏障（Write Barrier）

**写屏障 = 在修改引用时插入额外代码，防止漏标**

Go 使用**混合写屏障**：

```go
// 伪代码：写屏障
func writePointer(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    // 1. 把旧值标记为灰色（删除写屏障）
    shade(*slot)
    
    // 2. 把新值标记为灰色（插入写屏障）
    shade(ptr)
    
    // 3. 执行实际的赋值
    *slot = ptr
}
```

**有了写屏障，上面的问题就解决了：**

```
用户程序执行 A.ref = C 时：
1. 写屏障把 C 标记为灰色
2. C 会被扫描，不会被漏掉
```

---

## 动手实验：观察写屏障

```go
package main

import (
    "fmt"
    "runtime"
)

type Node struct {
    data int
    next *Node
}

func main() {
    // 开启 GC 追踪
    // GODEBUG=gctrace=1 go run main.go
    
    var head *Node
    
    // 不断创建和修改引用
    for i := 0; i < 1000000; i++ {
        node := &Node{data: i}
        node.next = head  // 修改引用，触发写屏障
        head = node
        
        if i%100000 == 0 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            fmt.Printf("i=%d, Alloc=%dMB, NumGC=%d\n", 
                i, m.Alloc/1024/1024, m.NumGC)
        }
    }
}
```

---

## 三色不变式

**为了保证正确性，GC 必须维护"三色不变式"：**

**强三色不变式：** 黑色对象不能直接引用白色对象

**弱三色不变式：** 黑色对象可以引用白色对象，但白色对象必须有灰色对象保护（能通过灰色对象访问到）

Go 使用**弱三色不变式 + 写屏障**来保证正确性。

```
允许：
  A(⬛) ──→ C(⚪)
            ↑
  B(⚫) ────┘     C 有灰色对象 B 保护，不会被漏标

不允许（写屏障会阻止）：
  A(⬛) ──→ C(⚪)    C 没有灰色保护，会被漏标
```

---

## 练习

### 练习1：手动模拟三色标记

```
给定对象和引用关系：
- 根对象引用 A
- A 引用 B 和 C
- B 引用 D
- E 没有被引用

请写出三色标记的每一步状态变化。
```

<details>
<summary>点击查看答案</summary>

```
初始：A(⚪) B(⚪) C(⚪) D(⚪) E(⚪)

第1步：从根开始，A 变灰
       A(⚫) B(⚪) C(⚪) D(⚪) E(⚪)

第2步：处理 A，A 变黑，B、C 变灰
       A(⬛) B(⚫) C(⚫) D(⚪) E(⚪)

第3步：处理 B，B 变黑，D 变灰
       A(⬛) B(⬛) C(⚫) D(⚫) E(⚪)

第4步：处理 C，C 变黑（C 没有引用其他对象）
       A(⬛) B(⬛) C(⬛) D(⚫) E(⚪)

第5步：处理 D，D 变黑
       A(⬛) B(⬛) C(⬛) D(⬛) E(⚪)

结束：E 是白色，是垃圾
```

</details>

### 练习2：判断是否会漏标

```
GC 正在进行，当前状态：
  A(⬛) ──→ B(⚫)
            │
            ↓
           C(⚪)

用户程序执行：B.next = nil（B 不再引用 C）

问：C 会被漏标吗？
```

<details>
<summary>点击查看答案</summary>

**不会漏标**

因为 B 是灰色，还没处理完。GC 会继续处理 B，发现 B 引用 C（在删除之前的快照），C 会被标记为灰色。

写屏障也会在 `B.next = nil` 时把旧值 C 标记为灰色。

</details>

---

## 本章总结

| 概念 | 说明 |
|------|------|
| 三色标记 | 白色（可能垃圾）、灰色（处理中）、黑色（存活） |
| 标记过程 | 从根开始，灰色变黑色，引用变灰色，直到没有灰色 |
| 写屏障 | 修改引用时插入代码，防止漏标 |
| 三色不变式 | 保证黑色对象引用的白色对象不会被漏标 |

---

## 下一篇

[GC教程3-GC流程详解](./GC教程3-GC流程详解.md) - Go GC 的完整执行流程
