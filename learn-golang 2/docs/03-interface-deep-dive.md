# Interface 深入理解

## 1. 接口的底层结构

Go 中有两种接口：

```go
// 空接口 (interface{} 或 any)
type eface struct {
    _type *_type         // 类型信息
    data  unsafe.Pointer // 数据指针
}

// 非空接口
type iface struct {
    tab  *itab           // 类型和方法表
    data unsafe.Pointer  // 数据指针
}

type itab struct {
    inter *interfacetype // 接口类型
    _type *_type         // 具体类型
    hash  uint32         // 类型 hash，用于快速判断
    _     [4]byte
    fun   [1]uintptr     // 方法表（变长数组）
}
```

## 2. 接口的隐式实现

```go
// 定义接口
type Writer interface {
    Write([]byte) (int, error)
}

// 实现接口（无需显式声明）
type MyWriter struct{}

func (w MyWriter) Write(data []byte) (int, error) {
    fmt.Println(string(data))
    return len(data), nil
}

// 编译时检查是否实现接口
var _ Writer = (*MyWriter)(nil)
var _ Writer = MyWriter{}
```

## 3. 值接收者 vs 指针接收者

```go
type Counter interface {
    Increment()
    Value() int
}

type MyCounter struct {
    count int
}

// 值接收者
func (c MyCounter) Value() int {
    return c.count
}

// 指针接收者
func (c *MyCounter) Increment() {
    c.count++
}

func main() {
    // 指针类型实现了所有方法
    var c1 Counter = &MyCounter{}  // OK
    
    // 值类型只实现了值接收者的方法
    // var c2 Counter = MyCounter{}  // 编译错误！
    // 因为 MyCounter 没有实现 Increment（指针接收者）
}
```

**规则：**
- 值接收者：值和指针都能调用
- 指针接收者：只有指针能调用
- 接口变量存储的是值的副本，如果方法需要修改，必须用指针接收者

## 4. 接口的零值

```go
// 接口的零值是 nil
var w Writer  // w == nil

// nil 接口调用方法会 panic
// w.Write([]byte("hello"))  // panic!

// 但是：接口值为 nil 和接口内部数据为 nil 是不同的！
var p *MyWriter = nil
var w Writer = p

fmt.Println(w == nil)  // false！
// 因为 w 的类型信息不为空，只是数据为空

// 这是常见的坑
func returnsNil() Writer {
    var p *MyWriter = nil
    return p  // 返回的不是 nil 接口！
}

w := returnsNil()
fmt.Println(w == nil)  // false
```

## 5. 类型断言

```go
var i interface{} = "hello"

// 方式1: 直接断言（失败会 panic）
s := i.(string)
fmt.Println(s)

// 方式2: 安全断言
s, ok := i.(string)
if ok {
    fmt.Println(s)
}

// 方式3: type switch
switch v := i.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
case nil:
    fmt.Println("nil")
default:
    fmt.Printf("unknown: %T\n", v)
}
```

## 6. 接口组合

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// 接口组合
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// 标准库中的例子
// io.ReadWriter, io.ReadCloser, io.WriteCloser, io.ReadWriteCloser
```

## 7. 空接口的使用

```go
// 空接口可以存储任何类型
var any interface{}
any = 42
any = "hello"
any = []int{1, 2, 3}

// Go 1.18+ 可以用 any 代替 interface{}
var a any = "hello"

// 常见用途
func Print(args ...interface{}) {
    for _, arg := range args {
        fmt.Println(arg)
    }
}

// JSON 解析
var data map[string]interface{}
json.Unmarshal([]byte(`{"name":"Go","age":15}`), &data)
```

## 8. 接口设计原则

```go
// 原则1: 接口应该小而专注
// 好的设计
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 不好的设计
type FileOperator interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
    Close() error
    Seek(offset int64, whence int) (int64, error)
    Stat() (FileInfo, error)
    // ... 太多方法
}

// 原则2: 在使用方定义接口，而不是实现方
// 不好：在实现方定义
package mydb
type Database interface {
    Query(sql string) ([]Row, error)
    Exec(sql string) error
}
type MySQL struct{}
func (m *MySQL) Query(sql string) ([]Row, error) { ... }

// 好：在使用方定义
package myservice
type Querier interface {
    Query(sql string) ([]Row, error)
}
func NewService(db Querier) *Service { ... }

// 原则3: 返回具体类型，接受接口类型
func NewBuffer() *Buffer { ... }  // 返回具体类型
func Copy(dst Writer, src Reader) { ... }  // 接受接口
```

## 9. 接口的性能

```go
// 接口调用有一定开销（动态分派）
// 但通常可以忽略

// 避免不必要的接口装箱
func process(data interface{}) {
    // 每次调用都会装箱
}

// 如果类型确定，使用具体类型
func processInt(data int) {
    // 无装箱开销
}

// 使用泛型（Go 1.18+）避免接口开销
func processSlice[T any](data []T) {
    // 编译时确定类型，无运行时开销
}
```

## 10. 常见面试题

**Q1: 接口值什么时候等于 nil？**
```go
// 只有当类型和值都为 nil 时，接口才等于 nil
var i interface{} = nil  // i == nil: true

var p *int = nil
var i interface{} = p    // i == nil: false
```

**Q2: 如何判断接口内部的值是否为 nil？**
```go
func isNil(i interface{}) bool {
    if i == nil {
        return true
    }
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
        return v.IsNil()
    }
    return false
}
```

**Q3: 为什么 Go 没有泛型时大量使用 interface{}？**
```go
// 因为 interface{} 可以存储任何类型
// 但有类型安全问题，需要运行时断言
// Go 1.18 引入泛型后，很多场景可以用泛型替代
```

## 11. 实用模式

```go
// 模式1: 接口检查
type Validator interface {
    Validate() error
}

func validate(v interface{}) error {
    if validator, ok := v.(Validator); ok {
        return validator.Validate()
    }
    return nil
}

// 模式2: 可选接口
type Handler interface {
    Handle(ctx context.Context) error
}

type Initializer interface {
    Init() error
}

func run(h Handler) error {
    // 检查是否实现了可选接口
    if init, ok := h.(Initializer); ok {
        if err := init.Init(); err != nil {
            return err
        }
    }
    return h.Handle(context.Background())
}

// 模式3: 函数类型实现接口
type HandlerFunc func(ctx context.Context) error

func (f HandlerFunc) Handle(ctx context.Context) error {
    return f(ctx)
}

// 使用
var h Handler = HandlerFunc(func(ctx context.Context) error {
    return nil
})
```

## 练习

1. 实现一个 `Stringer` 接口，让自定义类型可以打印
2. 实现一个简单的依赖注入容器
3. 设计一个插件系统，使用接口定义插件规范
4. 实现一个 `io.Reader` 包装器，统计读取的字节数
