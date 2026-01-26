# Go GC 教程4 - 性能优化

## GC 的代价

虽然 Go 的 GC 很快，但还是有代价：

1. **CPU 消耗**：GC 占用约 25% CPU
2. **STW 暂停**：虽然很短，但高频场景可能有影响
3. **内存开销**：需要额外内存存储标记信息

**优化目标：减少 GC 次数和 GC 工作量**

---

## 优化原则

```
GC 工作量 = 对象数量 × 对象存活时间

减少 GC 压力的方法：
1. 减少对象分配数量
2. 减少对象大小
3. 复用对象
4. 让对象分配在栈上（不需要 GC）
```

---

## 优化技巧1：对象复用（sync.Pool）

### 问题：频繁创建临时对象

```go
// 不好：每次请求都创建新的 buffer
func handleRequest() {
    buf := make([]byte, 1024)  // 每次都分配
    // 使用 buf
}   // buf 变成垃圾
```

### 解决：使用 sync.Pool 复用对象

```go
// 好：复用 buffer
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func handleRequest() {
    buf := bufPool.Get().([]byte)  // 从池中获取
    defer bufPool.Put(buf)          // 用完放回池中
    
    // 使用 buf
}
```

### 完整示例

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

// 不使用 Pool
func withoutPool() {
    for i := 0; i < 1000000; i++ {
        buf := make([]byte, 1024)
        _ = buf
    }
}

// 使用 Pool
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func withPool() {
    for i := 0; i < 1000000; i++ {
        buf := pool.Get().([]byte)
        pool.Put(buf)
    }
}

func main() {
    // 测试不使用 Pool
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    gc1 := m1.NumGC
    
    withoutPool()
    
    runtime.ReadMemStats(&m1)
    fmt.Printf("不使用 Pool: GC 次数 = %d\n", m1.NumGC-gc1)
    
    // 测试使用 Pool
    runtime.GC()
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    gc2 := m2.NumGC
    
    withPool()
    
    runtime.ReadMemStats(&m2)
    fmt.Printf("使用 Pool: GC 次数 = %d\n", m2.NumGC-gc2)
}
```

**运行结果：**
```
不使用 Pool: GC 次数 = 15
使用 Pool: GC 次数 = 0
```

---

## 优化技巧2：预分配切片

### 问题：切片动态扩容

```go
// 不好：多次扩容，产生多个废弃的底层数组
func bad() []int {
    var result []int
    for i := 0; i < 10000; i++ {
        result = append(result, i)  // 多次扩容
    }
    return result
}
```

### 解决：预分配容量

```go
// 好：一次分配，不扩容
func good() []int {
    result := make([]int, 0, 10000)  // 预分配容量
    for i := 0; i < 10000; i++ {
        result = append(result, i)  // 不会扩容
    }
    return result
}
```

### 性能对比

```go
package main

import (
    "fmt"
    "runtime"
)

func withoutPrealloc() {
    for i := 0; i < 10000; i++ {
        var s []int
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

func withPrealloc() {
    for i := 0; i < 10000; i++ {
        s := make([]int, 0, 1000)
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

func main() {
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    withoutPrealloc()
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("不预分配: 分配次数=%d, GC次数=%d\n", 
        m2.Mallocs-m1.Mallocs, m2.NumGC-m1.NumGC)
    
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    withPrealloc()
    
    runtime.ReadMemStats(&m2)
    fmt.Printf("预分配: 分配次数=%d, GC次数=%d\n", 
        m2.Mallocs-m1.Mallocs, m2.NumGC-m1.NumGC)
}
```

---

## 优化技巧3：避免字符串拼接

### 问题：字符串拼接产生大量临时对象

```go
// 不好：每次 += 都创建新字符串
func bad() string {
    s := ""
    for i := 0; i < 10000; i++ {
        s += "a"  // 每次都分配新字符串
    }
    return s
}
```

### 解决：使用 strings.Builder

```go
// 好：使用 Builder
func good() string {
    var builder strings.Builder
    builder.Grow(10000)  // 预分配
    for i := 0; i < 10000; i++ {
        builder.WriteString("a")
    }
    return builder.String()
}
```

### 性能对比

```go
package main

import (
    "fmt"
    "runtime"
    "strings"
)

func concatWithPlus() string {
    s := ""
    for i := 0; i < 10000; i++ {
        s += "a"
    }
    return s
}

func concatWithBuilder() string {
    var builder strings.Builder
    for i := 0; i < 10000; i++ {
        builder.WriteString("a")
    }
    return builder.String()
}

func main() {
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    _ = concatWithPlus()
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("用 +: 分配=%d bytes\n", m2.TotalAlloc-m1.TotalAlloc)
    
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    _ = concatWithBuilder()
    
    runtime.ReadMemStats(&m2)
    fmt.Printf("用 Builder: 分配=%d bytes\n", m2.TotalAlloc-m1.TotalAlloc)
}
```

**结果：Builder 比 + 快 100 倍以上！**

---

## 优化技巧4：让变量分配在栈上

### 逃逸分析

```go
// 分配在栈上（不需要 GC）
func stackAlloc() int {
    x := 10  // x 不逃逸，分配在栈上
    return x
}

// 分配在堆上（需要 GC）
func heapAlloc() *int {
    x := 10   // x 逃逸到堆上
    return &x // 返回指针，x 必须在堆上
}
```

### 查看逃逸分析

```bash
go build -gcflags="-m" main.go
```

```go
package main

func stackAlloc() int {
    x := 10
    return x
}

func heapAlloc() *int {
    x := 10
    return &x
}

func main() {
    _ = stackAlloc()
    _ = heapAlloc()
}
```

**输出：**
```
./main.go:9:2: moved to heap: x    ← x 逃逸到堆
```

### 避免不必要的逃逸

```go
// 不好：返回指针，导致逃逸
func bad() *User {
    u := User{Name: "test"}
    return &u  // u 逃逸到堆
}

// 好：返回值，不逃逸
func good() User {
    u := User{Name: "test"}
    return u  // u 在栈上，返回时复制
}

// 或者：传入指针，让调用者决定分配位置
func better(u *User) {
    u.Name = "test"
}
```

---

## 优化技巧5：减少指针使用

### 问题：指针多 = GC 扫描工作量大

```go
// 不好：大量指针
type Bad struct {
    A *int
    B *string
    C *float64
}

// 好：直接存值
type Good struct {
    A int
    B string
    C float64
}
```

### 切片的选择

```go
// 不好：指针切片，每个元素都要扫描
users := make([]*User, 10000)

// 好：值切片，连续内存，扫描更快
users := make([]User, 10000)
```

---

## 优化技巧6：使用 bytes.Buffer 复用

```go
package main

import (
    "bytes"
    "sync"
)

var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processData(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()  // 重置，复用底层数组
    defer bufferPool.Put(buf)
    
    // 处理数据
    buf.Write(data)
    buf.WriteString(" processed")
    
    // 返回副本（因为 buf 会被复用）
    return append([]byte(nil), buf.Bytes()...)
}
```

---

## 实战：优化一个函数

### 优化前

```go
func processUsers(ids []int) []string {
    var results []string
    for _, id := range ids {
        user := fmt.Sprintf("user_%d", id)
        results = append(results, user)
    }
    return results
}
```

### 优化后

```go
func processUsersOptimized(ids []int) []string {
    // 1. 预分配切片
    results := make([]string, 0, len(ids))
    
    // 2. 复用 Builder
    var builder strings.Builder
    
    for _, id := range ids {
        builder.Reset()
        builder.WriteString("user_")
        builder.WriteString(strconv.Itoa(id))
        results = append(results, builder.String())
    }
    return results
}
```

### 性能测试

```go
package main

import (
    "fmt"
    "runtime"
    "strconv"
    "strings"
)

func processUsers(ids []int) []string {
    var results []string
    for _, id := range ids {
        user := fmt.Sprintf("user_%d", id)
        results = append(results, user)
    }
    return results
}

func processUsersOptimized(ids []int) []string {
    results := make([]string, 0, len(ids))
    var builder strings.Builder
    
    for _, id := range ids {
        builder.Reset()
        builder.WriteString("user_")
        builder.WriteString(strconv.Itoa(id))
        results = append(results, builder.String())
    }
    return results
}

func main() {
    ids := make([]int, 100000)
    for i := range ids {
        ids[i] = i
    }
    
    // 测试原版
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    _ = processUsers(ids)
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("原版: 分配=%dMB, GC=%d\n", 
        (m2.TotalAlloc-m1.TotalAlloc)/1024/1024, m2.NumGC-m1.NumGC)
    
    // 测试优化版
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    _ = processUsersOptimized(ids)
    
    runtime.ReadMemStats(&m2)
    fmt.Printf("优化: 分配=%dMB, GC=%d\n", 
        (m2.TotalAlloc-m1.TotalAlloc)/1024/1024, m2.NumGC-m1.NumGC)
}
```

---

## GC 调优参数

### GOGC

```bash
# 减少内存占用（GC 更频繁）
GOGC=50 ./myapp

# 减少 GC 次数（内存占用更高）
GOGC=200 ./myapp
```

### GOMEMLIMIT（Go 1.19+）

```bash
# 设置内存上限
GOMEMLIMIT=1GiB ./myapp
```

```go
import "runtime/debug"

// 代码中设置
debug.SetMemoryLimit(1 << 30)  // 1GB
```

---

## 监控 GC

### 使用 pprof

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // 你的程序...
}
```

**访问：**
- `http://localhost:6060/debug/pprof/heap` - 堆内存
- `http://localhost:6060/debug/pprof/allocs` - 分配统计

### 使用 trace

```bash
go test -trace trace.out
go tool trace trace.out
```

---

## 练习

### 练习1：优化这段代码

```go
func buildJSON(users []User) string {
    result := "["
    for i, u := range users {
        if i > 0 {
            result += ","
        }
        result += fmt.Sprintf(`{"name":"%s","age":%d}`, u.Name, u.Age)
    }
    result += "]"
    return result
}
```

<details>
<summary>点击查看答案</summary>

```go
func buildJSONOptimized(users []User) string {
    var builder strings.Builder
    builder.Grow(len(users) * 50)  // 预估大小
    
    builder.WriteByte('[')
    for i, u := range users {
        if i > 0 {
            builder.WriteByte(',')
        }
        builder.WriteString(`{"name":"`)
        builder.WriteString(u.Name)
        builder.WriteString(`","age":`)
        builder.WriteString(strconv.Itoa(u.Age))
        builder.WriteByte('}')
    }
    builder.WriteByte(']')
    
    return builder.String()
}
```

</details>

### 练习2：找出问题

```go
func getUsers() []*User {
    users := make([]*User, 1000)
    for i := 0; i < 1000; i++ {
        users[i] = &User{ID: i}
    }
    return users
}
```

问题在哪？怎么优化？

<details>
<summary>点击查看答案</summary>

**问题：** 创建了 1000 个指针，每个 User 都在堆上分配，GC 需要扫描 1000 个指针。

**优化：**
```go
func getUsersOptimized() []User {
    users := make([]User, 1000)  // 值切片，连续内存
    for i := 0; i < 1000; i++ {
        users[i] = User{ID: i}
    }
    return users
}
```

值切片在连续内存中，GC 扫描更快。

</details>

---

## 本章总结

| 优化技巧 | 方法 |
|---------|------|
| 对象复用 | sync.Pool |
| 预分配 | make([]T, 0, cap) |
| 字符串拼接 | strings.Builder |
| 避免逃逸 | 返回值而非指针 |
| 减少指针 | 值类型优于指针类型 |
| 复用 Buffer | bytes.Buffer + sync.Pool |

| 调优参数 | 说明 |
|---------|------|
| GOGC | 控制 GC 频率 |
| GOMEMLIMIT | 内存上限 |
| -gcflags="-m" | 查看逃逸分析 |
| GODEBUG=gctrace=1 | GC 日志 |

**记住：先测量，再优化。不要过早优化！**
