# Context 正确使用

## 1. Context 是什么

Context 用于在 API 边界和进程之间传递：
- 截止时间（Deadline）
- 取消信号（Cancellation）
- 请求范围的值（Request-scoped values）

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

## 2. 创建 Context

```go
// 根 Context（永不取消）
ctx := context.Background()  // 主函数、初始化
ctx := context.TODO()        // 不确定用什么时的占位符

// 带取消的 Context
ctx, cancel := context.WithCancel(parent)
defer cancel()  // 必须调用！

// 带超时的 Context
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// 带截止时间的 Context
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()

// 带值的 Context
ctx := context.WithValue(parent, key, value)
```


## 3. 取消传播

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx, "worker1")
    go worker(ctx, "worker2")
    
    time.Sleep(2 * time.Second)
    cancel()  // 取消所有 worker
    
    time.Sleep(time.Second)  // 等待 worker 退出
}

func worker(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: stopped, reason: %v\n", name, ctx.Err())
            return
        default:
            fmt.Printf("%s: working...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

## 4. 超时控制

```go
func fetchWithTimeout(url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return nil, fmt.Errorf("request timeout")
        }
        return nil, err
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}
```

## 5. 传递请求范围的值

```go
// 定义 key 类型（避免冲突）
type contextKey string

const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

// 设置值
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

// 获取值
func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}

// 使用
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := WithUserID(r.Context(), "user123")
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    if userID, ok := GetUserID(ctx); ok {
        fmt.Println("User:", userID)
    }
}
```

## 6. 最佳实践

```go
// 1. Context 作为第一个参数
func DoSomething(ctx context.Context, arg string) error {
    // ...
}

// 2. 不要存储 Context
// 不好
type Server struct {
    ctx context.Context  // 不要这样做！
}

// 好
func (s *Server) Handle(ctx context.Context) {
    // 每次调用传入 ctx
}

// 3. 总是调用 cancel
ctx, cancel := context.WithTimeout(parent, time.Second)
defer cancel()  // 即使操作成功也要调用

// 4. 检查 ctx.Done()
func longOperation(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 继续工作
        }
    }
    return nil
}

// 5. 不要传 nil Context
// 不好
DoSomething(nil, "arg")

// 好
DoSomething(context.Background(), "arg")
DoSomething(context.TODO(), "arg")
```

## 7. 常见错误

```go
// 错误1: 忘记调用 cancel
func bad() {
    ctx, _ := context.WithCancel(context.Background())
    // 内存泄漏！
}

// 错误2: 在 goroutine 中使用父 context
func bad(ctx context.Context) {
    go func() {
        // 如果父 ctx 被取消，这里也会被取消
        // 可能不是你想要的
        doWork(ctx)
    }()
}

// 正确: 根据需要创建新的 context
func good(ctx context.Context) {
    go func() {
        // 独立的 context
        newCtx := context.Background()
        doWork(newCtx)
    }()
}

// 错误3: 用 context 传递可选参数
// 不好
ctx = context.WithValue(ctx, "debug", true)

// 好: 使用函数参数或配置结构体
func DoSomething(ctx context.Context, opts Options) {}
```

## 8. 数据库操作中的 Context

```go
func queryUser(ctx context.Context, db *sql.DB, id int) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
    
    var user User
    if err := row.Scan(&user.ID, &user.Name); err != nil {
        return nil, err
    }
    return &user, nil
}
```

## 9. HTTP 中间件中的 Context

```go
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## 10. 面试题

**Q1: context.Background() 和 context.TODO() 的区别？**
```
功能上没有区别，都返回空 Context
语义上：
- Background: 用于主函数、初始化、测试
- TODO: 不确定用什么 Context 时的占位符
```

**Q2: 为什么要用自定义类型作为 Context key？**
```go
// 避免不同包之间的 key 冲突
type myKey string
const key myKey = "myKey"

// 如果用 string，可能和其他包冲突
ctx.Value("userID")  // 可能冲突
```

## 练习

1. 实现一个带超时的 HTTP 客户端
2. 实现一个可取消的文件下载器
3. 在 HTTP 服务中实现请求追踪（传递 request ID）
