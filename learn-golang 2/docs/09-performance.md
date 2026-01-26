# 性能分析与优化

## 1. pprof 基础

```go
import (
    "net/http"
    _ "net/http/pprof"  // 自动注册 pprof 路由
)

func main() {
    // 启动 pprof 服务
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // 你的应用代码...
}

// 访问 http://localhost:6060/debug/pprof/
```

## 2. CPU 分析

```bash
# 采集 30 秒 CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 常用命令
(pprof) top        # 查看 CPU 占用最高的函数
(pprof) top10      # 前 10 个
(pprof) list main  # 查看 main 函数的详细信息
(pprof) web        # 在浏览器中查看调用图

# 命令行直接查看
go tool pprof -top http://localhost:6060/debug/pprof/profile?seconds=30
```

## 3. 内存分析

```bash
# 堆内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看内存分配
(pprof) top
(pprof) list functionName

# 查看 inuse（当前使用）vs alloc（总分配）
go tool pprof -inuse_space http://localhost:6060/debug/pprof/heap
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap
```

## 4. Goroutine 分析

```bash
# 查看 goroutine 状态
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 查看阻塞情况
go tool pprof http://localhost:6060/debug/pprof/block

# 查看互斥锁竞争
go tool pprof http://localhost:6060/debug/pprof/mutex
```

## 5. Trace 分析

```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // 你的代码...
}

// 或者通过 HTTP
// curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
// go tool trace trace.out
```


## 6. Benchmark 测试

```go
func BenchmarkSliceAppend(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var s []int
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

func BenchmarkSlicePrealloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := make([]int, 0, 1000)
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

// 运行
// go test -bench=. -benchmem
// BenchmarkSliceAppend-8      10000    150000 ns/op    40000 B/op    20 allocs/op
// BenchmarkSlicePrealloc-8    50000     30000 ns/op     8192 B/op     1 allocs/op
```

## 7. 常见优化技巧

### 内存优化
```go
// 1. 预分配切片
s := make([]int, 0, expectedSize)

// 2. 使用 sync.Pool 复用对象
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()
    // 使用 buf...
}

// 3. 避免不必要的内存分配
// 不好
func concat(a, b string) string {
    return a + b  // 每次都分配新内存
}

// 好
func concat(a, b string) string {
    var builder strings.Builder
    builder.WriteString(a)
    builder.WriteString(b)
    return builder.String()
}

// 4. 使用指针避免复制大结构体
type BigStruct struct {
    data [1024]byte
}

func process(b *BigStruct) {}  // 好
func process(b BigStruct) {}   // 不好，会复制
```

### CPU 优化
```go
// 1. 减少锁竞争
// 使用 sync.RWMutex 代替 sync.Mutex
// 使用分片减少锁粒度

// 2. 使用 atomic 代替锁
var counter int64
atomic.AddInt64(&counter, 1)

// 3. 避免反射
// 反射很慢，热路径上避免使用

// 4. 内联小函数
// 编译器会自动内联小函数
// 可以用 //go:noinline 禁止内联（调试用）
```

## 8. 逃逸分析

```go
// 查看逃逸分析结果
// go build -gcflags="-m" main.go

func noEscape() int {
    x := 42  // 栈上分配
    return x
}

func escape() *int {
    x := 42   // 逃逸到堆上
    return &x // 返回指针导致逃逸
}

// 减少逃逸的技巧
// 1. 避免返回局部变量的指针
// 2. 避免在闭包中捕获局部变量
// 3. 使用值传递而不是指针（小结构体）
```

## 9. 编译优化

```bash
# 查看编译器优化
go build -gcflags="-m -m" main.go

# 禁用优化（调试用）
go build -gcflags="-N -l" main.go

# 查看汇编
go build -gcflags="-S" main.go

# 查看内联决策
go build -gcflags="-m=2" main.go
```

## 10. 实战优化案例

```go
// 案例1: JSON 序列化优化
// 使用 json-iterator 代替标准库
import jsoniter "github.com/json-iterator/go"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 案例2: 字符串拼接
// 不好
s := ""
for i := 0; i < 1000; i++ {
    s += strconv.Itoa(i)
}

// 好
var builder strings.Builder
for i := 0; i < 1000; i++ {
    builder.WriteString(strconv.Itoa(i))
}
s := builder.String()

// 案例3: 避免 interface{} 装箱
// 不好
func sum(nums []interface{}) int {
    total := 0
    for _, n := range nums {
        total += n.(int)
    }
    return total
}

// 好
func sum(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 或使用泛型
func sum[T int | int64 | float64](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

## 11. 性能检查清单

```
□ 是否预分配了切片和 map？
□ 是否使用了 sync.Pool 复用对象？
□ 是否避免了不必要的内存分配？
□ 是否使用了合适的数据结构？
□ 是否减少了锁竞争？
□ 是否避免了热路径上的反射？
□ 是否使用了 pprof 分析过？
□ 是否有 benchmark 测试？
```

## 练习

1. 使用 pprof 分析一个 HTTP 服务的性能
2. 优化一个字符串处理函数，减少内存分配
3. 实现一个对象池，对比使用前后的性能
4. 分析一个程序的逃逸情况，尝试减少堆分配
