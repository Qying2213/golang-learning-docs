# Go 标准库教程

> 目标：熟悉 Go 最常用的标准库，能在实际开发中灵活运用

---

## 1. fmt - 格式化输入输出

最常用的库，没有之一。

### 1.1 输出函数

```go
package main

import "fmt"

func main() {
    name := "秦阳"
    age := 22
    
    // Print 系列 - 输出到控制台
    fmt.Print("不换行")
    fmt.Println("自动换行")
    fmt.Printf("格式化: %s, %d岁\n", name, age)
    
    // Sprint 系列 - 返回字符串（不输出）
    s := fmt.Sprintf("我是%s", name)
    fmt.Println(s)
    
    // Fprint 系列 - 输出到指定 Writer（文件等）
    // fmt.Fprintf(file, "写入文件: %s", name)
}
```

### 1.2 常用格式化占位符

```go
// 通用
%v   // 默认格式
%+v  // 结构体会打印字段名
%#v  // Go 语法格式
%T   // 类型

// 整数
%d   // 十进制
%b   // 二进制
%x   // 十六进制（小写）
%X   // 十六进制（大写）

// 浮点数
%f   // 小数，默认精度
%.2f // 保留2位小数
%e   // 科学计数法

// 字符串
%s   // 字符串
%q   // 带引号的字符串

// 布尔
%t   // true 或 false

// 指针
%p   // 指针地址
```

### 1.3 实际例子

```go
type User struct {
    Name string
    Age  int
}

func main() {
    u := User{"秦阳", 22}
    
    fmt.Printf("%v\n", u)   // {秦阳 22}
    fmt.Printf("%+v\n", u)  // {Name:秦阳 Age:22}
    fmt.Printf("%#v\n", u)  // main.User{Name:"秦阳", Age:22}
    
    price := 99.5
    fmt.Printf("价格: ￥%.2f\n", price)  // 价格: ￥99.50
}
```

---

## 2. strings - 字符串操作

### 2.1 常用函数

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, World"
    
    // 查找
    fmt.Println(strings.Contains(s, "World"))  // true - 是否包含
    fmt.Println(strings.HasPrefix(s, "Hello")) // true - 前缀
    fmt.Println(strings.HasSuffix(s, "World")) // true - 后缀
    fmt.Println(strings.Index(s, "o"))         // 4 - 第一次出现位置
    fmt.Println(strings.Count(s, "o"))         // 2 - 出现次数
    
    // 转换
    fmt.Println(strings.ToUpper(s))            // HELLO, WORLD
    fmt.Println(strings.ToLower(s))            // hello, world
    fmt.Println(strings.TrimSpace("  hi  "))   // "hi" - 去空格
    fmt.Println(strings.Trim("##hi##", "#"))   // "hi" - 去指定字符
    
    // 分割和拼接
    parts := strings.Split("a,b,c", ",")       // ["a", "b", "c"]
    joined := strings.Join(parts, "-")         // "a-b-c"
    
    // 替换
    fmt.Println(strings.Replace(s, "o", "0", 1))  // Hell0, World（替换1次）
    fmt.Println(strings.ReplaceAll(s, "o", "0")) // Hell0, W0rld（全部替换）
}
```

### 2.2 strings.Builder - 高效拼接

```go
// ❌ 低效：每次 += 都会创建新字符串
func bad() string {
    s := ""
    for i := 0; i < 1000; i++ {
        s += "a"
    }
    return s
}

// ✅ 高效：用 Builder
func good() string {
    var builder strings.Builder
    for i := 0; i < 1000; i++ {
        builder.WriteString("a")
    }
    return builder.String()
}
```

---

## 3. strconv - 字符串转换

### 3.1 字符串 ↔ 数字

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // 字符串 → 整数
    n, err := strconv.Atoi("123")
    if err != nil {
        fmt.Println("转换失败")
    }
    fmt.Println(n)  // 123
    
    // 整数 → 字符串
    s := strconv.Itoa(456)
    fmt.Println(s)  // "456"
    
    // 字符串 → 浮点数
    f, _ := strconv.ParseFloat("3.14", 64)
    fmt.Println(f)  // 3.14
    
    // 浮点数 → 字符串
    fs := strconv.FormatFloat(3.14159, 'f', 2, 64)
    fmt.Println(fs)  // "3.14"
    
    // 字符串 → 布尔
    b, _ := strconv.ParseBool("true")
    fmt.Println(b)  // true
}
```

### 3.2 记忆技巧

```
Atoi = Ascii to Integer（字符串转整数）
Itoa = Integer to Ascii（整数转字符串）
Parse = 解析（字符串 → 其他类型）
Format = 格式化（其他类型 → 字符串）
```

---

## 4. time - 时间处理

### 4.1 获取时间

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 当前时间
    now := time.Now()
    fmt.Println(now)  // 2024-12-24 10:30:00.123456 +0800 CST
    
    // 获取各部分
    fmt.Println(now.Year())        // 2024
    fmt.Println(now.Month())       // December
    fmt.Println(int(now.Month()))  // 12
    fmt.Println(now.Day())         // 24
    fmt.Println(now.Hour())        // 10
    fmt.Println(now.Minute())      // 30
    fmt.Println(now.Second())      // 0
    fmt.Println(now.Weekday())     // Tuesday
    
    // Unix 时间戳
    fmt.Println(now.Unix())        // 秒级时间戳
    fmt.Println(now.UnixMilli())   // 毫秒级时间戳
}
```

### 4.2 时间格式化（重要！）

Go 的时间格式化很特别，用的是 **固定的参考时间**：`2006-01-02 15:04:05`

```go
func main() {
    now := time.Now()
    
    // 格式化输出
    fmt.Println(now.Format("2006-01-02"))          // 2024-12-24
    fmt.Println(now.Format("2006/01/02 15:04:05")) // 2024/12/24 10:30:00
    fmt.Println(now.Format("15:04"))               // 10:30
    fmt.Println(now.Format("2006年01月02日"))       // 2024年12月24日
    
    // 解析字符串为时间
    t, _ := time.Parse("2006-01-02", "2024-12-25")
    fmt.Println(t)
}
```

**记忆口诀**：`2006-01-02 15:04:05` = 1月2日下午3点4分5秒2006年

### 4.3 时间计算

```go
func main() {
    now := time.Now()
    
    // 加减时间
    after := now.Add(2 * time.Hour)           // 2小时后
    before := now.Add(-30 * time.Minute)      // 30分钟前
    tomorrow := now.AddDate(0, 0, 1)          // 明天
    nextMonth := now.AddDate(0, 1, 0)         // 下个月
    
    // 计算时间差
    diff := after.Sub(now)
    fmt.Println(diff)          // 2h0m0s
    fmt.Println(diff.Hours())  // 2
    
    // 比较时间
    fmt.Println(after.After(now))   // true
    fmt.Println(before.Before(now)) // true
}
```

### 4.4 定时器和休眠

```go
func main() {
    // 休眠
    time.Sleep(1 * time.Second)
    
    // 定时器（一次性）
    timer := time.NewTimer(2 * time.Second)
    <-timer.C  // 2秒后收到信号
    fmt.Println("2秒到了")
    
    // Ticker（周期性）
    ticker := time.NewTicker(1 * time.Second)
    for i := 0; i < 3; i++ {
        <-ticker.C
        fmt.Println("tick", i)
    }
    ticker.Stop()
}
```

---

## 5. os - 操作系统交互

### 5.1 环境变量和参数

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // 命令行参数
    fmt.Println(os.Args)     // [程序名, 参数1, 参数2...]
    
    // 环境变量
    fmt.Println(os.Getenv("HOME"))        // 获取
    os.Setenv("MY_VAR", "hello")          // 设置
    
    // 当前工作目录
    dir, _ := os.Getwd()
    fmt.Println(dir)
    
    // 退出程序
    // os.Exit(1)  // 非0表示异常退出
}
```

### 5.2 文件操作

```go
func main() {
    // 创建文件
    file, err := os.Create("test.txt")
    if err != nil {
        panic(err)
    }
    file.WriteString("Hello, Go!")
    file.Close()
    
    // 读取文件（简单方式）
    data, _ := os.ReadFile("test.txt")
    fmt.Println(string(data))
    
    // 写入文件（简单方式）
    os.WriteFile("test2.txt", []byte("内容"), 0644)
    
    // 删除文件
    os.Remove("test.txt")
    
    // 创建目录
    os.Mkdir("mydir", 0755)
    os.MkdirAll("a/b/c", 0755)  // 递归创建
    
    // 删除目录
    os.RemoveAll("mydir")
    
    // 判断文件是否存在
    _, err = os.Stat("test.txt")
    if os.IsNotExist(err) {
        fmt.Println("文件不存在")
    }
}
```

---

## 6. io 和 bufio - 输入输出

### 6.1 io 基础

```go
package main

import (
    "io"
    "os"
    "strings"
)

func main() {
    // 从 Reader 读取全部
    reader := strings.NewReader("Hello, Go!")
    data, _ := io.ReadAll(reader)
    fmt.Println(string(data))
    
    // 复制
    src, _ := os.Open("source.txt")
    dst, _ := os.Create("dest.txt")
    io.Copy(dst, src)
    src.Close()
    dst.Close()
}
```

### 6.2 bufio - 带缓冲的读写

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // 按行读取文件
    file, _ := os.Open("test.txt")
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
    }
    
    // 从控制台读取输入
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("请输入: ")
    input, _ := reader.ReadString('\n')
    fmt.Println("你输入了:", input)
}
```

---

## 7. encoding/json - JSON 处理

### 7.1 结构体 ↔ JSON

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name  string `json:"name"`           // 指定 JSON 字段名
    Age   int    `json:"age"`
    Email string `json:"email,omitempty"` // 空值时省略
}

func main() {
    // 结构体 → JSON（序列化）
    u := User{Name: "秦阳", Age: 22}
    data, _ := json.Marshal(u)
    fmt.Println(string(data))  // {"name":"秦阳","age":22}
    
    // 格式化输出
    pretty, _ := json.MarshalIndent(u, "", "  ")
    fmt.Println(string(pretty))
    
    // JSON → 结构体（反序列化）
    jsonStr := `{"name":"张三","age":25,"email":"test@example.com"}`
    var u2 User
    json.Unmarshal([]byte(jsonStr), &u2)
    fmt.Printf("%+v\n", u2)  // {Name:张三 Age:25 Email:test@example.com}
}
```

### 7.2 处理动态 JSON

```go
func main() {
    jsonStr := `{"name":"秦阳","scores":[90,85,92]}`
    
    // 用 map 接收
    var data map[string]interface{}
    json.Unmarshal([]byte(jsonStr), &data)
    
    fmt.Println(data["name"])    // 秦阳
    fmt.Println(data["scores"])  // [90 85 92]
    
    // 类型断言获取具体值
    if name, ok := data["name"].(string); ok {
        fmt.Println("名字是:", name)
    }
}
```

---

## 8. net/http - HTTP 客户端和服务器

### 8.1 HTTP 客户端

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    // 简单 GET 请求
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
    
    // 检查状态码
    fmt.Println("状态码:", resp.StatusCode)
}
```

### 8.2 POST 请求

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

func main() {
    // POST JSON
    data := map[string]string{"name": "秦阳"}
    jsonData, _ := json.Marshal(data)
    
    resp, _ := http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    defer resp.Body.Close()
    
    // 处理响应...
}
```

### 8.3 HTTP 服务器

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    // 处理函数
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
    })
    
    http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
        user := map[string]interface{}{
            "name": "秦阳",
            "age":  22,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    })
    
    // 启动服务器
    fmt.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## 9. path/filepath - 路径处理

```go
package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    // 拼接路径（自动处理分隔符）
    p := filepath.Join("home", "user", "docs", "file.txt")
    fmt.Println(p)  // home/user/docs/file.txt
    
    // 获取目录和文件名
    fmt.Println(filepath.Dir(p))   // home/user/docs
    fmt.Println(filepath.Base(p))  // file.txt
    fmt.Println(filepath.Ext(p))   // .txt
    
    // 获取绝对路径
    abs, _ := filepath.Abs(".")
    fmt.Println(abs)
    
    // 遍历目录
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        fmt.Println(path)
        return nil
    })
}
```

---

## 10. sort - 排序

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // 基本类型排序
    nums := []int{3, 1, 4, 1, 5, 9}
    sort.Ints(nums)
    fmt.Println(nums)  // [1 1 3 4 5 9]
    
    strs := []string{"banana", "apple", "cherry"}
    sort.Strings(strs)
    fmt.Println(strs)  // [apple banana cherry]
    
    // 自定义排序
    type Person struct {
        Name string
        Age  int
    }
    
    people := []Person{
        {"秦阳", 22},
        {"张三", 25},
        {"李四", 20},
    }
    
    // 按年龄排序
    sort.Slice(people, func(i, j int) bool {
        return people[i].Age < people[j].Age
    })
    fmt.Println(people)  // [{李四 20} {秦阳 22} {张三 25}]
    
    // 检查是否已排序
    fmt.Println(sort.IntsAreSorted(nums))  // true
}
```

---

## 11. regexp - 正则表达式

```go
package main

import (
    "fmt"
    "regexp"
)

func main() {
    // 编译正则
    re := regexp.MustCompile(`\d+`)
    
    // 查找
    fmt.Println(re.FindString("abc123def"))      // 123
    fmt.Println(re.FindAllString("a1b2c3", -1))  // [1 2 3]
    
    // 匹配
    fmt.Println(re.MatchString("hello123"))  // true
    
    // 替换
    result := re.ReplaceAllString("a1b2c3", "X")
    fmt.Println(result)  // aXbXcX
    
    // 验证邮箱
    emailRe := regexp.MustCompile(`^[\w.-]+@[\w.-]+\.\w+$`)
    fmt.Println(emailRe.MatchString("test@example.com"))  // true
    
    // 提取分组
    re2 := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
    matches := re2.FindStringSubmatch("test@example.com")
    fmt.Println(matches)  // [test@example.com test example com]
}
```

---

## 12. log - 日志

```go
package main

import (
    "log"
    "os"
)

func main() {
    // 基本使用
    log.Println("普通日志")
    log.Printf("格式化: %s\n", "hello")
    
    // 设置前缀和格式
    log.SetPrefix("[APP] ")
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Println("带前缀的日志")
    // 输出: [APP] 2024/12/24 10:30:00 main.go:15: 带前缀的日志
    
    // 输出到文件
    file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    log.SetOutput(file)
    log.Println("写入文件")
    
    // 致命错误（会调用 os.Exit(1)）
    // log.Fatal("致命错误")
    
    // Panic（会触发 panic）
    // log.Panic("panic!")
}
```

---

## 13. errors - 错误处理

```go
package main

import (
    "errors"
    "fmt"
)

// 定义错误
var ErrNotFound = errors.New("not found")
var ErrPermission = errors.New("permission denied")

func findUser(id int) (string, error) {
    if id <= 0 {
        return "", ErrNotFound
    }
    return "秦阳", nil
}

func main() {
    // 创建错误
    err := errors.New("something went wrong")
    fmt.Println(err)
    
    // 格式化错误
    err2 := fmt.Errorf("user %d not found", 123)
    fmt.Println(err2)
    
    // 错误判断
    _, err = findUser(0)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("用户不存在")
    }
    
    // 错误包装（Go 1.13+）
    wrapped := fmt.Errorf("查询失败: %w", ErrNotFound)
    fmt.Println(errors.Is(wrapped, ErrNotFound))  // true
    fmt.Println(errors.Unwrap(wrapped))           // not found
}
```

---

## 14. context - 上下文控制

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 带超时的 context
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // 模拟耗时操作
    select {
    case <-time.After(3 * time.Second):
        fmt.Println("操作完成")
    case <-ctx.Done():
        fmt.Println("超时了:", ctx.Err())
    }
    
    // 带取消的 context
    ctx2, cancel2 := context.WithCancel(context.Background())
    go func() {
        time.Sleep(1 * time.Second)
        cancel2()  // 手动取消
    }()
    <-ctx2.Done()
    fmt.Println("被取消了")
    
    // 带值的 context
    ctx3 := context.WithValue(context.Background(), "userID", 123)
    fmt.Println(ctx3.Value("userID"))  // 123
}
```

---

## 练习题

### 练习 1：字符串处理
```go
// 实现一个函数，统计字符串中每个单词出现的次数
// 输入: "hello world hello go world world"
// 输出: map[string]int{"hello": 2, "world": 3, "go": 1}
func WordCount(s string) map[string]int {
    // 你来实现
}
```

### 练习 2：时间计算
```go
// 计算两个日期之间相差多少天
// 输入: "2024-01-01", "2024-12-25"
// 输出: 359
func DaysBetween(date1, date2 string) int {
    // 你来实现
}
```

### 练习 3：JSON 处理
```go
// 解析 JSON 并提取信息
// 输入: `{"users":[{"name":"秦阳","age":22},{"name":"张三","age":25}]}`
// 输出: 所有用户的名字
func GetUserNames(jsonStr string) []string {
    // 你来实现
}
```

### 练习 4：HTTP 服务
```go
// 写一个简单的 HTTP 服务器
// GET /hello?name=xxx 返回 "Hello, xxx!"
// GET /time 返回当前时间
func main() {
    // 你来实现
}
```

### 练习 5：文件操作
```go
// 读取一个文件，统计行数、单词数、字符数
// 类似 Linux 的 wc 命令
func FileStats(filename string) (lines, words, chars int) {
    // 你来实现
}
```

---

## 标准库速查表

| 库 | 用途 | 常用函数 |
|---|---|---|
| fmt | 格式化输出 | Printf, Sprintf, Println |
| strings | 字符串操作 | Contains, Split, Join, Replace |
| strconv | 类型转换 | Atoi, Itoa, ParseFloat |
| time | 时间处理 | Now, Format, Parse, Sleep |
| os | 系统操作 | Open, Create, ReadFile, Getenv |
| io | 输入输出 | ReadAll, Copy |
| bufio | 缓冲读写 | NewScanner, NewReader |
| encoding/json | JSON处理 | Marshal, Unmarshal |
| net/http | HTTP | Get, Post, ListenAndServe |
| filepath | 路径处理 | Join, Dir, Base, Walk |
| sort | 排序 | Ints, Strings, Slice |
| regexp | 正则表达式 | MatchString, FindString |
| log | 日志 | Println, Printf, Fatal |
| errors | 错误处理 | New, Is, Unwrap |
| context | 上下文 | WithTimeout, WithCancel |

---

## 学习建议

1. **先掌握这几个**：fmt, strings, strconv, time, encoding/json
2. **写项目时学**：net/http, os, io, context
3. **遇到再查**：regexp, sort, filepath

多写代码，遇到问题查文档：https://pkg.go.dev/std
