# 错误处理最佳实践

## 1. 基本错误处理

```go
// Go 的错误是值
type error interface {
    Error() string
}

// 基本模式
result, err := someFunction()
if err != nil {
    return err  // 或处理错误
}
// 使用 result

// 不要忽略错误
_ = someFunction()  // 不好！
```

## 2. 创建错误

```go
import "errors"
import "fmt"

// 方式1: errors.New
err := errors.New("something went wrong")

// 方式2: fmt.Errorf
err := fmt.Errorf("failed to process %s: %v", name, originalErr)

// 方式3: 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
```

## 3. 错误包装（Go 1.13+）

```go
// 包装错误，保留原始错误信息
originalErr := errors.New("database connection failed")
wrappedErr := fmt.Errorf("failed to get user: %w", originalErr)

// 解包错误
fmt.Println(errors.Unwrap(wrappedErr))  // database connection failed

// 检查错误链中是否包含特定错误
if errors.Is(wrappedErr, originalErr) {
    fmt.Println("contains original error")
}

// 检查错误链中是否包含特定类型
var validationErr *ValidationError
if errors.As(wrappedErr, &validationErr) {
    fmt.Println("field:", validationErr.Field)
}
```

## 4. 哨兵错误

```go
// 定义包级别的错误变量
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)

// 使用
func GetUser(id int) (*User, error) {
    user := db.Find(id)
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}

// 检查
user, err := GetUser(123)
if errors.Is(err, ErrNotFound) {
    // 处理未找到的情况
}
```


## 5. 错误处理模式

```go
// 模式1: 立即返回
func process() error {
    if err := step1(); err != nil {
        return fmt.Errorf("step1 failed: %w", err)
    }
    if err := step2(); err != nil {
        return fmt.Errorf("step2 failed: %w", err)
    }
    return nil
}

// 模式2: 延迟处理
func processFile(path string) (err error) {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr
        }
    }()
    // 处理文件...
    return nil
}

// 模式3: 错误聚合
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    var msgs []string
    for _, err := range m.Errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.Errors = append(m.Errors, err)
    }
}

func (m *MultiError) ErrorOrNil() error {
    if len(m.Errors) == 0 {
        return nil
    }
    return m
}
```

## 6. panic 和 recover

```go
// panic 用于不可恢复的错误
func mustParse(s string) int {
    n, err := strconv.Atoi(s)
    if err != nil {
        panic(fmt.Sprintf("invalid number: %s", s))
    }
    return n
}

// recover 用于捕获 panic
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    fn()
    return nil
}

// HTTP 服务中的 recover 中间件
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## 7. 错误日志

```go
// 结构化日志
import "log/slog"

func processOrder(orderID string) error {
    logger := slog.With("orderID", orderID)
    
    if err := validateOrder(orderID); err != nil {
        logger.Error("validation failed", "error", err)
        return fmt.Errorf("validate order: %w", err)
    }
    
    logger.Info("order processed successfully")
    return nil
}

// 错误追踪
type tracedError struct {
    err   error
    stack []byte
}

func Trace(err error) error {
    if err == nil {
        return nil
    }
    return &tracedError{
        err:   err,
        stack: debug.Stack(),
    }
}

func (e *tracedError) Error() string {
    return e.err.Error()
}

func (e *tracedError) Unwrap() error {
    return e.err
}
```

## 8. 最佳实践

```go
// 1. 错误信息要有上下文
// 不好
return err

// 好
return fmt.Errorf("failed to get user %d: %w", userID, err)

// 2. 只处理一次错误
// 不好
if err != nil {
    log.Println(err)  // 记录
    return err        // 又返回
}

// 好：要么处理，要么返回
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// 3. 使用 errors.Is 而不是 ==
// 不好
if err == sql.ErrNoRows {}

// 好
if errors.Is(err, sql.ErrNoRows) {}

// 4. 错误变量命名
var ErrNotFound = errors.New("not found")  // 以 Err 开头

// 5. 错误类型命名
type NotFoundError struct{}  // 以 Error 结尾
```

## 9. 常见面试题

**Q1: error 和 panic 的区别？什么时候用 panic？**
```
error: 可预期的错误，调用者可以处理
panic: 不可恢复的错误，程序无法继续

使用 panic 的场景：
1. 程序初始化失败（配置错误、必要资源不可用）
2. 编程错误（数组越界、空指针）
3. 不可能发生的情况（用于断言）
```

**Q2: 如何实现错误链？**
```go
// 使用 %w 包装
err := fmt.Errorf("outer: %w", innerErr)

// 使用 errors.Is 检查
if errors.Is(err, targetErr) {}

// 使用 errors.As 获取特定类型
var myErr *MyError
if errors.As(err, &myErr) {}
```

## 练习

1. 实现一个带重试的函数，支持自定义重试次数和间隔
2. 实现一个错误收集器，收集多个 goroutine 的错误
3. 为你的项目设计一套错误码系统
