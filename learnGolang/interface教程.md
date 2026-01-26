# Go Interface 完全教程

> 从零开始理解 interface，适合实习生

---

## 一、什么是 Interface

**一句话：Interface 是一组方法的集合。**

```go
// 定义一个接口
type Speaker interface {
    Speak() string  // 只要有 Speak() 方法，就实现了这个接口
}
```

**关键点：Go 的接口是隐式实现的，不需要写 `implements`。**

```go
// 定义一个结构体
type Dog struct {
    Name string
}

// Dog 实现了 Speak() 方法，自动就实现了 Speaker 接口
func (d Dog) Speak() string {
    return d.Name + ": 汪汪汪"
}

type Cat struct {
    Name string
}

// Cat 也实现了 Speak() 方法
func (c Cat) Speak() string {
    return c.Name + ": 喵喵喵"
}
```

**使用：**

```go
func main() {
    // 接口变量可以存储任何实现了该接口的类型
    var s Speaker
    
    s = Dog{Name: "旺财"}
    fmt.Println(s.Speak())  // 旺财: 汪汪汪
    
    s = Cat{Name: "咪咪"}
    fmt.Println(s.Speak())  // 咪咪: 喵喵喵
}
```

---

## 二、为什么需要 Interface

### 2.1 不用接口的问题

```go
// 假设你要写一个函数，让动物叫
func MakeDogSpeak(d Dog) {
    fmt.Println(d.Speak())
}

func MakeCatSpeak(c Cat) {
    fmt.Println(c.Speak())
}

// 如果有 100 种动物，你要写 100 个函数？
```

### 2.2 用接口解决

```go
// 一个函数搞定所有
func MakeSpeak(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    MakeSpeak(Dog{Name: "旺财"})  // 旺财: 汪汪汪
    MakeSpeak(Cat{Name: "咪咪"})  // 咪咪: 喵喵喵
    // 以后加新动物，这个函数不用改
}
```

**接口的好处：解耦，让代码更灵活。**

---

## 三、常见的内置接口

### 3.1 Stringer（最常用）

```go
// fmt 包里定义的
type Stringer interface {
    String() string
}
```

**实现 Stringer 后，fmt.Println 会自动调用你的 String() 方法：**

```go
type Person struct {
    Name string
    Age  int
}

// 实现 Stringer 接口
func (p Person) String() string {
    return fmt.Sprintf("姓名: %s, 年龄: %d", p.Name, p.Age)
}

func main() {
    p := Person{"秦阳", 22}
    fmt.Println(p)  // 姓名: 秦阳, 年龄: 22（自动调用 String()）
}
```

---

### 3.2 error 接口

```go
// 内置的 error 接口
type error interface {
    Error() string
}
```

**自定义错误：**

```go
type MyError struct {
    Code    int
    Message string
}

// 实现 error 接口
func (e MyError) Error() string {
    return fmt.Sprintf("错误码: %d, 信息: %s", e.Code, e.Message)
}

func doSomething() error {
    return MyError{Code: 1001, Message: "出错了"}
}

func main() {
    err := doSomething()
    if err != nil {
        fmt.Println(err)  // 错误码: 1001, 信息: 出错了
    }
}
```

---

### 3.3 io.Reader 和 io.Writer

```go
// io 包里定义的
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**很多类型都实现了这两个接口：**
- `os.File`（文件）
- `bytes.Buffer`（内存缓冲区）
- `strings.Reader`（字符串）
- `net.Conn`（网络连接）

**好处：你可以用同样的代码处理不同的数据源。**

```go
// 这个函数可以从任何 Reader 读取数据
func ReadAll(r io.Reader) ([]byte, error) {
    return io.ReadAll(r)
}

func main() {
    // 从字符串读
    data1, _ := ReadAll(strings.NewReader("hello"))
    
    // 从文件读
    file, _ := os.Open("test.txt")
    data2, _ := ReadAll(file)
}
```

---

## 四、空接口 interface{}

**空接口没有任何方法，所以任何类型都实现了它。**

```go
var anything interface{}

anything = 42
anything = "hello"
anything = []int{1, 2, 3}
anything = Dog{Name: "旺财"}

// Go 1.18+ 可以用 any 代替 interface{}
var anything2 any
```

**常见用途：**

```go
// 1. 接收任意类型的参数
func Print(v interface{}) {
    fmt.Println(v)
}

// 2. 存储任意类型的值
data := map[string]interface{}{
    "name": "秦阳",
    "age":  22,
    "scores": []int{90, 85, 88},
}
```

---

## 五、类型断言

**从接口变量中取出具体类型的值。**

### 5.1 基本用法

```go
var s Speaker = Dog{Name: "旺财"}

// 类型断言
dog := s.(Dog)  // 把 s 断言为 Dog 类型
fmt.Println(dog.Name)  // 旺财
```

### 5.2 安全的类型断言

```go
var s Speaker = Dog{Name: "旺财"}

// 如果断言错误，会 panic
// cat := s.(Cat)  // panic!

// 安全写法：用两个返回值
dog, ok := s.(Dog)
if ok {
    fmt.Println("是狗:", dog.Name)
} else {
    fmt.Println("不是狗")
}
```

### 5.3 type switch

```go
func WhatIsIt(v interface{}) {
    switch x := v.(type) {
    case int:
        fmt.Println("是整数:", x)
    case string:
        fmt.Println("是字符串:", x)
    case Dog:
        fmt.Println("是狗:", x.Name)
    default:
        fmt.Printf("未知类型: %T\n", x)
    }
}

func main() {
    WhatIsIt(42)           // 是整数: 42
    WhatIsIt("hello")      // 是字符串: hello
    WhatIsIt(Dog{"旺财"})  // 是狗: 旺财
    WhatIsIt(3.14)         // 未知类型: float64
}
```

---

## 六、值接收者 vs 指针接收者

**这是一个重要的坑！**

```go
type Speaker interface {
    Speak() string
}

type Dog struct {
    Name string
}

// 值接收者
func (d Dog) Speak() string {
    return d.Name + ": 汪汪汪"
}
```

**值接收者：值和指针都能用**

```go
var s Speaker

s = Dog{Name: "旺财"}   // ✅ OK
s = &Dog{Name: "旺财"}  // ✅ OK
```

---

```go
type Cat struct {
    Name string
}

// 指针接收者
func (c *Cat) Speak() string {
    return c.Name + ": 喵喵喵"
}
```

**指针接收者：只有指针能用**

```go
var s Speaker

s = &Cat{Name: "咪咪"}  // ✅ OK
s = Cat{Name: "咪咪"}   // ❌ 编译错误！
```

**记忆口诀：**
```
值接收者：值和指针都行
指针接收者：只能用指针
```

---

## 七、接口组合

**小接口组合成大接口。**

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}
```

**Go 推荐小接口：**

```go
// 好：小而专一
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 不好：太大太杂
type FileOperator interface {
    Read() error
    Write() error
    Close() error
    Seek() error
    Stat() error
    // ...
}
```

---

## 八、接口的零值是 nil

```go
var s Speaker  // nil

if s == nil {
    fmt.Println("s 是 nil")
}

// 调用 nil 接口的方法会 panic
// s.Speak()  // panic: nil pointer dereference
```

**注意这个坑：**

```go
var dog *Dog = nil
var s Speaker = dog

// s 不是 nil！
fmt.Println(s == nil)  // false

// 但调用方法还是会 panic
// s.Speak()  // panic
```

---

## 九、实战练习

### 练习1：实现 Stringer ⭐

```go
type Student struct {
    Name  string
    Score int
}

// 实现 String() 方法，让 fmt.Println 输出：
// "学生: xxx, 成绩: xxx"

func main() {
    s := Student{"秦阳", 95}
    fmt.Println(s)  // 学生: 秦阳, 成绩: 95
}
```

---

### 练习2：形状接口 ⭐⭐

```go
type Shape interface {
    Area() float64
}

// 实现 Rectangle 和 Circle，都实现 Shape 接口

type Rectangle struct {
    Width, Height float64
}

type Circle struct {
    Radius float64
}

// 写一个函数，计算多个形状的总面积
func TotalArea(shapes []Shape) float64 {
    // 你来实现
}

func main() {
    shapes := []Shape{
        Rectangle{10, 5},
        Circle{3},
        Rectangle{4, 6},
    }
    fmt.Println(TotalArea(shapes))  // 50 + 28.27 + 24 = 102.27
}
```

---

### 练习3：类型断言 ⭐⭐

```go
// 写一个函数，接收 interface{}
// 如果是 int，返回它的两倍
// 如果是 string，返回它的长度
// 其他类型返回 -1

func Process(v interface{}) int {
    // 你来实现
}

func main() {
    fmt.Println(Process(10))      // 20
    fmt.Println(Process("hello")) // 5
    fmt.Println(Process(3.14))    // -1
}
```

---

## 十、总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Interface 要点                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   1. 接口是方法的集合                                                        │
│                                                                             │
│   2. 隐式实现，不需要 implements                                             │
│                                                                             │
│   3. 空接口 interface{} / any 可以存任何类型                                 │
│                                                                             │
│   4. 类型断言：v.(Type) 或 v.(type)                                          │
│                                                                             │
│   5. 值接收者：值和指针都行                                                   │
│      指针接收者：只能用指针                                                   │
│                                                                             │
│   6. 接口的零值是 nil                                                        │
│                                                                             │
│   7. 推荐小接口，可以组合                                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**常用内置接口：**
- `Stringer` - 自定义打印格式
- `error` - 错误处理
- `io.Reader` / `io.Writer` - IO 操作
- `sort.Interface` - 排序
