# Go 错误处理教程

---

## 1. Go 的错误处理方式

Go 没有 try-catch，用**返回值**来处理错误。

```go
// 其他语言
try {
    result = doSomething()
} catch (error) {
    // 处理错误
}

// Go
result, err := doSomething()
if err != nil {
    // 处理错误
}
```

---

## 2. 基本用法

### 2.1 检查错误

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("不存在的文件.txt")
    if err != nil {
        fmt.Println("打开失败:", err)
        return
    }
    defer file.Close()
    
    fmt.Println("打开成功")
}
```

输出：
```
打开失败: open 不存在的文件.txt: no such file or directory
```

### 2.2 error 是什么

`error` 是一个接口，只有一个方法：

```go
type error interface {
    Error() string
}
```

任何实现了 `Error()` 方法的类型都是 error。

---

## 3. 创建错误

### 3.1 errors.New

最简单的方式：

```go
import "errors"

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为0")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("错误:", err)
        return
    }
    fmt.Println("结果:", result)
}
```

### 3.2 fmt.Errorf

可以格式化错误信息：

```go
func getUser(id int) (string, error) {
    if id <= 0 {
        return "", fmt.Errorf("无效的用户ID: %d", id)
    }
    return "秦阳", nil
}

func main() {
    user, err := getUser(-1)
    if err != nil {
        fmt.Println(err)  // 无效的用户ID: -1
    }
}
```

---

## 4. 错误处理模式

### 4.1 直接返回

```go
func doA() error {
    return errors.New("A 出错了")
}

func doB() error {
    err := doA()
    if err != nil {
        return err  // 直接返回
    }
    return nil
}
```

### 4.2 包装错误（推荐）

用 `%w` 包装错误，保留原始错误信息：

```go
func readConfig() error {
    _, err := os.Open("config.json")
    if err != nil {
        return fmt.Errorf("读取配置失败: %w", err)
    }
    return nil
}

func main() {
    err := readConfig()
    if err != nil {
        fmt.Println(err)
        // 输出: 读取配置失败: open config.json: no such file or directory
    }
}
```

### 4.3 多个错误检查

```go
func process() error {
    data, err := readFile()
    if err != nil {
        return fmt.Errorf("读取文件失败: %w", err)
    }
    
    result, err := parseData(data)
    if err != nil {
        return fmt.Errorf("解析数据失败: %w", err)
    }
    
    err = saveResult(result)
    if err != nil {
        return fmt.Errorf("保存结果失败: %w", err)
    }
    
    return nil
}
```

---

## 5. 预定义错误

定义一些常用的错误，方便判断：

```go
var (
    ErrNotFound     = errors.New("未找到")
    ErrUnauthorized = errors.New("未授权")
    ErrInvalidInput = errors.New("输入无效")
)

func findUser(id int) (string, error) {
    if id <= 0 {
        return "", ErrInvalidInput
    }
    if id == 999 {
        return "", ErrNotFound
    }
    return "秦阳", nil
}

func main() {
    _, err := findUser(999)
    if err == ErrNotFound {
        fmt.Println("用户不存在")
    }
}
```

---

## 6. errors.Is 和 errors.As

### 6.1 errors.Is - 判断错误类型

```go
var ErrNotFound = errors.New("未找到")

func getUser(id int) error {
    return fmt.Errorf("查询用户失败: %w", ErrNotFound)
}

func main() {
    err := getUser(1)
    
    // ❌ 直接比较不行，因为被包装了
    if err == ErrNotFound {
        fmt.Println("不会执行")
    }
    
    // ✅ 用 errors.Is，能穿透包装
    if errors.Is(err, ErrNotFound) {
        fmt.Println("用户不存在")  // 会执行
    }
}
```

### 6.2 errors.As - 获取特定类型的错误

```go
// 自定义错误类型
type MyError struct {
    Code    int
    Message string
}

func (e *MyError) Error() string {
    return fmt.Sprintf("错误码%d: %s", e.Code, e.Message)
}

func doSomething() error {
    return &MyError{Code: 404, Message: "未找到"}
}

func main() {
    err := doSomething()
    
    var myErr *MyError
    if errors.As(err, &myErr) {
        fmt.Println("错误码:", myErr.Code)      // 404
        fmt.Println("错误信息:", myErr.Message) // 未找到
    }
}
```

**简单记**：
- `errors.Is` - 判断是不是某个错误
- `errors.As` - 把错误转成某个类型

---

## 7. 自定义错误类型

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateUser(name string, age int) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "不能为空"}
    }
    if age < 0 {
        return &ValidationError{Field: "age", Message: "不能为负数"}
    }
    return nil
}

func main() {
    err := validateUser("", 20)
    if err != nil {
        fmt.Println(err)  // name: 不能为空
        
        var ve *ValidationError
        if errors.As(err, &ve) {
            fmt.Println("字段:", ve.Field)  // name
        }
    }
}
```

---

## 8. panic 和 recover

### 8.1 panic - 程序崩溃

```go
func main() {
    panic("出大事了！")
    fmt.Println("这行不会执行")
}
```

**什么时候用 panic**：
- 程序遇到无法恢复的错误
- 初始化失败（比如配置文件不存在）
- 不应该发生的情况（程序 bug）

**一般情况用 error，不要用 panic**。

### 8.2 recover - 捕获 panic

```go
func mayPanic() {
    panic("出错了！")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("捕获到 panic:", r)
        }
    }()
    
    mayPanic()
    
    fmt.Println("这行不会执行")
}
```

输出：
```
捕获到 panic: 出错了！
```

### 8.3 实际用法：防止服务崩溃

```go
func handleRequest() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("请求处理出错:", r)
            // 记录日志，返回 500 错误
        }
    }()
    
    // 处理请求...
    // 如果这里 panic 了，不会让整个服务挂掉
}
```

---

## 9. 错误处理最佳实践

### 9.1 不要忽略错误

```go
// ❌ 错误：忽略错误
file, _ := os.Open("file.txt")

// ✅ 正确：处理错误
file, err := os.Open("file.txt")
if err != nil {
    return err
}
```

### 9.2 错误信息要有上下文

```go
// ❌ 不好：不知道哪里出错
return err

// ✅ 好：知道是什么操作出错
return fmt.Errorf("读取用户配置失败: %w", err)
```

### 9.3 在合适的层级处理错误

```go
// 底层：返回错误
func readFile(path string) ([]byte, error) {
    return os.ReadFile(path)
}

// 中层：包装错误
func loadConfig() (*Config, error) {
    data, err := readFile("config.json")
    if err != nil {
        return nil, fmt.Errorf("加载配置失败: %w", err)
    }
    // ...
}

// 顶层：处理错误（打印、记录日志等）
func main() {
    config, err := loadConfig()
    if err != nil {
        log.Fatal(err)  // 在最顶层处理
    }
}
```

### 9.4 用预定义错误方便判断

```go
var ErrUserNotFound = errors.New("用户不存在")

func getUser(id int) (*User, error) {
    // ...
    return nil, ErrUserNotFound
}

func main() {
    user, err := getUser(123)
    if errors.Is(err, ErrUserNotFound) {
        // 特殊处理：用户不存在
    }
}
```

---

## 10. 常见错误处理场景

### 10.1 文件操作

```go
func readConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("配置文件不存在: %s", path)
        }
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    return data, nil
}
```

### 10.2 HTTP 请求

```go
func fetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("请求失败: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("请求返回错误状态码: %d", resp.StatusCode)
    }
    
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("读取响应失败: %w", err)
    }
    
    return data, nil
}
```

### 10.3 JSON 解析

```go
func parseUser(data []byte) (*User, error) {
    var user User
    err := json.Unmarshal(data, &user)
    if err != nil {
        return nil, fmt.Errorf("解析用户数据失败: %w", err)
    }
    return &user, nil
}
```

---

## 11. 练习题

### 练习 1：基本错误处理

```go
// 实现一个除法函数，除数为 0 时返回错误
func divide(a, b float64) (float64, error) {
    // 你来实现
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("错误:", err)
    } else {
        fmt.Println("结果:", result)
    }
}
```

### 练习 2：错误包装

```go
// 实现一个读取文件并解析 JSON 的函数
// 要求：错误信息要包含上下文
func loadUserFromFile(path string) (*User, error) {
    // 你来实现
}
```

### 练习 3：自定义错误

```go
// 实现一个验证函数，返回自定义的 ValidationError
type ValidationError struct {
    Field   string
    Message string
}

func validateAge(age int) error {
    // 你来实现
    // age < 0: "age 不能为负数"
    // age > 150: "age 不能超过150"
}
```

---

## 总结

| 概念 | 说明 |
|------|------|
| `error` | 接口，只有 `Error() string` 方法 |
| `errors.New()` | 创建简单错误 |
| `fmt.Errorf()` | 创建格式化错误 |
| `fmt.Errorf("%w", err)` | 包装错误 |
| `errors.Is(err, target)` | 判断是否是某个错误 |
| `errors.As(err, &target)` | 转换成某个错误类型 |
| `panic` | 程序崩溃（少用） |
| `recover` | 捕获 panic |

**核心原则**：
1. 不要忽略错误
2. 错误信息要有上下文
3. 用 `%w` 包装错误
4. 用 `errors.Is` 判断错误类型
